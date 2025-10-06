package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"tea-logistics/pkg/cache"
	"tea-logistics/pkg/logger"
)

// RateLimitStrategy レート制限戦略
type RateLimitStrategy string

const (
	StrategyFixedWindow   RateLimitStrategy = "fixed_window"
	StrategySlidingWindow RateLimitStrategy = "sliding_window"
	StrategyTokenBucket   RateLimitStrategy = "token_bucket"
	StrategyLeakyBucket   RateLimitStrategy = "leaky_bucket"
)

// RateLimitConfig レート制限設定
type RateLimitConfig struct {
	Strategy    RateLimitStrategy `json:"strategy"`
	Limit       int               `json:"limit"`
	Window      time.Duration     `json:"window"`
	Burst       int               `json:"burst"`
	RefillRate  int               `json:"refill_rate"`
	RefillTime  time.Duration     `json:"refill_time"`
	KeyPrefix   string            `json:"key_prefix"`
	SkipOnError bool              `json:"skip_on_error"`
}

// DefaultRateLimitConfig デフォルトレート制限設定
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		Strategy:    StrategyFixedWindow,
		Limit:       100,
		Window:      1 * time.Minute,
		Burst:       10,
		RefillRate:  10,
		RefillTime:  1 * time.Second,
		KeyPrefix:   "rate_limit",
		SkipOnError: true,
	}
}

// RateLimitResult レート制限結果
type RateLimitResult struct {
	Allowed    bool          `json:"allowed"`
	Limit      int           `json:"limit"`
	Remaining  int           `json:"remaining"`
	ResetTime  time.Time     `json:"reset_time"`
	RetryAfter time.Duration `json:"retry_after"`
	Strategy   string        `json:"strategy"`
	Key        string        `json:"key"`
}

// RateLimiter レート制限インターフェース
type RateLimiter interface {
	Allow(ctx context.Context, key string) (*RateLimitResult, error)
	Reset(ctx context.Context, key string) error
	GetLimit(ctx context.Context, key string) (*RateLimitResult, error)
}

// FixedWindowRateLimiter 固定ウィンドウレート制限
type FixedWindowRateLimiter struct {
	config *RateLimitConfig
	cache  *cache.CacheManager
}

// NewFixedWindowRateLimiter 新しい固定ウィンドウレート制限を作成
func NewFixedWindowRateLimiter(config *RateLimitConfig, cache *cache.CacheManager) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		config: config,
		cache:  cache,
	}
}

// Allow リクエストを許可するかチェック
func (r *FixedWindowRateLimiter) Allow(ctx context.Context, key string) (*RateLimitResult, error) {
	cacheKey := r.buildCacheKey(key)

	// 現在のウィンドウの開始時間を計算
	now := time.Now()
	windowStart := now.Truncate(r.config.Window)

	// ウィンドウキーを作成
	windowKey := fmt.Sprintf("%s:window:%d", cacheKey, windowStart.Unix())

	// 現在のカウントを取得
	count, err := r.getCurrentCount(ctx, windowKey)
	if err != nil {
		if r.config.SkipOnError {
			logger.Warn("レート制限エラー（スキップ）", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
			return &RateLimitResult{
				Allowed:   true,
				Limit:     r.config.Limit,
				Remaining: r.config.Limit,
				ResetTime: windowStart.Add(r.config.Window),
				Strategy:  string(r.config.Strategy),
				Key:       key,
			}, nil
		}
		return nil, err
	}

	// 制限をチェック
	if count >= r.config.Limit {
		resetTime := windowStart.Add(r.config.Window)
		retryAfter := resetTime.Sub(now)

		logger.Warn("レート制限に達しました", map[string]interface{}{
			"key":         key,
			"count":       count,
			"limit":       r.config.Limit,
			"reset_time":  resetTime,
			"retry_after": retryAfter.String(),
		})

		return &RateLimitResult{
			Allowed:    false,
			Limit:      r.config.Limit,
			Remaining:  0,
			ResetTime:  resetTime,
			RetryAfter: retryAfter,
			Strategy:   string(r.config.Strategy),
			Key:        key,
		}, nil
	}

	// カウントを増加
	newCount, err := r.incrementCount(ctx, windowKey)
	if err != nil {
		if r.config.SkipOnError {
			logger.Warn("レート制限カウント増加エラー（スキップ）", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
			return &RateLimitResult{
				Allowed:   true,
				Limit:     r.config.Limit,
				Remaining: r.config.Limit - count,
				ResetTime: windowStart.Add(r.config.Window),
				Strategy:  string(r.config.Strategy),
				Key:       key,
			}, nil
		}
		return nil, err
	}

	remaining := r.config.Limit - newCount
	if remaining < 0 {
		remaining = 0
	}

	logger.Debug("レート制限チェック成功", map[string]interface{}{
		"key":       key,
		"count":     newCount,
		"limit":     r.config.Limit,
		"remaining": remaining,
	})

	return &RateLimitResult{
		Allowed:   true,
		Limit:     r.config.Limit,
		Remaining: remaining,
		ResetTime: windowStart.Add(r.config.Window),
		Strategy:  string(r.config.Strategy),
		Key:       key,
	}, nil
}

// Reset レート制限をリセット
func (r *FixedWindowRateLimiter) Reset(ctx context.Context, key string) error {
	cacheKey := r.buildCacheKey(key)
	now := time.Now()
	windowStart := now.Truncate(r.config.Window)
	windowKey := fmt.Sprintf("%s:window:%d", cacheKey, windowStart.Unix())

	err := r.cache.Delete(ctx, windowKey)
	if err != nil {
		logger.Error("レート制限リセットエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return err
	}

	logger.Info("レート制限をリセットしました", map[string]interface{}{
		"key": key,
	})

	return nil
}

// GetLimit 現在の制限情報を取得
func (r *FixedWindowRateLimiter) GetLimit(ctx context.Context, key string) (*RateLimitResult, error) {
	cacheKey := r.buildCacheKey(key)
	now := time.Now()
	windowStart := now.Truncate(r.config.Window)
	windowKey := fmt.Sprintf("%s:window:%d", cacheKey, windowStart.Unix())

	count, err := r.getCurrentCount(ctx, windowKey)
	if err != nil {
		if r.config.SkipOnError {
			count = 0
		} else {
			return nil, err
		}
	}

	remaining := r.config.Limit - count
	if remaining < 0 {
		remaining = 0
	}

	return &RateLimitResult{
		Allowed:   count < r.config.Limit,
		Limit:     r.config.Limit,
		Remaining: remaining,
		ResetTime: windowStart.Add(r.config.Window),
		Strategy:  string(r.config.Strategy),
		Key:       key,
	}, nil
}

// buildCacheKey キャッシュキーを構築
func (r *FixedWindowRateLimiter) buildCacheKey(key string) string {
	return fmt.Sprintf("%s:%s", r.config.KeyPrefix, key)
}

// getCurrentCount 現在のカウントを取得
func (r *FixedWindowRateLimiter) getCurrentCount(ctx context.Context, windowKey string) (int, error) {
	val, err := r.cache.Get(ctx, windowKey)
	if err != nil {
		// キーが存在しない場合は0を返す
		if err.Error() == fmt.Sprintf("キー '%s' が見つかりません", windowKey) {
			return 0, nil
		}
		return 0, err
	}

	// 文字列を整数に変換
	var count int
	_, err = fmt.Sscanf(val, "%d", &count)
	if err != nil {
		return 0, fmt.Errorf("カウント値の解析に失敗しました: %v", err)
	}

	return count, nil
}

// incrementCount カウントを増加
func (r *FixedWindowRateLimiter) incrementCount(ctx context.Context, windowKey string) (int, error) {
	// ウィンドウの有効期限を設定
	expiration := r.config.Window

	// カウントを増加
	count, err := r.cache.Increment(ctx, windowKey)
	if err != nil {
		// キーが存在しない場合は新しく作成
		if err.Error() == fmt.Sprintf("キー '%s' が見つかりません", windowKey) {
			err = r.cache.Set(ctx, windowKey, "1", expiration)
			if err != nil {
				return 0, err
			}
			return 1, nil
		}
		return 0, err
	}

	// 有効期限を更新
	err = r.cache.Expire(ctx, windowKey, expiration)
	if err != nil {
		logger.Warn("有効期限の更新に失敗しました", map[string]interface{}{
			"key":        windowKey,
			"expiration": expiration.String(),
			"error":      err.Error(),
		})
	}

	return int(count), nil
}

// SlidingWindowRateLimiter スライディングウィンドウレート制限
type SlidingWindowRateLimiter struct {
	config *RateLimitConfig
	cache  *cache.CacheManager
}

// NewSlidingWindowRateLimiter 新しいスライディングウィンドウレート制限を作成
func NewSlidingWindowRateLimiter(config *RateLimitConfig, cache *cache.CacheManager) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		config: config,
		cache:  cache,
	}
}

// Allow リクエストを許可するかチェック
func (r *SlidingWindowRateLimiter) Allow(ctx context.Context, key string) (*RateLimitResult, error) {
	cacheKey := r.buildCacheKey(key)
	now := time.Now()

	// スライディングウィンドウの開始時間を計算
	windowStart := now.Add(-r.config.Window)

	// 古いエントリを削除
	err := r.cleanOldEntries(ctx, cacheKey, windowStart)
	if err != nil {
		if r.config.SkipOnError {
			logger.Warn("古いエントリの削除エラー（スキップ）", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
		} else {
			return nil, err
		}
	}

	// 現在のカウントを取得
	count, err := r.getCurrentCount(ctx, cacheKey)
	if err != nil {
		if r.config.SkipOnError {
			logger.Warn("レート制限カウント取得エラー（スキップ）", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
			return &RateLimitResult{
				Allowed:   true,
				Limit:     r.config.Limit,
				Remaining: r.config.Limit,
				ResetTime: now.Add(r.config.Window),
				Strategy:  string(r.config.Strategy),
				Key:       key,
			}, nil
		}
		return nil, err
	}

	// 制限をチェック
	if count >= r.config.Limit {
		retryAfter := r.config.Window

		logger.Warn("スライディングウィンドウレート制限に達しました", map[string]interface{}{
			"key":         key,
			"count":       count,
			"limit":       r.config.Limit,
			"retry_after": retryAfter.String(),
		})

		return &RateLimitResult{
			Allowed:    false,
			Limit:      r.config.Limit,
			Remaining:  0,
			ResetTime:  now.Add(r.config.Window),
			RetryAfter: retryAfter,
			Strategy:   string(r.config.Strategy),
			Key:        key,
		}, nil
	}

	// 新しいエントリを追加
	err = r.addEntry(ctx, cacheKey, now)
	if err != nil {
		if r.config.SkipOnError {
			logger.Warn("エントリ追加エラー（スキップ）", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
			return &RateLimitResult{
				Allowed:   true,
				Limit:     r.config.Limit,
				Remaining: r.config.Limit - count,
				ResetTime: now.Add(r.config.Window),
				Strategy:  string(r.config.Strategy),
				Key:       key,
			}, nil
		}
		return nil, err
	}

	remaining := r.config.Limit - count - 1
	if remaining < 0 {
		remaining = 0
	}

	logger.Debug("スライディングウィンドウレート制限チェック成功", map[string]interface{}{
		"key":       key,
		"count":     count + 1,
		"limit":     r.config.Limit,
		"remaining": remaining,
	})

	return &RateLimitResult{
		Allowed:   true,
		Limit:     r.config.Limit,
		Remaining: remaining,
		ResetTime: now.Add(r.config.Window),
		Strategy:  string(r.config.Strategy),
		Key:       key,
	}, nil
}

// Reset レート制限をリセット
func (r *SlidingWindowRateLimiter) Reset(ctx context.Context, key string) error {
	cacheKey := r.buildCacheKey(key)

	// パターンに一致するキーを削除
	pattern := fmt.Sprintf("%s:*", cacheKey)
	keys, err := r.cache.Keys(ctx, pattern)
	if err != nil {
		logger.Error("レート制限リセットエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return err
	}

	for _, k := range keys {
		err = r.cache.Delete(ctx, k)
		if err != nil {
			logger.Warn("キーの削除に失敗しました", map[string]interface{}{
				"key":   k,
				"error": err.Error(),
			})
		}
	}

	logger.Info("スライディングウィンドウレート制限をリセットしました", map[string]interface{}{
		"key": key,
	})

	return nil
}

// GetLimit 現在の制限情報を取得
func (r *SlidingWindowRateLimiter) GetLimit(ctx context.Context, key string) (*RateLimitResult, error) {
	cacheKey := r.buildCacheKey(key)
	now := time.Now()
	windowStart := now.Add(-r.config.Window)

	// 古いエントリを削除
	err := r.cleanOldEntries(ctx, cacheKey, windowStart)
	if err != nil && !r.config.SkipOnError {
		return nil, err
	}

	count, err := r.getCurrentCount(ctx, cacheKey)
	if err != nil {
		if r.config.SkipOnError {
			count = 0
		} else {
			return nil, err
		}
	}

	remaining := r.config.Limit - count
	if remaining < 0 {
		remaining = 0
	}

	return &RateLimitResult{
		Allowed:   count < r.config.Limit,
		Limit:     r.config.Limit,
		Remaining: remaining,
		ResetTime: now.Add(r.config.Window),
		Strategy:  string(r.config.Strategy),
		Key:       key,
	}, nil
}

// buildCacheKey キャッシュキーを構築
func (r *SlidingWindowRateLimiter) buildCacheKey(key string) string {
	return fmt.Sprintf("%s:%s", r.config.KeyPrefix, key)
}

// cleanOldEntries 古いエントリを削除
func (r *SlidingWindowRateLimiter) cleanOldEntries(ctx context.Context, cacheKey string, windowStart time.Time) error {
	// 古いエントリのキーを生成
	oldKey := fmt.Sprintf("%s:entry:%d", cacheKey, windowStart.Unix())
	return r.cache.Delete(ctx, oldKey)
}

// getCurrentCount 現在のカウントを取得
func (r *SlidingWindowRateLimiter) getCurrentCount(ctx context.Context, cacheKey string) (int, error) {
	// パターンに一致するキーを検索
	pattern := fmt.Sprintf("%s:entry:*", cacheKey)
	keys, err := r.cache.Keys(ctx, pattern)
	if err != nil {
		return 0, err
	}

	return len(keys), nil
}

// addEntry 新しいエントリを追加
func (r *SlidingWindowRateLimiter) addEntry(ctx context.Context, cacheKey string, timestamp time.Time) error {
	entryKey := fmt.Sprintf("%s:entry:%d", cacheKey, timestamp.Unix())
	expiration := r.config.Window

	return r.cache.Set(ctx, entryKey, "1", expiration)
}

// TokenBucketRateLimiter トークンバケットレート制限
type TokenBucketRateLimiter struct {
	config *RateLimitConfig
	cache  *cache.CacheManager
	mu     sync.Mutex
}

// NewTokenBucketRateLimiter 新しいトークンバケットレート制限を作成
func NewTokenBucketRateLimiter(config *RateLimitConfig, cache *cache.CacheManager) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		config: config,
		cache:  cache,
	}
}

// Allow リクエストを許可するかチェック
func (r *TokenBucketRateLimiter) Allow(ctx context.Context, key string) (*RateLimitResult, error) {
	cacheKey := r.buildCacheKey(key)
	now := time.Now()

	// トークンバケットの状態を取得
	bucket, err := r.getBucketState(ctx, cacheKey)
	if err != nil {
		if r.config.SkipOnError {
			logger.Warn("トークンバケット状態取得エラー（スキップ）", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
			return &RateLimitResult{
				Allowed:   true,
				Limit:     r.config.Burst,
				Remaining: r.config.Burst,
				ResetTime: now.Add(r.config.RefillTime),
				Strategy:  string(r.config.Strategy),
				Key:       key,
			}, nil
		}
		return nil, err
	}

	// トークンを補充
	bucket = r.refillTokens(bucket, now)

	// トークンが利用可能かチェック
	if bucket.Tokens <= 0 {
		retryAfter := r.config.RefillTime

		logger.Warn("トークンバケットが空です", map[string]interface{}{
			"key":         key,
			"tokens":      bucket.Tokens,
			"burst":       r.config.Burst,
			"retry_after": retryAfter.String(),
		})

		return &RateLimitResult{
			Allowed:    false,
			Limit:      r.config.Burst,
			Remaining:  0,
			ResetTime:  bucket.LastRefill.Add(r.config.RefillTime),
			RetryAfter: retryAfter,
			Strategy:   string(r.config.Strategy),
			Key:        key,
		}, nil
	}

	// トークンを消費
	bucket.Tokens--
	bucket.LastRefill = now

	// バケット状態を保存
	err = r.saveBucketState(ctx, cacheKey, bucket)
	if err != nil {
		if r.config.SkipOnError {
			logger.Warn("トークンバケット状態保存エラー（スキップ）", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
		} else {
			return nil, err
		}
	}

	logger.Debug("トークンバケットチェック成功", map[string]interface{}{
		"key":       key,
		"tokens":    bucket.Tokens,
		"burst":     r.config.Burst,
		"remaining": bucket.Tokens,
	})

	return &RateLimitResult{
		Allowed:   true,
		Limit:     r.config.Burst,
		Remaining: bucket.Tokens,
		ResetTime: bucket.LastRefill.Add(r.config.RefillTime),
		Strategy:  string(r.config.Strategy),
		Key:       key,
	}, nil
}

// Reset レート制限をリセット
func (r *TokenBucketRateLimiter) Reset(ctx context.Context, key string) error {
	cacheKey := r.buildCacheKey(key)

	err := r.cache.Delete(ctx, cacheKey)
	if err != nil {
		logger.Error("トークンバケットリセットエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return err
	}

	logger.Info("トークンバケットをリセットしました", map[string]interface{}{
		"key": key,
	})

	return nil
}

// GetLimit 現在の制限情報を取得
func (r *TokenBucketRateLimiter) GetLimit(ctx context.Context, key string) (*RateLimitResult, error) {
	cacheKey := r.buildCacheKey(key)
	now := time.Now()

	bucket, err := r.getBucketState(ctx, cacheKey)
	if err != nil {
		if r.config.SkipOnError {
			bucket = &TokenBucket{
				Tokens:     r.config.Burst,
				LastRefill: now,
			}
		} else {
			return nil, err
		}
	}

	bucket = r.refillTokens(bucket, now)

	return &RateLimitResult{
		Allowed:   bucket.Tokens > 0,
		Limit:     r.config.Burst,
		Remaining: bucket.Tokens,
		ResetTime: bucket.LastRefill.Add(r.config.RefillTime),
		Strategy:  string(r.config.Strategy),
		Key:       key,
	}, nil
}

// TokenBucket トークンバケットの状態
type TokenBucket struct {
	Tokens     int       `json:"tokens"`
	LastRefill time.Time `json:"last_refill"`
}

// buildCacheKey キャッシュキーを構築
func (r *TokenBucketRateLimiter) buildCacheKey(key string) string {
	return fmt.Sprintf("%s:%s", r.config.KeyPrefix, key)
}

// getBucketState バケット状態を取得
func (r *TokenBucketRateLimiter) getBucketState(ctx context.Context, cacheKey string) (*TokenBucket, error) {
	var bucket TokenBucket
	err := r.cache.GetObject(ctx, cacheKey, &bucket)
	if err != nil {
		// キーが存在しない場合は新しいバケットを作成
		if err.Error() == fmt.Sprintf("キー '%s' が見つかりません", cacheKey) {
			return &TokenBucket{
				Tokens:     r.config.Burst,
				LastRefill: time.Now(),
			}, nil
		}
		return nil, err
	}

	return &bucket, nil
}

// saveBucketState バケット状態を保存
func (r *TokenBucketRateLimiter) saveBucketState(ctx context.Context, cacheKey string, bucket *TokenBucket) error {
	expiration := r.config.RefillTime * 2 // バケットの有効期限を設定
	return r.cache.Set(ctx, cacheKey, bucket, expiration)
}

// refillTokens トークンを補充
func (r *TokenBucketRateLimiter) refillTokens(bucket *TokenBucket, now time.Time) *TokenBucket {
	elapsed := now.Sub(bucket.LastRefill)
	refillCount := int(elapsed / r.config.RefillTime)

	if refillCount > 0 {
		bucket.Tokens += refillCount * r.config.RefillRate
		if bucket.Tokens > r.config.Burst {
			bucket.Tokens = r.config.Burst
		}
		bucket.LastRefill = bucket.LastRefill.Add(time.Duration(refillCount) * r.config.RefillTime)
	}

	return bucket
}

// RateLimitManager レート制限管理
type RateLimitManager struct {
	limiters map[RateLimitStrategy]RateLimiter
	cache    *cache.CacheManager
}

// NewRateLimitManager 新しいレート制限管理を作成
func NewRateLimitManager(cache *cache.CacheManager) *RateLimitManager {
	return &RateLimitManager{
		limiters: make(map[RateLimitStrategy]RateLimiter),
		cache:    cache,
	}
}

// RegisterLimiter レート制限を登録
func (m *RateLimitManager) RegisterLimiter(strategy RateLimitStrategy, config *RateLimitConfig) {
	switch strategy {
	case StrategyFixedWindow:
		m.limiters[strategy] = NewFixedWindowRateLimiter(config, m.cache)
	case StrategySlidingWindow:
		m.limiters[strategy] = NewSlidingWindowRateLimiter(config, m.cache)
	case StrategyTokenBucket:
		m.limiters[strategy] = NewTokenBucketRateLimiter(config, m.cache)
	default:
		logger.Warn("未対応のレート制限戦略", map[string]interface{}{
			"strategy": strategy,
		})
	}
}

// Allow リクエストを許可するかチェック
func (m *RateLimitManager) Allow(ctx context.Context, strategy RateLimitStrategy, key string) (*RateLimitResult, error) {
	limiter, exists := m.limiters[strategy]
	if !exists {
		return nil, fmt.Errorf("レート制限戦略 '%s' が登録されていません", strategy)
	}

	return limiter.Allow(ctx, key)
}

// Reset レート制限をリセット
func (m *RateLimitManager) Reset(ctx context.Context, strategy RateLimitStrategy, key string) error {
	limiter, exists := m.limiters[strategy]
	if !exists {
		return fmt.Errorf("レート制限戦略 '%s' が登録されていません", strategy)
	}

	return limiter.Reset(ctx, key)
}

// GetLimit 現在の制限情報を取得
func (m *RateLimitManager) GetLimit(ctx context.Context, strategy RateLimitStrategy, key string) (*RateLimitResult, error) {
	limiter, exists := m.limiters[strategy]
	if !exists {
		return nil, fmt.Errorf("レート制限戦略 '%s' が登録されていません", strategy)
	}

	return limiter.GetLimit(ctx, key)
}

// グローバルレート制限管理
var globalRateLimitManager *RateLimitManager

// InitGlobalRateLimitManager グローバルレート制限管理を初期化
func InitGlobalRateLimitManager(cache *cache.CacheManager) {
	globalRateLimitManager = NewRateLimitManager(cache)
}

// GetGlobalRateLimitManager グローバルレート制限管理を取得
func GetGlobalRateLimitManager() *RateLimitManager {
	if globalRateLimitManager == nil {
		globalRateLimitManager = NewRateLimitManager(cache.GetGlobalCacheManager())
	}
	return globalRateLimitManager
}

// AllowGlobal グローバルレート制限をチェック
func AllowGlobal(ctx context.Context, strategy RateLimitStrategy, key string) (*RateLimitResult, error) {
	return GetGlobalRateLimitManager().Allow(ctx, strategy, key)
}

// ResetGlobal グローバルレート制限をリセット
func ResetGlobal(ctx context.Context, strategy RateLimitStrategy, key string) error {
	return GetGlobalRateLimitManager().Reset(ctx, strategy, key)
}

// GetLimitGlobal グローバルレート制限情報を取得
func GetLimitGlobal(ctx context.Context, strategy RateLimitStrategy, key string) (*RateLimitResult, error) {
	return GetGlobalRateLimitManager().GetLimit(ctx, strategy, key)
}
