package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// CacheConfigManager キャッシュ設定管理
type CacheConfigManager struct {
	config *CacheConfig
}

// NewCacheConfigManager 新しいキャッシュ設定管理を作成
func NewCacheConfigManager() *CacheConfigManager {
	return &CacheConfigManager{
		config: DefaultCacheConfig(),
	}
}

// LoadFromEnv 環境変数から設定を読み込み
func (c *CacheConfigManager) LoadFromEnv() {
	if host := os.Getenv("REDIS_HOST"); host != "" {
		c.config.Host = host
	}

	if portStr := os.Getenv("REDIS_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			c.config.Port = port
		}
	}

	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		c.config.Password = password
	}

	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			c.config.DB = db
		}
	}

	if poolSizeStr := os.Getenv("REDIS_POOL_SIZE"); poolSizeStr != "" {
		if poolSize, err := strconv.Atoi(poolSizeStr); err == nil {
			c.config.PoolSize = poolSize
		}
	}

	if timeoutStr := os.Getenv("REDIS_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			c.config.Timeout = timeout
		}
	}
}

// LoadFromFile ファイルから設定を読み込み
func (c *CacheConfigManager) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("設定ファイルの読み込みに失敗しました: %v", err)
	}

	var config CacheConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("設定ファイルの解析に失敗しました: %v", err)
	}

	c.config = &config
	return nil
}

// SaveToFile 設定をファイルに保存
func (c *CacheConfigManager) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		return fmt.Errorf("設定のシリアライズに失敗しました: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("設定ファイルの保存に失敗しました: %v", err)
	}

	return nil
}

// GetConfig 設定を取得
func (c *CacheConfigManager) GetConfig() *CacheConfig {
	return c.config
}

// SetConfig 設定を設定
func (c *CacheConfigManager) SetConfig(config *CacheConfig) {
	c.config = config
}

// ValidateConfig 設定の妥当性を検証
func (c *CacheConfigManager) ValidateConfig() error {
	if c.config.Host == "" {
		return fmt.Errorf("ホストが設定されていません")
	}

	if c.config.Port <= 0 || c.config.Port > 65535 {
		return fmt.Errorf("無効なポート番号: %d", c.config.Port)
	}

	if c.config.DB < 0 {
		return fmt.Errorf("無効なデータベース番号: %d", c.config.DB)
	}

	if c.config.PoolSize <= 0 {
		return fmt.Errorf("無効なプールサイズ: %d", c.config.PoolSize)
	}

	if c.config.Timeout <= 0 {
		return fmt.Errorf("無効なタイムアウト: %v", c.config.Timeout)
	}

	return nil
}

// CacheKeyBuilder キャッシュキービルダー
type CacheKeyBuilder struct {
	prefix string
	parts  []string
}

// NewCacheKeyBuilder 新しいキャッシュキービルダーを作成
func NewCacheKeyBuilder(prefix string) *CacheKeyBuilder {
	return &CacheKeyBuilder{
		prefix: prefix,
		parts:  make([]string, 0),
	}
}

// Add キーの部分を追加
func (b *CacheKeyBuilder) Add(part string) *CacheKeyBuilder {
	b.parts = append(b.parts, part)
	return b
}

// AddInt 整数を追加
func (b *CacheKeyBuilder) AddInt(value int) *CacheKeyBuilder {
	b.parts = append(b.parts, strconv.Itoa(value))
	return b
}

// AddInt64 64ビット整数を追加
func (b *CacheKeyBuilder) AddInt64(value int64) *CacheKeyBuilder {
	b.parts = append(b.parts, strconv.FormatInt(value, 10))
	return b
}

// Build キーを構築
func (b *CacheKeyBuilder) Build() string {
	if len(b.parts) == 0 {
		return b.prefix
	}

	key := b.prefix + ":" + strings.Join(b.parts, ":")
	return key
}

// CacheStrategy キャッシュ戦略
type CacheStrategy struct {
	TTL           time.Duration
	RefreshTTL    time.Duration
	MaxRetries    int
	RetryInterval time.Duration
}

// DefaultCacheStrategy デフォルトキャッシュ戦略
func DefaultCacheStrategy() *CacheStrategy {
	return &CacheStrategy{
		TTL:           1 * time.Hour,
		RefreshTTL:   30 * time.Minute,
		MaxRetries:    3,
		RetryInterval: 100 * time.Millisecond,
	}
}

// CacheWithFallback フォールバック付きキャッシュ
type CacheWithFallback struct {
	cache    *CacheManager
	strategy *CacheStrategy
}

// NewCacheWithFallback 新しいフォールバック付きキャッシュを作成
func NewCacheWithFallback(cache *CacheManager, strategy *CacheStrategy) *CacheWithFallback {
	if strategy == nil {
		strategy = DefaultCacheStrategy()
	}

	return &CacheWithFallback{
		cache:    cache,
		strategy: strategy,
	}
}

// GetOrSet キャッシュから取得、なければ設定
func (c *CacheWithFallback) GetOrSet(ctx context.Context, key string, fallback func() (interface{}, error)) (string, error) {
	// キャッシュから取得を試行
	val, err := c.cache.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	// フォールバック関数を実行
	data, err := fallback()
	if err != nil {
		return "", fmt.Errorf("フォールバック関数の実行に失敗しました: %v", err)
	}

	// キャッシュに設定
	if err := c.cache.Set(ctx, key, data, c.strategy.TTL); err != nil {
		// キャッシュの設定に失敗しても、データは返す
		fmt.Printf("警告: キャッシュの設定に失敗しました: %v\n", err)
	}

	// データを文字列として返す
	switch v := data.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("データのシリアライズに失敗しました: %v", err)
		}
		return string(jsonData), nil
	}
}

// GetOrSetObject オブジェクトをキャッシュから取得、なければ設定
func (c *CacheWithFallback) GetOrSetObject(ctx context.Context, key string, dest interface{}, fallback func() (interface{}, error)) error {
	// キャッシュから取得を試行
	err := c.cache.GetObject(ctx, key, dest)
	if err == nil {
		return nil
	}

	// フォールバック関数を実行
	data, err := fallback()
	if err != nil {
		return fmt.Errorf("フォールバック関数の実行に失敗しました: %v", err)
	}

	// キャッシュに設定
	if err := c.cache.Set(ctx, key, data, c.strategy.TTL); err != nil {
		// キャッシュの設定に失敗しても、データは設定する
		fmt.Printf("警告: キャッシュの設定に失敗しました: %v\n", err)
	}

	// データをdestに設定
	switch v := data.(type) {
	case string:
		return json.Unmarshal([]byte(v), dest)
	case []byte:
		return json.Unmarshal(v, dest)
	default:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("データのシリアライズに失敗しました: %v", err)
		}
		return json.Unmarshal(jsonData, dest)
	}
}

// Refresh キャッシュを更新
func (c *CacheWithFallback) Refresh(ctx context.Context, key string, fallback func() (interface{}, error)) error {
	// フォールバック関数を実行
	data, err := fallback()
	if err != nil {
		return fmt.Errorf("フォールバック関数の実行に失敗しました: %v", err)
	}

	// キャッシュを更新
	return c.cache.Set(ctx, key, data, c.strategy.TTL)
}

// Invalidate キャッシュを無効化
func (c *CacheWithFallback) Invalidate(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, key)
}

// CachePattern キャッシュパターン
type CachePattern struct {
	KeyPattern string
	TTL        time.Duration
	Strategy   *CacheStrategy
}

// CachePatternManager キャッシュパターン管理
type CachePatternManager struct {
	patterns map[string]*CachePattern
	cache    *CacheManager
}

// NewCachePatternManager 新しいキャッシュパターン管理を作成
func NewCachePatternManager(cache *CacheManager) *CachePatternManager {
	return &CachePatternManager{
		patterns: make(map[string]*CachePattern),
		cache:    cache,
	}
}

// RegisterPattern パターンを登録
func (m *CachePatternManager) RegisterPattern(name string, pattern *CachePattern) {
	m.patterns[name] = pattern
}

// GetPattern パターンを取得
func (m *CachePatternManager) GetPattern(name string) (*CachePattern, bool) {
	pattern, exists := m.patterns[name]
	return pattern, exists
}

// InvalidatePattern パターンに一致するキーを無効化
func (m *CachePatternManager) InvalidatePattern(ctx context.Context, patternName string) error {
	pattern, exists := m.GetPattern(patternName)
	if !exists {
		return fmt.Errorf("パターン '%s' が見つかりません", patternName)
	}

	keys, err := m.cache.Keys(ctx, pattern.KeyPattern)
	if err != nil {
		return fmt.Errorf("キーの検索に失敗しました: %v", err)
	}

	for _, key := range keys {
		if err := m.cache.Delete(ctx, key); err != nil {
			fmt.Printf("警告: キー '%s' の削除に失敗しました: %v\n", key, err)
		}
	}

	return nil
}

// CacheMetrics キャッシュメトリクス
type CacheMetrics struct {
	Hits       int64
	Misses     int64
	Sets       int64
	Deletes    int64
	Errors     int64
	TotalSize  int64
	KeyCount   int64
}

// CacheMetricsCollector キャッシュメトリクスコレクター
type CacheMetricsCollector struct {
	metrics CacheMetrics
	cache   *CacheManager
}

// NewCacheMetricsCollector 新しいキャッシュメトリクスコレクターを作成
func NewCacheMetricsCollector(cache *CacheManager) *CacheMetricsCollector {
	return &CacheMetricsCollector{
		metrics: CacheMetrics{},
		cache:   cache,
	}
}

// RecordHit ヒットを記録
func (c *CacheMetricsCollector) RecordHit() {
	c.metrics.Hits++
}

// RecordMiss ミスを記録
func (c *CacheMetricsCollector) RecordMiss() {
	c.metrics.Misses++
}

// RecordSet セットを記録
func (c *CacheMetricsCollector) RecordSet() {
	c.metrics.Sets++
}

// RecordDelete 削除を記録
func (c *CacheMetricsCollector) RecordDelete() {
	c.metrics.Deletes++
}

// RecordError エラーを記録
func (c *CacheMetricsCollector) RecordError() {
	c.metrics.Errors++
}

// GetMetrics メトリクスを取得
func (c *CacheMetricsCollector) GetMetrics() CacheMetrics {
	return c.metrics
}

// GetHitRate ヒット率を取得
func (c *CacheMetricsCollector) GetHitRate() float64 {
	total := c.metrics.Hits + c.metrics.Misses
	if total == 0 {
		return 0.0
	}
	return float64(c.metrics.Hits) / float64(total)
}

// ResetMetrics メトリクスをリセット
func (c *CacheMetricsCollector) ResetMetrics() {
	c.metrics = CacheMetrics{}
}
