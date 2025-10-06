package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultRequestLogConfig(t *testing.T) {
	config := DefaultRequestLogConfig()

	assert.Contains(t, config.SkipPaths, "/health")
	assert.Contains(t, config.SkipPaths, "/metrics")
	assert.Contains(t, config.SkipPaths, "/favicon.ico")

	assert.Contains(t, config.SkipHeaders, "Authorization")
	assert.Contains(t, config.SkipHeaders, "Cookie")
	assert.Contains(t, config.SkipHeaders, "X-API-Key")

	assert.True(t, config.LogRequestBody)
	assert.False(t, config.LogResponseBody)
	assert.Equal(t, 1024, config.MaxBodySize)
}

func TestRequestLogger(t *testing.T) {
	// ログ出力をキャプチャ
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)
	SetGlobalLogger(logger)

	// Ginをテストモードに設定
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// リクエストログミドルウェアを追加
	router.Use(RequestLogger(nil))

	// テストルート
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// リクエストを送信
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "test-request-id")
	req.Header.Set("X-User-ID", "test-user-id")
	req.Header.Set("X-Trace-ID", "test-trace-id")

	router.ServeHTTP(w, req)

	// ログを確認
	logOutput := buf.String()
	assert.NotEmpty(t, logOutput)

	// JSONログをパース
	var entry LogEntry
	err := json.Unmarshal([]byte(logOutput), &entry)
	require.NoError(t, err)

	assert.Equal(t, "GET /test", entry.Message)
	assert.Equal(t, INFO, entry.Level)
	assert.Equal(t, "test-request-id", entry.RequestID)
	assert.Equal(t, "test-user-id", entry.UserID)
	assert.Equal(t, "test-trace-id", entry.TraceID)
	assert.Equal(t, "GET", entry.Fields["method"])
	assert.Equal(t, "/test", entry.Fields["path"])
	assert.Equal(t, float64(200), entry.Fields["status"])
}

func TestRequestLoggerSkipPaths(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)
	SetGlobalLogger(logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	config := DefaultRequestLogConfig()
	router.Use(RequestLogger(config))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	router.ServeHTTP(w, req)

	// スキップパスなのでログは出力されない
	assert.Empty(t, buf.String())
}

func TestRequestLoggerWithRequestBody(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)
	SetGlobalLogger(logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	config := DefaultRequestLogConfig()
	config.LogRequestBody = true
	router.Use(RequestLogger(config))

	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	reqBody := `{"key": "value"}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	logOutput := buf.String()
	assert.NotEmpty(t, logOutput)

	var entry LogEntry
	err := json.Unmarshal([]byte(logOutput), &entry)
	require.NoError(t, err)

	assert.Equal(t, "POST /test", entry.Message)
	assert.Contains(t, entry.Fields["request_body"], "key")
	assert.Contains(t, entry.Fields["request_body"], "value")
}

func TestRequestLoggerWithLargeBody(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)
	SetGlobalLogger(logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	config := DefaultRequestLogConfig()
	config.LogRequestBody = true
	config.MaxBodySize = 10 // 小さなサイズ制限
	router.Use(RequestLogger(config))

	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	reqBody := `{"key": "very_long_value_that_exceeds_max_body_size"}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	logOutput := buf.String()
	assert.NotEmpty(t, logOutput)

	var entry LogEntry
	err := json.Unmarshal([]byte(logOutput), &entry)
	require.NoError(t, err)

	requestBody := entry.Fields["request_body"].(string)
	assert.True(t, len(requestBody) <= config.MaxBodySize+3) // +3 for "..."
	assert.Contains(t, requestBody, "...")
}

func TestResponseLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)
	SetGlobalLogger(logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	config := DefaultRequestLogConfig()
	config.LogResponseBody = true
	router.Use(ResponseLogger(config))

	router.GET("/error", func(c *gin.Context) {
		c.JSON(400, gin.H{"error": "bad request"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	req.Header.Set("X-Request-ID", "test-request-id")

	router.ServeHTTP(w, req)

	// エラーレスポンスなのでログが出力される
	logOutput := buf.String()
	assert.NotEmpty(t, logOutput)

	var entry LogEntry
	err := json.Unmarshal([]byte(logOutput), &entry)
	require.NoError(t, err)

	assert.Equal(t, "Response error", entry.Message)
	assert.Equal(t, WARN, entry.Level)
	assert.Equal(t, float64(400), entry.Fields["status"])
	assert.Contains(t, entry.Fields["response_body"], "bad request")
}

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		requestID := c.GetString("request_id")
		c.JSON(200, gin.H{"request_id": requestID})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	requestID := response["request_id"].(string)
	assert.NotEmpty(t, requestID)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}

func TestTraceIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(TraceIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		traceID := c.GetString("trace_id")
		c.JSON(200, gin.H{"trace_id": traceID})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	traceID := response["trace_id"].(string)
	assert.NotEmpty(t, traceID)
	assert.NotEmpty(t, w.Header().Get("X-Trace-ID"))
}

func TestUserIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(UserIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		userID := c.GetString("user_id")
		c.JSON(200, gin.H{"user_id": userID})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-User-ID", "test-user-id")

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	userID := response["user_id"].(string)
	assert.Equal(t, "test-user-id", userID)
}

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)
	SetGlobalLogger(logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(LoggingMiddleware(nil))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "test-request-id")
	req.Header.Set("X-User-ID", "test-user-id")
	req.Header.Set("X-Trace-ID", "test-trace-id")

	router.ServeHTTP(w, req)

	// 開始と完了の両方のログが出力される
	logOutput := buf.String()
	lines := strings.Split(strings.TrimSpace(logOutput), "\n")
	assert.Len(t, lines, 2)

	// 開始ログ
	var startEntry LogEntry
	err := json.Unmarshal([]byte(lines[0]), &startEntry)
	require.NoError(t, err)
	assert.Equal(t, "Request started", startEntry.Message)
	assert.Equal(t, "test-request-id", startEntry.RequestID)
	assert.Equal(t, "test-user-id", startEntry.UserID)
	assert.Equal(t, "test-trace-id", startEntry.TraceID)

	// 完了ログ
	var completeEntry LogEntry
	err = json.Unmarshal([]byte(lines[1]), &completeEntry)
	require.NoError(t, err)
	assert.Equal(t, "Request completed: GET /test", completeEntry.Message)
	assert.Equal(t, float64(200), completeEntry.Fields["status"])
}

func TestErrorLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)
	SetGlobalLogger(logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(ErrorLogger())
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	req.Header.Set("X-Request-ID", "test-request-id")

	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)

	// パニックログが出力される
	logOutput := buf.String()
	assert.NotEmpty(t, logOutput)

	var entry LogEntry
	err := json.Unmarshal([]byte(logOutput), &entry)
	require.NoError(t, err)

	assert.Equal(t, "Panic recovered", entry.Message)
	assert.Equal(t, ERROR, entry.Level)
	assert.Equal(t, "test-request-id", entry.RequestID)
	assert.Contains(t, entry.Fields["error"], "test panic")
}
