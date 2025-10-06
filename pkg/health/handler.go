package health

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"tea-logistics/pkg/logger"

	"github.com/gin-gonic/gin"
)

// HealthHandler ヘルスチェックハンドラー
type HealthHandler struct {
	healthChecker *HealthChecker
	metricsManager *MetricsManager
}

// NewHealthHandler 新しいヘルスチェックハンドラーを作成
func NewHealthHandler(healthChecker *HealthChecker, metricsManager *MetricsManager) *HealthHandler {
	return &HealthHandler{
		healthChecker:  healthChecker,
		metricsManager: metricsManager,
	}
}

// HealthResponse ヘルスチェックレスポンス
type HealthResponse struct {
	Status    HealthStatus           `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration"`
	Checks    map[string]HealthResult `json:"checks,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// MetricsResponse メトリクスレスポンス
type MetricsResponse struct {
	Timestamp time.Time     `json:"timestamp"`
	Metrics   map[string]*Metric `json:"metrics"`
}

// HealthCheck ヘルスチェックエンドポイント
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	start := time.Now()
	ctx := c.Request.Context()

	// ヘルスチェックを実行
	results := h.healthChecker.CheckAll(ctx)
	overallStatus := h.healthChecker.GetOverallStatus(results)
	duration := time.Since(start)

	// レスポンスの作成
	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Duration:  duration,
		Checks:    results,
		Details: map[string]interface{}{
			"total_checks": len(results),
		},
	}

	// ログ出力
	logger.WithRequestID(c.GetString("request_id")).
		WithUserID(c.GetString("user_id")).
		Info("ヘルスチェック実行", map[string]interface{}{
			"status":        overallStatus,
			"duration":      duration.String(),
			"total_checks":  len(results),
		})

	// ステータスコードの決定
	statusCode := http.StatusOK
	if overallStatus == StatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	} else if overallStatus == StatusDegraded {
		statusCode = http.StatusOK // 部分的に利用可能
	}

	c.JSON(statusCode, response)
}

// LivenessCheck ライブネスチェックエンドポイント
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	// 基本的なライブネスチェック
	response := HealthResponse{
		Status:    StatusHealthy,
		Timestamp: time.Now(),
		Duration:  0,
		Details: map[string]interface{}{
			"message": "アプリケーションは稼働中です",
		},
	}

	logger.WithRequestID(c.GetString("request_id")).
		Info("ライブネスチェック実行")

	c.JSON(http.StatusOK, response)
}

// ReadinessCheck レディネスチェックエンドポイント
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	start := time.Now()
	ctx := c.Request.Context()

	// 重要なヘルスチェックのみを実行
	criticalChecks := []string{"database"}
	results := make(map[string]HealthResult)

	for _, checkName := range criticalChecks {
		if result, exists := h.healthChecker.Check(ctx, checkName); exists {
			results[checkName] = result
		}
	}

	overallStatus := h.healthChecker.GetOverallStatus(results)
	duration := time.Since(start)

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Duration:  duration,
		Checks:    results,
		Details: map[string]interface{}{
			"message": "レディネスチェック完了",
		},
	}

	// ステータスコードの決定
	statusCode := http.StatusOK
	if overallStatus == StatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	logger.WithRequestID(c.GetString("request_id")).
		Info("レディネスチェック実行", map[string]interface{}{
			"status":   overallStatus,
			"duration": duration.String(),
		})

	c.JSON(statusCode, response)
}

// Metrics メトリクスエンドポイント
func (h *HealthHandler) Metrics(c *gin.Context) {
	// メトリクスを収集
	metrics := h.metricsManager.CollectAllMetrics()

	response := MetricsResponse{
		Timestamp: time.Now(),
		Metrics:   metrics,
	}

	logger.WithRequestID(c.GetString("request_id")).
		Debug("メトリクス取得", map[string]interface{}{
			"metric_count": len(metrics),
		})

	c.JSON(http.StatusOK, response)
}

// PrometheusMetrics Prometheus形式のメトリクスエンドポイント
func (h *HealthHandler) PrometheusMetrics(c *gin.Context) {
	// メトリクスを収集
	metrics := h.metricsManager.CollectAllMetrics()

	// Prometheus形式で出力
	var output string
	for name, metric := range metrics {
		// ラベルの処理
		labels := ""
		if len(metric.Labels) > 0 {
			var labelPairs []string
			for key, value := range metric.Labels {
				labelPairs = append(labelPairs, fmt.Sprintf(`%s="%s"`, key, value))
			}
			labels = "{" + strings.Join(labelPairs, ",") + "}"
		}

		// Prometheus形式のメトリクス行
		output += fmt.Sprintf("# HELP %s %s\n", name, metric.Name)
		output += fmt.Sprintf("# TYPE %s %s\n", name, metric.Type)
		output += fmt.Sprintf("%s%s %f\n", name, labels, metric.Value)
	}

	c.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	c.String(http.StatusOK, output)
}

// HealthCheckDetail 特定のヘルスチェック詳細エンドポイント
func (h *HealthHandler) HealthCheckDetail(c *gin.Context) {
	checkName := c.Param("name")
	ctx := c.Request.Context()

	result, exists := h.healthChecker.Check(ctx, checkName)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ヘルスチェックが見つかりません",
			"name":  checkName,
		})
		return
	}

	logger.WithRequestID(c.GetString("request_id")).
		Info("ヘルスチェック詳細取得", map[string]interface{}{
			"check_name": checkName,
			"status":     result.Status,
		})

	c.JSON(http.StatusOK, result)
}

// SystemInfo システム情報エンドポイント
func (h *HealthHandler) SystemInfo(c *gin.Context) {
	// システムメトリクスを取得
	systemMetrics := h.metricsManager.GetSystemMetrics()
	metrics := systemMetrics.GetMetrics()

	// システム情報を構築
	systemInfo := map[string]interface{}{
		"timestamp": time.Now(),
		"metrics":   metrics,
		"runtime": map[string]interface{}{
			"go_version": runtime.Version(),
			"go_os":      runtime.GOOS,
			"go_arch":    runtime.GOARCH,
		},
	}

	logger.WithRequestID(c.GetString("request_id")).
		Info("システム情報取得")

	c.JSON(http.StatusOK, systemInfo)
}

// SetupHealthRoutes ヘルスチェックルートを設定
func SetupHealthRoutes(router *gin.Engine, healthChecker *HealthChecker, metricsManager *MetricsManager) {
	handler := NewHealthHandler(healthChecker, metricsManager)

	// ヘルスチェックルートグループ
	health := router.Group("/health")
	{
		// 基本的なヘルスチェック
		health.GET("/", handler.HealthCheck)
		
		// Kubernetes用のヘルスチェック
		health.GET("/live", handler.LivenessCheck)
		health.GET("/ready", handler.ReadinessCheck)
		
		// 特定のヘルスチェック詳細
		health.GET("/check/:name", handler.HealthCheckDetail)
		
		// システム情報
		health.GET("/system", handler.SystemInfo)
	}

	// メトリクスルートグループ
	metrics := router.Group("/metrics")
	{
		// JSON形式のメトリクス
		metrics.GET("/", handler.Metrics)
		
		// Prometheus形式のメトリクス
		metrics.GET("/prometheus", handler.PrometheusMetrics)
	}
}

// グローバルヘルスハンドラー
var globalHealthHandler *HealthHandler

// InitGlobalHealthHandler グローバルヘルスハンドラーを初期化
func InitGlobalHealthHandler() {
	globalHealthHandler = NewHealthHandler(
		GetGlobalHealthChecker(),
		GetGlobalMetricsManager(),
	)
}

// GetGlobalHealthHandler グローバルヘルスハンドラーを取得
func GetGlobalHealthHandler() *HealthHandler {
	if globalHealthHandler == nil {
		InitGlobalHealthHandler()
	}
	return globalHealthHandler
}

// SetupGlobalHealthRoutes グローバルヘルスルートを設定
func SetupGlobalHealthRoutes(router *gin.Engine) {
	SetupHealthRoutes(router, GetGlobalHealthChecker(), GetGlobalMetricsManager())
}
