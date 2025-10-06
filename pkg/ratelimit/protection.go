package ratelimit

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tea-logistics/pkg/logger"

	"github.com/gin-gonic/gin"
)

// APIProtectionConfig API保護設定
type APIProtectionConfig struct {
	Enabled           bool              `json:"enabled"`
	DefaultStrategy   RateLimitStrategy `json:"default_strategy"`
	DefaultLimit      int               `json:"default_limit"`
	DefaultWindow     time.Duration     `json:"default_window"`
	IPWhitelist       []string          `json:"ip_whitelist"`
	IPBlacklist       []string          `json:"ip_blacklist"`
	UserWhitelist     []string          `json:"user_whitelist"`
	UserBlacklist     []string          `json:"user_blacklist"`
	SkipPaths         []string          `json:"skip_paths"`
	CustomLimits      map[string]int    `json:"custom_limits"`
	CustomStrategies  map[string]string `json:"custom_strategies"`
	CustomWindows     map[string]string `json:"custom_windows"`
	ErrorResponse     string            `json:"error_response"`
	Headers           map[string]string `json:"headers"`
}

// DefaultAPIProtectionConfig デフォルトAPI保護設定
func DefaultAPIProtectionConfig() *APIProtectionConfig {
	return &APIProtectionConfig{
		Enabled:         true,
		DefaultStrategy: StrategyFixedWindow,
		DefaultLimit:    100,
		DefaultWindow:   1 * time.Minute,
		IPWhitelist:     []string{},
		IPBlacklist:     []string{},
		UserWhitelist:   []string{},
		UserBlacklist:   []string{},
		SkipPaths:       []string{"/health", "/metrics"},
		CustomLimits:    make(map[string]int),
		CustomStrategies: make(map[string]string),
		CustomWindows:    make(map[string]string),
		ErrorResponse:   "レート制限に達しました。しばらくしてから再試行してください。",
		Headers: map[string]string{
			"X-RateLimit-Limit":     "X-RateLimit-Limit",
			"X-RateLimit-Remaining": "X-RateLimit-Remaining",
			"X-RateLimit-Reset":     "X-RateLimit-Reset",
			"X-RateLimit-Retry":     "X-RateLimit-Retry-After",
		},
	}
}

// APIProtector API保護機能
type APIProtector struct {
	config  *APIProtectionConfig
	manager *RateLimitManager
}

// NewAPIProtector 新しいAPI保護機能を作成
func NewAPIProtector(config *APIProtectionConfig, manager *RateLimitManager) *APIProtector {
	return &APIProtector{
		config:  config,
		manager: manager,
	}
}

// Middleware Ginミドルウェア
func (p *APIProtector) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !p.config.Enabled {
			c.Next()
			return
		}

		// スキップパスのチェック
		if p.shouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// クライアント情報を取得
		clientInfo := p.getClientInfo(c)
		
		// ブラックリストのチェック
		if p.isBlacklisted(clientInfo) {
			logger.Warn("ブラックリストからのアクセス", map[string]interface{}{
				"ip":   clientInfo.IP,
				"user": clientInfo.UserID,
				"path": c.Request.URL.Path,
			})
			c.JSON(http.StatusForbidden, gin.H{
				"error": "アクセスが拒否されました",
			})
			c.Abort()
			return
		}

		// ホワイトリストのチェック
		if p.isWhitelisted(clientInfo) {
			logger.Debug("ホワイトリストからのアクセス", map[string]interface{}{
				"ip":   clientInfo.IP,
				"user": clientInfo.UserID,
				"path": c.Request.URL.Path,
			})
			c.Next()
			return
		}

		// レート制限の設定を取得
		strategy, _, _ := p.getRateLimitConfig(c.Request.URL.Path, clientInfo)

		// レート制限キーを生成
		key := p.generateRateLimitKey(clientInfo, c.Request.URL.Path)

		// レート制限をチェック
		result, err := p.manager.Allow(c.Request.Context(), strategy, key)
		if err != nil {
			logger.Error("レート制限チェックエラー", map[string]interface{}{
				"key":      key,
				"strategy": strategy,
				"error":    err.Error(),
			})
			
			// エラー時は許可する（設定による）
			c.Next()
			return
		}

		// レスポンスヘッダーを設定
		p.setRateLimitHeaders(c, result)

		// レート制限に達した場合
		if !result.Allowed {
			logger.Warn("レート制限に達しました", map[string]interface{}{
				"key":         key,
				"strategy":    strategy,
				"limit":       result.Limit,
				"remaining":   result.Remaining,
				"retry_after": result.RetryAfter.String(),
			})

			c.Header("Retry-After", strconv.Itoa(int(result.RetryAfter.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       p.config.ErrorResponse,
				"retry_after": int(result.RetryAfter.Seconds()),
			})
			c.Abort()
			return
		}

		logger.Debug("レート制限チェック成功", map[string]interface{}{
			"key":       key,
			"strategy":  strategy,
			"limit":     result.Limit,
			"remaining": result.Remaining,
		})

		c.Next()
	}
}

// ClientInfo クライアント情報
type ClientInfo struct {
	IP      string
	UserID  string
	UserAgent string
	Path    string
}

// shouldSkipPath パスをスキップするかチェック
func (p *APIProtector) shouldSkipPath(path string) bool {
	for _, skipPath := range p.config.SkipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// getClientInfo クライアント情報を取得
func (p *APIProtector) getClientInfo(c *gin.Context) *ClientInfo {
	// IPアドレスを取得
	ip := c.ClientIP()
	
	// ユーザーIDを取得（認証済みの場合）
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "anonymous"
	}
	
	// User-Agentを取得
	userAgent := c.GetHeader("User-Agent")
	
	return &ClientInfo{
		IP:        ip,
		UserID:    userID,
		UserAgent: userAgent,
		Path:      c.Request.URL.Path,
	}
}

// isBlacklisted ブラックリストかチェック
func (p *APIProtector) isBlacklisted(clientInfo *ClientInfo) bool {
	// IPブラックリストのチェック
	for _, blacklistedIP := range p.config.IPBlacklist {
		if clientInfo.IP == blacklistedIP {
			return true
		}
	}
	
	// ユーザーブラックリストのチェック
	for _, blacklistedUser := range p.config.UserBlacklist {
		if clientInfo.UserID == blacklistedUser {
			return true
		}
	}
	
	return false
}

// isWhitelisted ホワイトリストかチェック
func (p *APIProtector) isWhitelisted(clientInfo *ClientInfo) bool {
	// IPホワイトリストのチェック
	for _, whitelistedIP := range p.config.IPWhitelist {
		if clientInfo.IP == whitelistedIP {
			return true
		}
	}
	
	// ユーザーホワイトリストのチェック
	for _, whitelistedUser := range p.config.UserWhitelist {
		if clientInfo.UserID == whitelistedUser {
			return true
		}
	}
	
	return false
}

// getRateLimitConfig レート制限設定を取得
func (p *APIProtector) getRateLimitConfig(path string, clientInfo *ClientInfo) (RateLimitStrategy, int, time.Duration) {
	strategy := p.config.DefaultStrategy
	limit := p.config.DefaultLimit
	window := p.config.DefaultWindow
	
	// パス固有の設定をチェック
	if customLimit, exists := p.config.CustomLimits[path]; exists {
		limit = customLimit
	}
	
	if customStrategy, exists := p.config.CustomStrategies[path]; exists {
		strategy = RateLimitStrategy(customStrategy)
	}
	
	if customWindow, exists := p.config.CustomWindows[path]; exists {
		if parsedWindow, err := time.ParseDuration(customWindow); err == nil {
			window = parsedWindow
		}
	}
	
	// ユーザー固有の設定をチェック
	userKey := fmt.Sprintf("user:%s", clientInfo.UserID)
	if customLimit, exists := p.config.CustomLimits[userKey]; exists {
		limit = customLimit
	}
	
	if customStrategy, exists := p.config.CustomStrategies[userKey]; exists {
		strategy = RateLimitStrategy(customStrategy)
	}
	
	if customWindow, exists := p.config.CustomWindows[userKey]; exists {
		if parsedWindow, err := time.ParseDuration(customWindow); err == nil {
			window = parsedWindow
		}
	}
	
	return strategy, limit, window
}

// generateRateLimitKey レート制限キーを生成
func (p *APIProtector) generateRateLimitKey(clientInfo *ClientInfo, path string) string {
	// ユーザーIDが匿名でない場合はユーザーIDを使用
	if clientInfo.UserID != "anonymous" {
		return fmt.Sprintf("user:%s:%s", clientInfo.UserID, path)
	}
	
	// 匿名ユーザーの場合はIPアドレスを使用
	return fmt.Sprintf("ip:%s:%s", clientInfo.IP, path)
}

// setRateLimitHeaders レート制限ヘッダーを設定
func (p *APIProtector) setRateLimitHeaders(c *gin.Context, result *RateLimitResult) {
	// カスタムヘッダー名を使用
	limitHeader := p.config.Headers["X-RateLimit-Limit"]
	remainingHeader := p.config.Headers["X-RateLimit-Remaining"]
	resetHeader := p.config.Headers["X-RateLimit-Reset"]
	retryHeader := p.config.Headers["X-RateLimit-Retry"]
	
	c.Header(limitHeader, strconv.Itoa(result.Limit))
	c.Header(remainingHeader, strconv.Itoa(result.Remaining))
	c.Header(resetHeader, strconv.FormatInt(result.ResetTime.Unix(), 10))
	
	if result.RetryAfter > 0 {
		c.Header(retryHeader, strconv.Itoa(int(result.RetryAfter.Seconds())))
	}
}

// DDoSProtectionConfig DDoS保護設定
type DDoSProtectionConfig struct {
	Enabled           bool          `json:"enabled"`
	MaxRequestsPerMin int           `json:"max_requests_per_min"`
	MaxRequestsPerSec int           `json:"max_requests_per_sec"`
	BlockDuration     time.Duration `json:"block_duration"`
	Whitelist         []string      `json:"whitelist"`
	Blacklist         []string      `json:"blacklist"`
	AutoBlock         bool          `json:"auto_block"`
	AlertThreshold    int           `json:"alert_threshold"`
}

// DefaultDDoSProtectionConfig デフォルトDDoS保護設定
func DefaultDDoSProtectionConfig() *DDoSProtectionConfig {
	return &DDoSProtectionConfig{
		Enabled:           true,
		MaxRequestsPerMin: 1000,
		MaxRequestsPerSec: 100,
		BlockDuration:     5 * time.Minute,
		Whitelist:         []string{},
		Blacklist:         []string{},
		AutoBlock:         true,
		AlertThreshold:    500,
	}
}

// DDoSProtector DDoS保護機能
type DDoSProtector struct {
	config  *DDoSProtectionConfig
	manager *RateLimitManager
}

// NewDDoSProtector 新しいDDoS保護機能を作成
func NewDDoSProtector(config *DDoSProtectionConfig, manager *RateLimitManager) *DDoSProtector {
	return &DDoSProtector{
		config:  config,
		manager: manager,
	}
}

// Middleware Ginミドルウェア
func (p *DDoSProtector) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !p.config.Enabled {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		
		// ホワイトリストのチェック
		if p.isWhitelisted(clientIP) {
			c.Next()
			return
		}
		
		// ブラックリストのチェック
		if p.isBlacklisted(clientIP) {
			logger.Warn("DDoSブラックリストからのアクセス", map[string]interface{}{
				"ip": clientIP,
			})
			c.JSON(http.StatusForbidden, gin.H{
				"error": "アクセスが拒否されました",
			})
			c.Abort()
			return
		}

		// レート制限をチェック
		perMinKey := fmt.Sprintf("ddos:min:%s", clientIP)
		perSecKey := fmt.Sprintf("ddos:sec:%s", clientIP)

		// 分単位のレート制限
		minResult, err := p.manager.Allow(c.Request.Context(), StrategyFixedWindow, perMinKey)
		if err != nil {
			logger.Error("DDoS保護分単位レート制限エラー", map[string]interface{}{
				"ip":    clientIP,
				"error": err.Error(),
			})
			c.Next()
			return
		}

		// 秒単位のレート制限
		secResult, err := p.manager.Allow(c.Request.Context(), StrategyFixedWindow, perSecKey)
		if err != nil {
			logger.Error("DDoS保護秒単位レート制限エラー", map[string]interface{}{
				"ip":    clientIP,
				"error": err.Error(),
			})
			c.Next()
			return
		}

		// DDoS攻撃の検出
		if !minResult.Allowed || !secResult.Allowed {
			logger.Error("DDoS攻撃を検出しました", map[string]interface{}{
				"ip":              clientIP,
				"min_requests":    minResult.Limit - minResult.Remaining,
				"sec_requests":    secResult.Limit - secResult.Remaining,
				"max_per_min":     p.config.MaxRequestsPerMin,
				"max_per_sec":     p.config.MaxRequestsPerSec,
			})

			// 自動ブロック
			if p.config.AutoBlock {
				p.addToBlacklist(clientIP)
			}

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "DDoS攻撃が検出されました",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isWhitelisted ホワイトリストかチェック
func (p *DDoSProtector) isWhitelisted(ip string) bool {
	for _, whitelistedIP := range p.config.Whitelist {
		if ip == whitelistedIP {
			return true
		}
	}
	return false
}

// isBlacklisted ブラックリストかチェック
func (p *DDoSProtector) isBlacklisted(ip string) bool {
	for _, blacklistedIP := range p.config.Blacklist {
		if ip == blacklistedIP {
			return true
		}
	}
	return false
}

// addToBlacklist ブラックリストに追加
func (p *DDoSProtector) addToBlacklist(ip string) {
	p.config.Blacklist = append(p.config.Blacklist, ip)
	logger.Info("IPをDDoSブラックリストに追加しました", map[string]interface{}{
		"ip": ip,
	})
}

// SecurityHeadersConfig セキュリティヘッダー設定
type SecurityHeadersConfig struct {
	Enabled                bool   `json:"enabled"`
	XFrameOptions          string `json:"x_frame_options"`
	XContentTypeOptions    string `json:"x_content_type_options"`
	XSSProtection          string `json:"xss_protection"`
	StrictTransportSecurity string `json:"strict_transport_security"`
	ContentSecurityPolicy  string `json:"content_security_policy"`
	ReferrerPolicy         string `json:"referrer_policy"`
	PermissionsPolicy      string `json:"permissions_policy"`
}

// DefaultSecurityHeadersConfig デフォルトセキュリティヘッダー設定
func DefaultSecurityHeadersConfig() *SecurityHeadersConfig {
	return &SecurityHeadersConfig{
		Enabled:                true,
		XFrameOptions:          "DENY",
		XContentTypeOptions:    "nosniff",
		XSSProtection:          "1; mode=block",
		StrictTransportSecurity: "max-age=31536000; includeSubDomains",
		ContentSecurityPolicy:  "default-src 'self'",
		ReferrerPolicy:         "strict-origin-when-cross-origin",
		PermissionsPolicy:      "geolocation=(), microphone=(), camera=()",
	}
}

// SecurityHeadersMiddleware セキュリティヘッダーミドルウェア
func SecurityHeadersMiddleware(config *SecurityHeadersConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Enabled {
			c.Next()
			return
		}

		// セキュリティヘッダーを設定
		c.Header("X-Frame-Options", config.XFrameOptions)
		c.Header("X-Content-Type-Options", config.XContentTypeOptions)
		c.Header("X-XSS-Protection", config.XSSProtection)
		c.Header("Strict-Transport-Security", config.StrictTransportSecurity)
		c.Header("Content-Security-Policy", config.ContentSecurityPolicy)
		c.Header("Referrer-Policy", config.ReferrerPolicy)
		c.Header("Permissions-Policy", config.PermissionsPolicy)

		c.Next()
	}
}

// グローバルAPI保護機能
var globalAPIProtector *APIProtector
var globalDDoSProtector *DDoSProtector

// InitGlobalAPIProtection グローバルAPI保護機能を初期化
func InitGlobalAPIProtection(config *APIProtectionConfig) {
	manager := GetGlobalRateLimitManager()
	globalAPIProtector = NewAPIProtector(config, manager)
}

// InitGlobalDDoSProtection グローバルDDoS保護機能を初期化
func InitGlobalDDoSProtection(config *DDoSProtectionConfig) {
	manager := GetGlobalRateLimitManager()
	globalDDoSProtector = NewDDoSProtector(config, manager)
}

// GetGlobalAPIProtector グローバルAPI保護機能を取得
func GetGlobalAPIProtector() *APIProtector {
	return globalAPIProtector
}

// GetGlobalDDoSProtector グローバルDDoS保護機能を取得
func GetGlobalDDoSProtector() *DDoSProtector {
	return globalDDoSProtector
}
