package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogConfig リクエストログ設定
type RequestLogConfig struct {
	SkipPaths     []string `json:"skip_paths"`
	SkipHeaders   []string `json:"skip_headers"`
	LogRequestBody bool    `json:"log_request_body"`
	LogResponseBody bool   `json:"log_response_body"`
	MaxBodySize   int      `json:"max_body_size"`
}

// DefaultRequestLogConfig デフォルトのリクエストログ設定
func DefaultRequestLogConfig() *RequestLogConfig {
	return &RequestLogConfig{
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/favicon.ico",
		},
		SkipHeaders: []string{
			"Authorization",
			"Cookie",
			"X-API-Key",
		},
		LogRequestBody:  true,
		LogResponseBody: false,
		MaxBodySize:     1024, // 1KB
	}
}

// RequestLogger リクエストログミドルウェア
func RequestLogger(config *RequestLogConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultRequestLogConfig()
	}
	
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// スキップパスのチェック
			for _, skipPath := range config.SkipPaths {
				if param.Path == skipPath {
					return ""
				}
			}
			
			// リクエストIDの生成
			requestID := param.Request.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}
			
			// ユーザーIDの取得
			userID := param.Request.Header.Get("X-User-ID")
			
			// トレースIDの取得
			traceID := param.Request.Header.Get("X-Trace-ID")
			
			// リクエストボディの取得
			var requestBody string
			if config.LogRequestBody && param.Request.Body != nil {
				bodyBytes, err := io.ReadAll(param.Request.Body)
				if err == nil {
					param.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					if len(bodyBytes) <= config.MaxBodySize {
						requestBody = string(bodyBytes)
					} else {
						requestBody = string(bodyBytes[:config.MaxBodySize]) + "..."
					}
				}
			}
			
			// ヘッダーのフィルタリング
			filteredHeaders := make(map[string]string)
			for name, values := range param.Request.Header {
				skip := false
				for _, skipHeader := range config.SkipHeaders {
					if name == skipHeader {
						skip = true
						break
					}
				}
				if !skip && len(values) > 0 {
					filteredHeaders[name] = values[0]
				}
			}
			
			// ログエントリの作成
			fields := map[string]interface{}{
				"method":      param.Method,
				"path":        param.Path,
				"status":      param.StatusCode,
				"latency":     param.Latency.String(),
				"client_ip":   param.ClientIP,
				"user_agent":  param.Request.UserAgent(),
				"headers":     filteredHeaders,
			}
			
			if requestBody != "" {
				fields["request_body"] = requestBody
			}
			
			// レスポンスボディのログ（エラー時のみ）
			if config.LogResponseBody && param.StatusCode >= 400 {
				// レスポンスボディは後で処理
			}
			
			// ログレベルの決定
			var level LogLevel
			switch {
			case param.StatusCode >= 500:
				level = ERROR
			case param.StatusCode >= 400:
				level = WARN
			default:
				level = INFO
			}
			
			// ログ出力
			logger := GetGlobalLogger().
				WithRequestID(requestID).
				WithUserID(userID).
				WithTraceID(traceID)
			
			message := fmt.Sprintf("%s %s", param.Method, param.Path)
			logger.log(level, message, fields)
			
			return "" // ginのデフォルトログは無効化
		},
		Output: io.Discard, // 出力は上記で処理
	})
}

// ResponseLogger レスポンスログミドルウェア
func ResponseLogger(config *RequestLogConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultRequestLogConfig()
	}
	
	return func(c *gin.Context) {
		// レスポンスボディをキャプチャするためのライター
		blw := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = blw
		
		c.Next()
		
		// エラー時のみレスポンスボディをログ
		if config.LogResponseBody && c.Writer.Status() >= 400 {
			requestID := c.GetHeader("X-Request-ID")
			userID := c.GetHeader("X-User-ID")
			traceID := c.GetHeader("X-Trace-ID")
			
			responseBody := blw.body.String()
			if len(responseBody) > config.MaxBodySize {
				responseBody = responseBody[:config.MaxBodySize] + "..."
			}
			
			logger := GetGlobalLogger().
				WithRequestID(requestID).
				WithUserID(userID).
				WithTraceID(traceID)
			
			logger.Warn("Response error", map[string]interface{}{
				"status":        c.Writer.Status(),
				"path":          c.Request.URL.Path,
				"response_body": responseBody,
			})
		}
	}
}

// bodyLogWriter レスポンスボディをキャプチャするライター
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ErrorLogger エラーログミドルウェア
func ErrorLogger() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := c.GetHeader("X-Request-ID")
		userID := c.GetHeader("X-User-ID")
		traceID := c.GetHeader("X-Trace-ID")
		
		logger := GetGlobalLogger().
			WithRequestID(requestID).
			WithUserID(userID).
			WithTraceID(traceID)
		
		logger.Error("Panic recovered", map[string]interface{}{
			"error": recovered,
			"path":  c.Request.URL.Path,
			"method": c.Request.Method,
		})
		
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// RequestIDMiddleware リクエストIDを生成・設定するミドルウェア
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// TraceIDMiddleware トレースIDを生成・設定するミドルウェア
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		
		c.Header("X-Trace-ID", traceID)
		c.Set("trace_id", traceID)
		c.Next()
	}
}

// UserIDMiddleware ユーザーIDを設定するミドルウェア
func UserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	}
}

// LoggingMiddleware 包括的なログミドルウェア
func LoggingMiddleware(config *RequestLogConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// リクエスト開始時のログ
		start := time.Now()
		
		requestID := c.GetHeader("X-Request-ID")
		userID := c.GetHeader("X-User-ID")
		traceID := c.GetHeader("X-Trace-ID")
		
		logger := GetGlobalLogger().
			WithRequestID(requestID).
			WithUserID(userID).
			WithTraceID(traceID)
		
		logger.Info("Request started", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"query":  c.Request.URL.RawQuery,
			"client_ip": c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})
		
		// リクエスト処理
		c.Next()
		
		// レスポンス完了時のログ
		latency := time.Since(start)
		status := c.Writer.Status()
		
		fields := map[string]interface{}{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  status,
			"latency": latency.String(),
			"size":    c.Writer.Size(),
		}
		
		// ログレベルの決定
		var level LogLevel
		switch {
		case status >= 500:
			level = ERROR
		case status >= 400:
			level = WARN
		default:
			level = INFO
		}
		
		message := fmt.Sprintf("Request completed: %s %s", c.Request.Method, c.Request.URL.Path)
		logger.log(level, message, fields)
	})
}
