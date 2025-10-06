package health

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"tea-logistics/pkg/logger"
)

// MetricType メトリクスタイプ
type MetricType string

const (
	MetricTypeCounter MetricType = "counter"
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
)

// Metric メトリクス
type Metric struct {
	Name      string                 `json:"name"`
	Type      MetricType             `json:"type"`
	Value     float64                `json:"value"`
	Labels    map[string]string      `json:"labels,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// MetricCollector メトリクスコレクター
type MetricCollector struct {
	metrics map[string]*Metric
	mu      sync.RWMutex
}

// NewMetricCollector 新しいメトリクスコレクターを作成
func NewMetricCollector() *MetricCollector {
	return &MetricCollector{
		metrics: make(map[string]*Metric),
	}
}

// SetCounter カウンターメトリクスを設定
func (m *MetricCollector) SetCounter(name string, value float64, labels map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.metrics[name] = &Metric{
		Name:      name,
		Type:      MetricTypeCounter,
		Value:     value,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}

// SetGauge ゲージメトリクスを設定
func (m *MetricCollector) SetGauge(name string, value float64, labels map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.metrics[name] = &Metric{
		Name:      name,
		Type:      MetricTypeGauge,
		Value:     value,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}

// IncrementCounter カウンターメトリクスを増加
func (m *MetricCollector) IncrementCounter(name string, labels map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if metric, exists := m.metrics[name]; exists {
		metric.Value++
		metric.Timestamp = time.Now()
	} else {
		m.metrics[name] = &Metric{
			Name:      name,
			Type:      MetricTypeCounter,
			Value:     1,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// GetMetric メトリクスを取得
func (m *MetricCollector) GetMetric(name string) (*Metric, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	metric, exists := m.metrics[name]
	if exists {
		// コピーを返す
		metricCopy := *metric
		return &metricCopy, true
	}
	return nil, false
}

// GetAllMetrics 全てのメトリクスを取得
func (m *MetricCollector) GetAllMetrics() map[string]*Metric {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	metrics := make(map[string]*Metric)
	for name, metric := range m.metrics {
		// コピーを作成
		metricCopy := *metric
		metrics[name] = &metricCopy
	}
	return metrics
}

// RemoveMetric メトリクスを削除
func (m *MetricCollector) RemoveMetric(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.metrics, name)
}

// ClearMetrics 全てのメトリクスをクリア
func (m *MetricCollector) ClearMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics = make(map[string]*Metric)
}

// SystemMetricsCollector システムメトリクスコレクター
type SystemMetricsCollector struct {
	collector *MetricCollector
	startTime time.Time
}

// NewSystemMetricsCollector 新しいシステムメトリクスコレクターを作成
func NewSystemMetricsCollector() *SystemMetricsCollector {
	return &SystemMetricsCollector{
		collector: NewMetricCollector(),
		startTime: time.Now(),
	}
}

// CollectSystemMetrics システムメトリクスを収集
func (s *SystemMetricsCollector) CollectSystemMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// メモリメトリクス
	s.collector.SetGauge("memory_alloc_bytes", float64(m.Alloc), nil)
	s.collector.SetGauge("memory_total_alloc_bytes", float64(m.TotalAlloc), nil)
	s.collector.SetGauge("memory_sys_bytes", float64(m.Sys), nil)
	s.collector.SetGauge("memory_num_gc", float64(m.NumGC), nil)
	s.collector.SetGauge("memory_gc_cpu_fraction", m.GCCPUFraction, nil)

	// ゴルーチンメトリクス
	s.collector.SetGauge("goroutines_count", float64(runtime.NumGoroutine()), nil)

	// 実行時間メトリクス
	uptime := time.Since(s.startTime)
	s.collector.SetGauge("uptime_seconds", uptime.Seconds(), nil)

	// CPU使用率（簡易版）
	s.collector.SetGauge("cpu_usage_percent", s.getCPUUsage(), nil)
}

// getCPUUsage CPU使用率を取得（簡易版）
func (s *SystemMetricsCollector) getCPUUsage() float64 {
	// 簡易的なCPU使用率計算
	// 実際の本番環境では、より精密な計算が必要
	return 0.0 // プレースホルダー
}

// GetMetrics メトリクスを取得
func (s *SystemMetricsCollector) GetMetrics() map[string]*Metric {
	return s.collector.GetAllMetrics()
}

// ApplicationMetricsCollector アプリケーションメトリクスコレクター
type ApplicationMetricsCollector struct {
	collector     *MetricCollector
	requestCount  int64
	errorCount    int64
	responseTime  int64 // ミリ秒
	activeUsers   int64
}

// NewApplicationMetricsCollector 新しいアプリケーションメトリクスコレクターを作成
func NewApplicationMetricsCollector() *ApplicationMetricsCollector {
	return &ApplicationMetricsCollector{
		collector: NewMetricCollector(),
	}
}

// RecordRequest リクエストを記録
func (a *ApplicationMetricsCollector) RecordRequest(responseTime time.Duration, statusCode int) {
	atomic.AddInt64(&a.requestCount, 1)
	atomic.AddInt64(&a.responseTime, int64(responseTime.Milliseconds()))

	labels := map[string]string{
		"status_code": fmt.Sprintf("%d", statusCode),
	}
	a.collector.IncrementCounter("http_requests_total", labels)

	// レスポンス時間メトリクス
	a.collector.SetGauge("http_response_time_ms", float64(responseTime.Milliseconds()), labels)
}

// RecordError エラーを記録
func (a *ApplicationMetricsCollector) RecordError(errorType string) {
	atomic.AddInt64(&a.errorCount, 1)

	labels := map[string]string{
		"error_type": errorType,
	}
	a.collector.IncrementCounter("errors_total", labels)
}

// SetActiveUsers アクティブユーザー数を設定
func (a *ApplicationMetricsCollector) SetActiveUsers(count int64) {
	atomic.StoreInt64(&a.activeUsers, count)
	a.collector.SetGauge("active_users", float64(count), nil)
}

// GetMetrics メトリクスを取得
func (a *ApplicationMetricsCollector) GetMetrics() map[string]*Metric {
	// アプリケーションメトリクスを更新
	a.collector.SetGauge("requests_total", float64(atomic.LoadInt64(&a.requestCount)), nil)
	a.collector.SetGauge("errors_total", float64(atomic.LoadInt64(&a.errorCount)), nil)
	a.collector.SetGauge("response_time_avg_ms", float64(atomic.LoadInt64(&a.responseTime)), nil)
	a.collector.SetGauge("active_users", float64(atomic.LoadInt64(&a.activeUsers)), nil)

	return a.collector.GetAllMetrics()
}

// DatabaseMetricsCollector データベースメトリクスコレクター
type DatabaseMetricsCollector struct {
	collector      *MetricCollector
	queryCount     int64
	queryTime      int64 // ミリ秒
	connectionCount int64
	errorCount     int64
}

// NewDatabaseMetricsCollector 新しいデータベースメトリクスコレクターを作成
func NewDatabaseMetricsCollector() *DatabaseMetricsCollector {
	return &DatabaseMetricsCollector{
		collector: NewMetricCollector(),
	}
}

// RecordQuery クエリを記録
func (d *DatabaseMetricsCollector) RecordQuery(queryType string, duration time.Duration, success bool) {
	atomic.AddInt64(&d.queryCount, 1)
	atomic.AddInt64(&d.queryTime, int64(duration.Milliseconds()))

	labels := map[string]string{
		"query_type": queryType,
		"success":    fmt.Sprintf("%t", success),
	}
	d.collector.IncrementCounter("database_queries_total", labels)

	if !success {
		atomic.AddInt64(&d.errorCount, 1)
		d.collector.IncrementCounter("database_errors_total", labels)
	}
}

// SetConnectionCount 接続数を設定
func (d *DatabaseMetricsCollector) SetConnectionCount(count int64) {
	atomic.StoreInt64(&d.connectionCount, count)
	d.collector.SetGauge("database_connections", float64(count), nil)
}

// GetMetrics メトリクスを取得
func (d *DatabaseMetricsCollector) GetMetrics() map[string]*Metric {
	// データベースメトリクスを更新
	d.collector.SetGauge("database_queries_total", float64(atomic.LoadInt64(&d.queryCount)), nil)
	d.collector.SetGauge("database_query_time_avg_ms", float64(atomic.LoadInt64(&d.queryTime)), nil)
	d.collector.SetGauge("database_connections", float64(atomic.LoadInt64(&d.connectionCount)), nil)
	d.collector.SetGauge("database_errors_total", float64(atomic.LoadInt64(&d.errorCount)), nil)

	return d.collector.GetAllMetrics()
}

// MetricsManager メトリクス管理
type MetricsManager struct {
	systemMetrics      *SystemMetricsCollector
	applicationMetrics *ApplicationMetricsCollector
	databaseMetrics    *DatabaseMetricsCollector
	mu                 sync.RWMutex
}

// NewMetricsManager 新しいメトリクス管理を作成
func NewMetricsManager() *MetricsManager {
	return &MetricsManager{
		systemMetrics:      NewSystemMetricsCollector(),
		applicationMetrics: NewApplicationMetricsCollector(),
		databaseMetrics:    NewDatabaseMetricsCollector(),
	}
}

// CollectAllMetrics 全てのメトリクスを収集
func (m *MetricsManager) CollectAllMetrics() map[string]*Metric {
	m.mu.Lock()
	defer m.mu.Unlock()

	// システムメトリクスを収集
	m.systemMetrics.CollectSystemMetrics()

	allMetrics := make(map[string]*Metric)

	// システムメトリクスを追加
	for name, metric := range m.systemMetrics.GetMetrics() {
		allMetrics["system_"+name] = metric
	}

	// アプリケーションメトリクスを追加
	for name, metric := range m.applicationMetrics.GetMetrics() {
		allMetrics["app_"+name] = metric
	}

	// データベースメトリクスを追加
	for name, metric := range m.databaseMetrics.GetMetrics() {
		allMetrics["db_"+name] = metric
	}

	return allMetrics
}

// GetSystemMetrics システムメトリクスを取得
func (m *MetricsManager) GetSystemMetrics() *SystemMetricsCollector {
	return m.systemMetrics
}

// GetApplicationMetrics アプリケーションメトリクスを取得
func (m *MetricsManager) GetApplicationMetrics() *ApplicationMetricsCollector {
	return m.applicationMetrics
}

// GetDatabaseMetrics データベースメトリクスを取得
func (m *MetricsManager) GetDatabaseMetrics() *DatabaseMetricsCollector {
	return m.databaseMetrics
}

// StartMetricsCollection メトリクス収集を開始
func (m *MetricsManager) StartMetricsCollection(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("メトリクス収集を停止しました")
			return
		case <-ticker.C:
			metrics := m.CollectAllMetrics()
			logger.Debug("メトリクス収集完了", map[string]interface{}{
				"metric_count": len(metrics),
			})
		}
	}
}

// グローバルメトリクス管理
var globalMetricsManager *MetricsManager

// InitGlobalMetricsManager グローバルメトリクス管理を初期化
func InitGlobalMetricsManager() {
	globalMetricsManager = NewMetricsManager()
}

// GetGlobalMetricsManager グローバルメトリクス管理を取得
func GetGlobalMetricsManager() *MetricsManager {
	if globalMetricsManager == nil {
		globalMetricsManager = NewMetricsManager()
	}
	return globalMetricsManager
}

// RecordGlobalRequest グローバルリクエストを記録
func RecordGlobalRequest(responseTime time.Duration, statusCode int) {
	GetGlobalMetricsManager().GetApplicationMetrics().RecordRequest(responseTime, statusCode)
}

// RecordGlobalError グローバルエラーを記録
func RecordGlobalError(errorType string) {
	GetGlobalMetricsManager().GetApplicationMetrics().RecordError(errorType)
}

// RecordGlobalQuery グローバルクエリを記録
func RecordGlobalQuery(queryType string, duration time.Duration, success bool) {
	GetGlobalMetricsManager().GetDatabaseMetrics().RecordQuery(queryType, duration, success)
}

// GetGlobalMetrics グローバルメトリクスを取得
func GetGlobalMetrics() map[string]*Metric {
	return GetGlobalMetricsManager().CollectAllMetrics()
}
