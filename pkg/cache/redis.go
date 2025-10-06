package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tea-logistics/pkg/logger"

	"github.com/go-redis/redis/v8"
)

// CacheConfig キャッシュ設定
type CacheConfig struct {
	Host     string        `json:"host"`
	Port     int           `json:"port"`
	Password string        `json:"password"`
	DB       int           `json:"db"`
	PoolSize int           `json:"pool_size"`
	Timeout  time.Duration `json:"timeout"`
}

// DefaultCacheConfig デフォルトキャッシュ設定
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
		PoolSize: 10,
		Timeout:  5 * time.Second,
	}
}

// CacheManager キャッシュ管理
type CacheManager struct {
	client *redis.Client
	config *CacheConfig
}

// NewCacheManager 新しいキャッシュ管理を作成
func NewCacheManager(config *CacheConfig) *CacheManager {
	if config == nil {
		config = DefaultCacheConfig()
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	return &CacheManager{
		client: client,
		config: config,
	}
}

// Connect キャッシュに接続
func (c *CacheManager) Connect(ctx context.Context) error {
	_, err := c.client.Ping(ctx).Result()
	if err != nil {
		logger.Error("Redis接続エラー", map[string]interface{}{
			"host": c.config.Host,
			"port": c.config.Port,
			"error": err.Error(),
		})
		return fmt.Errorf("Redis接続に失敗しました: %v", err)
	}

	logger.Info("Redis接続が確立されました", map[string]interface{}{
		"host": c.config.Host,
		"port": c.config.Port,
		"db":   c.config.DB,
	})

	return nil
}

// Close キャッシュ接続を閉じる
func (c *CacheManager) Close() error {
	err := c.client.Close()
	if err != nil {
		logger.Error("Redis接続のクローズエラー", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	logger.Info("Redis接続をクローズしました")
	return nil
}

// Set 値を設定
func (c *CacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var data []byte
	var err error

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			logger.Error("JSONマーシャルエラー", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
			return fmt.Errorf("値のシリアライズに失敗しました: %v", err)
		}
	}

	err = c.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		logger.Error("Redis Setエラー", map[string]interface{}{
			"key":        key,
			"expiration": expiration.String(),
			"error":      err.Error(),
		})
		return fmt.Errorf("キャッシュの設定に失敗しました: %v", err)
	}

	logger.Debug("キャッシュに値を設定", map[string]interface{}{
		"key":        key,
		"expiration": expiration.String(),
	})

	return nil
}

// Get 値を取得
func (c *CacheManager) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			logger.Debug("キャッシュキーが見つかりません", map[string]interface{}{
				"key": key,
			})
			return "", fmt.Errorf("キー '%s' が見つかりません", key)
		}

		logger.Error("Redis Getエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return "", fmt.Errorf("キャッシュの取得に失敗しました: %v", err)
	}

	logger.Debug("キャッシュから値を取得", map[string]interface{}{
		"key": key,
	})

	return val, nil
}

// GetObject オブジェクトを取得
func (c *CacheManager) GetObject(ctx context.Context, key string, dest interface{}) error {
	val, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		logger.Error("JSONアンマーシャルエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return fmt.Errorf("値のデシリアライズに失敗しました: %v", err)
	}

	return nil
}

// Delete 値を削除
func (c *CacheManager) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		logger.Error("Redis Deleteエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return fmt.Errorf("キャッシュの削除に失敗しました: %v", err)
	}

	logger.Debug("キャッシュから値を削除", map[string]interface{}{
		"key": key,
	})

	return nil
}

// Exists キーの存在確認
func (c *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		logger.Error("Redis Existsエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return false, fmt.Errorf("キーの存在確認に失敗しました: %v", err)
	}

	return count > 0, nil
}

// Expire キーの有効期限を設定
func (c *CacheManager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := c.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		logger.Error("Redis Expireエラー", map[string]interface{}{
			"key":        key,
			"expiration": expiration.String(),
			"error":      err.Error(),
		})
		return fmt.Errorf("有効期限の設定に失敗しました: %v", err)
	}

	logger.Debug("キャッシュの有効期限を設定", map[string]interface{}{
		"key":        key,
		"expiration": expiration.String(),
	})

	return nil
}

// TTL キーの残り有効期限を取得
func (c *CacheManager) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		logger.Error("Redis TTLエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return 0, fmt.Errorf("有効期限の取得に失敗しました: %v", err)
	}

	return ttl, nil
}

// Keys パターンに一致するキーを取得
func (c *CacheManager) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		logger.Error("Redis Keysエラー", map[string]interface{}{
			"pattern": pattern,
			"error":   err.Error(),
		})
		return nil, fmt.Errorf("キーの検索に失敗しました: %v", err)
	}

	logger.Debug("キャッシュキーを検索", map[string]interface{}{
		"pattern": pattern,
		"count":   len(keys),
	})

	return keys, nil
}

// FlushDB データベースをクリア
func (c *CacheManager) FlushDB(ctx context.Context) error {
	err := c.client.FlushDB(ctx).Err()
	if err != nil {
		logger.Error("Redis FlushDBエラー", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("データベースのクリアに失敗しました: %v", err)
	}

	logger.Info("Redisデータベースをクリアしました")
	return nil
}

// GetStats キャッシュ統計を取得
func (c *CacheManager) GetStats(ctx context.Context) (map[string]interface{}, error) {
	info, err := c.client.Info(ctx).Result()
	if err != nil {
		logger.Error("Redis Infoエラー", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("統計情報の取得に失敗しました: %v", err)
	}

	stats := map[string]interface{}{
		"info": info,
	}

	logger.Debug("Redis統計情報を取得")
	return stats, nil
}

// SetNX キーが存在しない場合のみ設定
func (c *CacheManager) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	var data []byte
	var err error

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			logger.Error("JSONマーシャルエラー", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
			return false, fmt.Errorf("値のシリアライズに失敗しました: %v", err)
		}
	}

	result, err := c.client.SetNX(ctx, key, data, expiration).Result()
	if err != nil {
		logger.Error("Redis SetNXエラー", map[string]interface{}{
			"key":        key,
			"expiration": expiration.String(),
			"error":      err.Error(),
		})
		return false, fmt.Errorf("キャッシュの設定に失敗しました: %v", err)
	}

	logger.Debug("キャッシュに値を設定（NX）", map[string]interface{}{
		"key":        key,
		"expiration": expiration.String(),
		"success":    result,
	})

	return result, nil
}

// Increment 数値を増加
func (c *CacheManager) Increment(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		logger.Error("Redis Incrエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return 0, fmt.Errorf("数値の増加に失敗しました: %v", err)
	}

	logger.Debug("キャッシュの数値を増加", map[string]interface{}{
		"key":   key,
		"value": val,
	})

	return val, nil
}

// IncrementBy 指定した値だけ数値を増加
func (c *CacheManager) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := c.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		logger.Error("Redis IncrByエラー", map[string]interface{}{
			"key":   key,
			"value": value,
			"error": err.Error(),
		})
		return 0, fmt.Errorf("数値の増加に失敗しました: %v", err)
	}

	logger.Debug("キャッシュの数値を増加", map[string]interface{}{
		"key":   key,
		"value": val,
	})

	return val, nil
}

// Decrement 数値を減少
func (c *CacheManager) Decrement(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Decr(ctx, key).Result()
	if err != nil {
		logger.Error("Redis Decrエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return 0, fmt.Errorf("数値の減少に失敗しました: %v", err)
	}

	logger.Debug("キャッシュの数値を減少", map[string]interface{}{
		"key":   key,
		"value": val,
	})

	return val, nil
}

// DecrementBy 指定した値だけ数値を減少
func (c *CacheManager) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := c.client.DecrBy(ctx, key, value).Result()
	if err != nil {
		logger.Error("Redis DecrByエラー", map[string]interface{}{
			"key":   key,
			"value": value,
			"error": err.Error(),
		})
		return 0, fmt.Errorf("数値の減少に失敗しました: %v", err)
	}

	logger.Debug("キャッシュの数値を減少", map[string]interface{}{
		"key":   key,
		"value": val,
	})

	return val, nil
}

// ListPush リストに値を追加
func (c *CacheManager) ListPush(ctx context.Context, key string, values ...interface{}) error {
	var data []interface{}
	for _, v := range values {
		switch val := v.(type) {
		case string:
			data = append(data, val)
		case []byte:
			data = append(data, string(val))
		default:
			jsonData, err := json.Marshal(val)
			if err != nil {
				logger.Error("JSONマーシャルエラー", map[string]interface{}{
					"key":   key,
					"error": err.Error(),
				})
				return fmt.Errorf("値のシリアライズに失敗しました: %v", err)
			}
			data = append(data, string(jsonData))
		}
	}

	err := c.client.LPush(ctx, key, data...).Err()
	if err != nil {
		logger.Error("Redis LPushエラー", map[string]interface{}{
			"key":    key,
			"values": len(values),
			"error":  err.Error(),
		})
		return fmt.Errorf("リストへの追加に失敗しました: %v", err)
	}

	logger.Debug("リストに値を追加", map[string]interface{}{
		"key":    key,
		"values": len(values),
	})

	return nil
}

// ListPop リストから値を取得
func (c *CacheManager) ListPop(ctx context.Context, key string) (string, error) {
	val, err := c.client.RPop(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			logger.Debug("リストが空です", map[string]interface{}{
				"key": key,
			})
			return "", fmt.Errorf("リスト '%s' が空です", key)
		}

		logger.Error("Redis RPopエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return "", fmt.Errorf("リストからの取得に失敗しました: %v", err)
	}

	logger.Debug("リストから値を取得", map[string]interface{}{
		"key": key,
	})

	return val, nil
}

// ListLength リストの長さを取得
func (c *CacheManager) ListLength(ctx context.Context, key string) (int64, error) {
	length, err := c.client.LLen(ctx, key).Result()
	if err != nil {
		logger.Error("Redis LLenエラー", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return 0, fmt.Errorf("リストの長さ取得に失敗しました: %v", err)
	}

	return length, nil
}

// グローバルキャッシュ管理
var globalCacheManager *CacheManager

// InitGlobalCacheManager グローバルキャッシュ管理を初期化
func InitGlobalCacheManager(config *CacheConfig) error {
	globalCacheManager = NewCacheManager(config)
	ctx := context.Background()
	return globalCacheManager.Connect(ctx)
}

// GetGlobalCacheManager グローバルキャッシュ管理を取得
func GetGlobalCacheManager() *CacheManager {
	if globalCacheManager == nil {
		globalCacheManager = NewCacheManager(DefaultCacheConfig())
	}
	return globalCacheManager
}

// CloseGlobalCacheManager グローバルキャッシュ管理を閉じる
func CloseGlobalCacheManager() error {
	if globalCacheManager != nil {
		return globalCacheManager.Close()
	}
	return nil
}
