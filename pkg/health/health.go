package health

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"tea-logistics/pkg/logger"
)

// HealthStatus ヘルスステータス
type HealthStatus string

const (
	StatusHealthy   HealthStatus = "healthy"
	StatusDegraded  HealthStatus = "degraded"
	StatusUnhealthy HealthStatus = "unhealthy"
)

// HealthCheck ヘルスチェックインターフェース
type HealthCheck interface {
	Name() string
	Check(ctx context.Context) HealthResult
	Timeout() time.Duration
}

// HealthResult ヘルスチェック結果
type HealthResult struct {
	Name      string                 `json:"name"`
	Status    HealthStatus           `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Duration  time.Duration          `json:"duration"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// HealthChecker ヘルスチェッカー
type HealthChecker struct {
	checks map[string]HealthCheck
	mu     sync.RWMutex
}

// NewHealthChecker 新しいヘルスチェッカーを作成
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]HealthCheck),
	}
}

// RegisterCheck ヘルスチェックを登録
func (h *HealthChecker) RegisterCheck(check HealthCheck) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[check.Name()] = check
}

// UnregisterCheck ヘルスチェックを削除
func (h *HealthChecker) UnregisterCheck(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.checks, name)
}

// CheckAll 全てのヘルスチェックを実行
func (h *HealthChecker) CheckAll(ctx context.Context) map[string]HealthResult {
	h.mu.RLock()
	checks := make(map[string]HealthCheck)
	for name, check := range h.checks {
		checks[name] = check
	}
	h.mu.RUnlock()

	results := make(map[string]HealthResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, check := range checks {
		wg.Add(1)
		go func(name string, check HealthCheck) {
			defer wg.Done()

			result := h.runCheck(ctx, check)

			mu.Lock()
			results[name] = result
			mu.Unlock()
		}(name, check)
	}

	wg.Wait()
	return results
}

// Check 特定のヘルスチェックを実行
func (h *HealthChecker) Check(ctx context.Context, name string) (HealthResult, bool) {
	h.mu.RLock()
	check, exists := h.checks[name]
	h.mu.RUnlock()

	if !exists {
		return HealthResult{
			Name:      name,
			Status:    StatusUnhealthy,
			Message:   "ヘルスチェックが見つかりません",
			Timestamp: time.Now(),
		}, false
	}

	return h.runCheck(ctx, check), true
}

// runCheck ヘルスチェックを実行
func (h *HealthChecker) runCheck(ctx context.Context, check HealthCheck) HealthResult {
	start := time.Now()

	// タイムアウト付きコンテキストを作成
	timeout := check.Timeout()
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	checkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// ヘルスチェックを実行
	result := check.Check(checkCtx)
	result.Duration = time.Since(start)
	result.Timestamp = time.Now()

	// ログ出力
	if result.Status == StatusUnhealthy {
		logger.Error("ヘルスチェック完了", map[string]interface{}{
			"name":     result.Name,
			"status":   result.Status,
			"message":  result.Message,
			"duration": result.Duration.String(),
		})
	} else if result.Status == StatusDegraded {
		logger.Warn("ヘルスチェック完了", map[string]interface{}{
			"name":     result.Name,
			"status":   result.Status,
			"message":  result.Message,
			"duration": result.Duration.String(),
		})
	} else {
		logger.Info("ヘルスチェック完了", map[string]interface{}{
			"name":     result.Name,
			"status":   result.Status,
			"message":  result.Message,
			"duration": result.Duration.String(),
		})
	}

	return result
}

// GetOverallStatus 全体のヘルスステータスを取得
func (h *HealthChecker) GetOverallStatus(results map[string]HealthResult) HealthStatus {
	if len(results) == 0 {
		return StatusUnhealthy
	}

	hasUnhealthy := false
	hasDegraded := false

	for _, result := range results {
		switch result.Status {
		case StatusUnhealthy:
			hasUnhealthy = true
		case StatusDegraded:
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return StatusUnhealthy
	}
	if hasDegraded {
		return StatusDegraded
	}
	return StatusHealthy
}

// DatabaseHealthCheck データベースヘルスチェック
type DatabaseHealthCheck struct {
	db      *sql.DB
	name    string
	timeout time.Duration
}

// NewDatabaseHealthCheck 新しいデータベースヘルスチェックを作成
func NewDatabaseHealthCheck(name string, db *sql.DB) *DatabaseHealthCheck {
	return &DatabaseHealthCheck{
		db:      db,
		name:    name,
		timeout: 10 * time.Second,
	}
}

// Name ヘルスチェック名を取得
func (d *DatabaseHealthCheck) Name() string {
	return d.name
}

// Timeout タイムアウトを取得
func (d *DatabaseHealthCheck) Timeout() time.Duration {
	return d.timeout
}

// Check データベースヘルスチェックを実行
func (d *DatabaseHealthCheck) Check(ctx context.Context) HealthResult {
	// データベース接続の確認
	if err := d.db.PingContext(ctx); err != nil {
		return HealthResult{
			Name:    d.name,
			Status:  StatusUnhealthy,
			Message: fmt.Sprintf("データベース接続エラー: %v", err),
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		}
	}

	// 基本的なクエリの実行
	var result int
	if err := d.db.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		return HealthResult{
			Name:    d.name,
			Status:  StatusDegraded,
			Message: fmt.Sprintf("データベースクエリエラー: %v", err),
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		}
	}

	// データベース統計情報の取得
	stats := d.db.Stats()

	return HealthResult{
		Name:    d.name,
		Status:  StatusHealthy,
		Message: "データベースは正常です",
		Details: map[string]interface{}{
			"open_connections":     stats.OpenConnections,
			"in_use":               stats.InUse,
			"idle":                 stats.Idle,
			"wait_count":           stats.WaitCount,
			"wait_duration":        stats.WaitDuration.String(),
			"max_idle_closed":      stats.MaxIdleClosed,
			"max_idle_time_closed": stats.MaxIdleTimeClosed,
			"max_lifetime_closed":  stats.MaxLifetimeClosed,
		},
	}
}

// HTTPHealthCheck HTTPエンドポイントヘルスチェック
type HTTPHealthCheck struct {
	name    string
	url     string
	timeout time.Duration
}

// NewHTTPHealthCheck 新しいHTTPヘルスチェックを作成
func NewHTTPHealthCheck(name, url string) *HTTPHealthCheck {
	return &HTTPHealthCheck{
		name:    name,
		url:     url,
		timeout: 5 * time.Second,
	}
}

// Name ヘルスチェック名を取得
func (h *HTTPHealthCheck) Name() string {
	return h.name
}

// Timeout タイムアウトを取得
func (h *HTTPHealthCheck) Timeout() time.Duration {
	return h.timeout
}

// Check HTTPヘルスチェックを実行
func (h *HTTPHealthCheck) Check(ctx context.Context) HealthResult {
	// HTTPクライアントの作成
	client := &http.Client{
		Timeout: h.timeout,
	}

	// HTTPリクエストの実行
	req, err := http.NewRequestWithContext(ctx, "GET", h.url, nil)
	if err != nil {
		return HealthResult{
			Name:    h.name,
			Status:  StatusUnhealthy,
			Message: fmt.Sprintf("HTTPリクエスト作成エラー: %v", err),
			Details: map[string]interface{}{
				"url":   h.url,
				"error": err.Error(),
			},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return HealthResult{
			Name:    h.name,
			Status:  StatusUnhealthy,
			Message: fmt.Sprintf("HTTPリクエストエラー: %v", err),
			Details: map[string]interface{}{
				"url":   h.url,
				"error": err.Error(),
			},
		}
	}
	defer resp.Body.Close()

	// レスポンスステータスの確認
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return HealthResult{
			Name:    h.name,
			Status:  StatusHealthy,
			Message: "HTTPエンドポイントは正常です",
			Details: map[string]interface{}{
				"url":          h.url,
				"status_code":  resp.StatusCode,
				"content_type": resp.Header.Get("Content-Type"),
			},
		}
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return HealthResult{
			Name:    h.name,
			Status:  StatusDegraded,
			Message: fmt.Sprintf("HTTPエンドポイントがクライアントエラーを返しました: %d", resp.StatusCode),
			Details: map[string]interface{}{
				"url":         h.url,
				"status_code": resp.StatusCode,
			},
		}
	} else {
		return HealthResult{
			Name:    h.name,
			Status:  StatusUnhealthy,
			Message: fmt.Sprintf("HTTPエンドポイントがサーバーエラーを返しました: %d", resp.StatusCode),
			Details: map[string]interface{}{
				"url":         h.url,
				"status_code": resp.StatusCode,
			},
		}
	}
}

// CustomHealthCheck カスタムヘルスチェック
type CustomHealthCheck struct {
	name    string
	checkFn func(ctx context.Context) HealthResult
	timeout time.Duration
}

// NewCustomHealthCheck 新しいカスタムヘルスチェックを作成
func NewCustomHealthCheck(name string, checkFn func(ctx context.Context) HealthResult, timeout time.Duration) *CustomHealthCheck {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &CustomHealthCheck{
		name:    name,
		checkFn: checkFn,
		timeout: timeout,
	}
}

// Name ヘルスチェック名を取得
func (c *CustomHealthCheck) Name() string {
	return c.name
}

// Timeout タイムアウトを取得
func (c *CustomHealthCheck) Timeout() time.Duration {
	return c.timeout
}

// Check カスタムヘルスチェックを実行
func (c *CustomHealthCheck) Check(ctx context.Context) HealthResult {
	return c.checkFn(ctx)
}

// グローバルヘルスチェッカー
var globalHealthChecker *HealthChecker

// InitGlobalHealthChecker グローバルヘルスチェッカーを初期化
func InitGlobalHealthChecker() {
	globalHealthChecker = NewHealthChecker()
}

// GetGlobalHealthChecker グローバルヘルスチェッカーを取得
func GetGlobalHealthChecker() *HealthChecker {
	if globalHealthChecker == nil {
		globalHealthChecker = NewHealthChecker()
	}
	return globalHealthChecker
}

// RegisterGlobalCheck グローバルヘルスチェックを登録
func RegisterGlobalCheck(check HealthCheck) {
	GetGlobalHealthChecker().RegisterCheck(check)
}

// CheckGlobalHealth グローバルヘルスチェックを実行
func CheckGlobalHealth(ctx context.Context) map[string]HealthResult {
	return GetGlobalHealthChecker().CheckAll(ctx)
}
