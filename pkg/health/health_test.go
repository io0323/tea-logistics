package health

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   HealthStatus
		expected string
	}{
		{"Healthy", StatusHealthy, "healthy"},
		{"Degraded", StatusDegraded, "degraded"},
		{"Unhealthy", StatusUnhealthy, "unhealthy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestHealthResult(t *testing.T) {
	result := HealthResult{
		Name:      "test_check",
		Status:    StatusHealthy,
		Message:   "テスト成功",
		Duration:  100 * time.Millisecond,
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"key": "value",
		},
	}

	assert.Equal(t, "test_check", result.Name)
	assert.Equal(t, StatusHealthy, result.Status)
	assert.Equal(t, "テスト成功", result.Message)
	assert.Equal(t, 100*time.Millisecond, result.Duration)
	assert.NotZero(t, result.Timestamp)
	assert.Equal(t, "value", result.Details["key"])
}

func TestHealthChecker(t *testing.T) {
	checker := NewHealthChecker()

	t.Run("ヘルスチェックの登録", func(t *testing.T) {
		check := &mockHealthCheck{
			name:    "test_check",
			timeout: 5 * time.Second,
		}
		checker.RegisterCheck(check)

		result, exists := checker.Check(context.Background(), "test_check")
		assert.True(t, exists)
		assert.Equal(t, "test_check", result.Name)
		assert.Equal(t, StatusHealthy, result.Status)
	})

	t.Run("存在しないヘルスチェック", func(t *testing.T) {
		result, exists := checker.Check(context.Background(), "nonexistent")
		assert.False(t, exists)
		assert.Equal(t, StatusUnhealthy, result.Status)
		assert.Contains(t, result.Message, "見つかりません")
	})

	t.Run("ヘルスチェックの削除", func(t *testing.T) {
		checker.UnregisterCheck("test_check")
		result, exists := checker.Check(context.Background(), "test_check")
		assert.False(t, exists)
		assert.Equal(t, StatusUnhealthy, result.Status)
	})

	t.Run("全体のヘルスステータス", func(t *testing.T) {
		// 正常なチェック
		healthyCheck := &mockHealthCheck{
			name:    "healthy_check",
			timeout: 5 * time.Second,
		}
		checker.RegisterCheck(healthyCheck)

		// 劣化したチェック
		degradedCheck := &mockHealthCheck{
			name:    "degraded_check",
			timeout: 5 * time.Second,
			status:  StatusDegraded,
		}
		checker.RegisterCheck(degradedCheck)

		results := checker.CheckAll(context.Background())
		overallStatus := checker.GetOverallStatus(results)

		assert.Equal(t, StatusDegraded, overallStatus)
		assert.Len(t, results, 2)
	})
}

func TestDatabaseHealthCheck(t *testing.T) {
	// モックデータベースの作成
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer db.Close()

	check := NewDatabaseHealthCheck("test_db", db)

	t.Run("正常なデータベースチェック", func(t *testing.T) {
		// モックの設定
		mock.ExpectPing()
		mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

		result := check.Check(context.Background())

		assert.Equal(t, "test_db", result.Name)
		assert.Equal(t, StatusHealthy, result.Status)
		assert.Contains(t, result.Message, "正常です")
		assert.NotNil(t, result.Details)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("データベース接続エラー", func(t *testing.T) {
		// 新しいモックデータベース
		db2, mock2, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		defer db2.Close()

		check2 := NewDatabaseHealthCheck("test_db2", db2)

		// Pingエラーを設定
		mock2.ExpectPing().WillReturnError(assert.AnError)

		result := check2.Check(context.Background())

		assert.Equal(t, "test_db2", result.Name)
		assert.Equal(t, StatusUnhealthy, result.Status)
		assert.Contains(t, result.Message, "接続エラー")
		assert.NoError(t, mock2.ExpectationsWereMet())
	})
}

func TestHTTPHealthCheck(t *testing.T) {
	// テストサーバーの作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	check := NewHTTPHealthCheck("test_http", server.URL)

	t.Run("正常なHTTPチェック", func(t *testing.T) {
		result := check.Check(context.Background())

		assert.Equal(t, "test_http", result.Name)
		assert.Equal(t, StatusHealthy, result.Status)
		assert.Contains(t, result.Message, "正常です")
		assert.NotNil(t, result.Details)
		assert.Equal(t, server.URL, result.Details["url"])
		assert.Equal(t, 200, result.Details["status_code"])
	})

	t.Run("存在しないURL", func(t *testing.T) {
		check2 := NewHTTPHealthCheck("test_http2", "http://nonexistent.example.com")
		result := check2.Check(context.Background())

		assert.Equal(t, "test_http2", result.Name)
		assert.Equal(t, StatusUnhealthy, result.Status)
		assert.Contains(t, result.Message, "リクエストエラー")
	})
}

func TestCustomHealthCheck(t *testing.T) {
	checkFn := func(ctx context.Context) HealthResult {
		return HealthResult{
			Name:    "custom_check",
			Status:  StatusHealthy,
			Message: "カスタムチェック成功",
		}
	}

	check := NewCustomHealthCheck("custom_check", checkFn, 5*time.Second)

	assert.Equal(t, "custom_check", check.Name())
	assert.Equal(t, 5*time.Second, check.Timeout())

	result := check.Check(context.Background())
	assert.Equal(t, "custom_check", result.Name)
	assert.Equal(t, StatusHealthy, result.Status)
	assert.Equal(t, "カスタムチェック成功", result.Message)
}

func TestMetricCollector(t *testing.T) {
	collector := NewMetricCollector()

	t.Run("カウンターメトリクス", func(t *testing.T) {
		collector.SetCounter("test_counter", 10, map[string]string{"label": "value"})
		
		metric, exists := collector.GetMetric("test_counter")
		require.True(t, exists)
		assert.Equal(t, "test_counter", metric.Name)
		assert.Equal(t, MetricTypeCounter, metric.Type)
		assert.Equal(t, 10.0, metric.Value)
		assert.Equal(t, "value", metric.Labels["label"])
	})

	t.Run("ゲージメトリクス", func(t *testing.T) {
		collector.SetGauge("test_gauge", 25.5, nil)
		
		metric, exists := collector.GetMetric("test_gauge")
		require.True(t, exists)
		assert.Equal(t, "test_gauge", metric.Name)
		assert.Equal(t, MetricTypeGauge, metric.Type)
		assert.Equal(t, 25.5, metric.Value)
	})

	t.Run("カウンターの増加", func(t *testing.T) {
		collector.IncrementCounter("test_counter", nil)
		
		metric, exists := collector.GetMetric("test_counter")
		require.True(t, exists)
		assert.Equal(t, 11.0, metric.Value)
	})

	t.Run("全てのメトリクス取得", func(t *testing.T) {
		metrics := collector.GetAllMetrics()
		assert.Len(t, metrics, 2)
		assert.Contains(t, metrics, "test_counter")
		assert.Contains(t, metrics, "test_gauge")
	})

	t.Run("メトリクスの削除", func(t *testing.T) {
		collector.RemoveMetric("test_counter")
		
		_, exists := collector.GetMetric("test_counter")
		assert.False(t, exists)
	})
}

func TestSystemMetricsCollector(t *testing.T) {
	collector := NewSystemMetricsCollector()
	collector.CollectSystemMetrics()

	metrics := collector.GetMetrics()

	// システムメトリクスが収集されていることを確認
	assert.Contains(t, metrics, "memory_alloc_bytes")
	assert.Contains(t, metrics, "goroutines_count")
	assert.Contains(t, metrics, "uptime_seconds")

	// メトリクスの型を確認
	memoryMetric := metrics["memory_alloc_bytes"]
	assert.Equal(t, MetricTypeGauge, memoryMetric.Type)
	assert.GreaterOrEqual(t, memoryMetric.Value, 0.0)
}

func TestApplicationMetricsCollector(t *testing.T) {
	collector := NewApplicationMetricsCollector()

	t.Run("リクエストの記録", func(t *testing.T) {
		collector.RecordRequest(100*time.Millisecond, 200)
		
		metrics := collector.GetMetrics()
		assert.Contains(t, metrics, "requests_total")
		assert.Contains(t, metrics, "http_requests_total")
		assert.Contains(t, metrics, "http_response_time_ms")
	})

	t.Run("エラーの記録", func(t *testing.T) {
		collector.RecordError("validation_error")
		
		metrics := collector.GetMetrics()
		assert.Contains(t, metrics, "errors_total")
	})

	t.Run("アクティブユーザーの設定", func(t *testing.T) {
		collector.SetActiveUsers(150)
		
		metrics := collector.GetMetrics()
		activeUsersMetric := metrics["active_users"]
		assert.Equal(t, 150.0, activeUsersMetric.Value)
	})
}

func TestDatabaseMetricsCollector(t *testing.T) {
	collector := NewDatabaseMetricsCollector()

	t.Run("クエリの記録", func(t *testing.T) {
		collector.RecordQuery("SELECT", 50*time.Millisecond, true)
		
		metrics := collector.GetMetrics()
		assert.Contains(t, metrics, "database_queries_total")
	})

	t.Run("エラークエリの記録", func(t *testing.T) {
		collector.RecordQuery("INSERT", 30*time.Millisecond, false)
		
		metrics := collector.GetMetrics()
		assert.Contains(t, metrics, "database_errors_total")
	})

	t.Run("接続数の設定", func(t *testing.T) {
		collector.SetConnectionCount(10)
		
		metrics := collector.GetMetrics()
		connectionsMetric := metrics["database_connections"]
		assert.Equal(t, 10.0, connectionsMetric.Value)
	})
}

func TestMetricsManager(t *testing.T) {
	manager := NewMetricsManager()

	t.Run("全てのメトリクス収集", func(t *testing.T) {
		metrics := manager.CollectAllMetrics()
		
		// プレフィックス付きのメトリクスが含まれていることを確認
		assert.Contains(t, metrics, "system_memory_alloc_bytes")
		assert.Contains(t, metrics, "app_requests_total")
		assert.Contains(t, metrics, "db_database_queries_total")
	})

	t.Run("各コレクターの取得", func(t *testing.T) {
		systemMetrics := manager.GetSystemMetrics()
		assert.NotNil(t, systemMetrics)

		appMetrics := manager.GetApplicationMetrics()
		assert.NotNil(t, appMetrics)

		dbMetrics := manager.GetDatabaseMetrics()
		assert.NotNil(t, dbMetrics)
	})
}

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// ヘルスチェッカーとメトリクス管理の作成
	healthChecker := NewHealthChecker()
	metricsManager := NewMetricsManager()

	// テスト用のヘルスチェックを登録
	testCheck := &mockHealthCheck{
		name:    "test_check",
		timeout: 5 * time.Second,
	}
	healthChecker.RegisterCheck(testCheck)

	// ハンドラーの作成
	handler := NewHealthHandler(healthChecker, metricsManager)

	t.Run("ヘルスチェックエンドポイント", func(t *testing.T) {
		router.GET("/health", handler.HealthCheck)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, StatusHealthy, response.Status)
		assert.NotZero(t, response.Timestamp)
		assert.NotZero(t, response.Duration)
		assert.Contains(t, response.Checks, "test_check")
	})

	t.Run("ライブネスチェックエンドポイント", func(t *testing.T) {
		router.GET("/live", handler.LivenessCheck)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/live", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, StatusHealthy, response.Status)
		assert.Contains(t, response.Details["message"], "稼働中")
	})

	t.Run("メトリクスエンドポイント", func(t *testing.T) {
		router.GET("/metrics", handler.Metrics)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/metrics", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response MetricsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotZero(t, response.Timestamp)
		assert.NotNil(t, response.Metrics)
	})
}

// モックヘルスチェック
type mockHealthCheck struct {
	name    string
	timeout time.Duration
	status  HealthStatus
}

func (m *mockHealthCheck) Name() string {
	return m.name
}

func (m *mockHealthCheck) Timeout() time.Duration {
	return m.timeout
}

func (m *mockHealthCheck) Check(ctx context.Context) HealthResult {
	status := m.status
	if status == "" {
		status = StatusHealthy
	}

	return HealthResult{
		Name:    m.name,
		Status:  status,
		Message: "モックチェック成功",
	}
}
