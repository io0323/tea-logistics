package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// DynamicConfigManager 動的設定管理
type DynamicConfigManager struct {
	config     *Config
	mu         sync.RWMutex
	watchers   []ConfigWatcher
	updateChan chan *Config
	ctx        context.Context
	cancel     context.CancelFunc
}

// ConfigWatcher 設定変更監視者
type ConfigWatcher interface {
	OnConfigChange(config *Config) error
}

// ConfigChangeCallback 設定変更コールバック
type ConfigChangeCallback func(config *Config) error

// NewDynamicConfigManager 新しい動的設定管理を作成
func NewDynamicConfigManager(initialConfig *Config) *DynamicConfigManager {
	ctx, cancel := context.WithCancel(context.Background())

	manager := &DynamicConfigManager{
		config:     initialConfig,
		watchers:   make([]ConfigWatcher, 0),
		updateChan: make(chan *Config, 10),
		ctx:        ctx,
		cancel:     cancel,
	}

	// 設定更新処理を開始
	go manager.processUpdates()

	return manager
}

// GetConfig 現在の設定を取得
func (d *DynamicConfigManager) GetConfig() *Config {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// 設定のコピーを返す
	configCopy := *d.config
	return &configCopy
}

// UpdateConfig 設定を更新
func (d *DynamicConfigManager) UpdateConfig(newConfig *Config) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 設定の妥当性を検証
	if _, err := newConfig.GetOutput(); err != nil {
		return fmt.Errorf("設定の妥当性検証に失敗: %v", err)
	}

	// 設定を更新
	d.config = newConfig

	// 非同期で更新を通知
	select {
	case d.updateChan <- newConfig:
	default:
		// チャンネルが満杯の場合は警告を出力
		Error("設定更新チャンネルが満杯です", map[string]interface{}{
			"config": newConfig,
		})
	}

	return nil
}

// UpdateConfigFields 設定の一部フィールドを更新
func (d *DynamicConfigManager) UpdateConfigFields(updates map[string]interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 現在の設定のコピーを作成
	newConfig := *d.config

	// フィールドを更新
	for key, value := range updates {
		switch key {
		case "level":
			if levelStr, ok := value.(string); ok {
				newConfig.Level = ParseLogLevel(levelStr)
			}
		case "output":
			if output, ok := value.(string); ok {
				newConfig.Output = output
			}
		case "caller":
			if caller, ok := value.(bool); ok {
				newConfig.Caller = caller
			}
		case "pretty":
			if pretty, ok := value.(bool); ok {
				newConfig.Pretty = pretty
			}
		case "time_format":
			if timeFormat, ok := value.(string); ok {
				newConfig.TimeFormat = timeFormat
			}
		default:
			return fmt.Errorf("不明な設定キー: %s", key)
		}
	}

	// 設定の妥当性を検証
	if _, err := newConfig.GetOutput(); err != nil {
		return fmt.Errorf("設定の妥当性検証に失敗: %v", err)
	}

	// 設定を更新
	d.config = &newConfig

	// 非同期で更新を通知
	select {
	case d.updateChan <- &newConfig:
	default:
		Error("設定更新チャンネルが満杯です", map[string]interface{}{
			"updates": updates,
		})
	}

	return nil
}

// AddWatcher 設定変更監視者を追加
func (d *DynamicConfigManager) AddWatcher(watcher ConfigWatcher) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.watchers = append(d.watchers, watcher)
}

// RemoveWatcher 設定変更監視者を削除
func (d *DynamicConfigManager) RemoveWatcher(watcher ConfigWatcher) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, w := range d.watchers {
		if w == watcher {
			d.watchers = append(d.watchers[:i], d.watchers[i+1:]...)
			break
		}
	}
}

// AddCallback 設定変更コールバックを追加
func (d *DynamicConfigManager) AddCallback(callback ConfigChangeCallback) {
	watcher := &callbackWatcher{callback: callback}
	d.AddWatcher(watcher)
}

// processUpdates 設定更新の処理
func (d *DynamicConfigManager) processUpdates() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case config := <-d.updateChan:
			d.notifyWatchers(config)
		}
	}
}

// notifyWatchers 監視者に設定変更を通知
func (d *DynamicConfigManager) notifyWatchers(config *Config) {
	d.mu.RLock()
	watchers := make([]ConfigWatcher, len(d.watchers))
	copy(watchers, d.watchers)
	d.mu.RUnlock()

	for _, watcher := range watchers {
		if err := watcher.OnConfigChange(config); err != nil {
			Error("設定変更通知エラー", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}

// Close 動的設定管理を閉じる
func (d *DynamicConfigManager) Close() {
	d.cancel()
	close(d.updateChan)
}

// callbackWatcher コールバック監視者
type callbackWatcher struct {
	callback ConfigChangeCallback
}

// OnConfigChange 設定変更時のコールバック実行
func (c *callbackWatcher) OnConfigChange(config *Config) error {
	return c.callback(config)
}

// ConfigReloader 設定再読み込み機能
type ConfigReloader struct {
	configPath string
	lastMod    time.Time
	manager    *DynamicConfigManager
}

// NewConfigReloader 新しい設定再読み込み機能を作成
func NewConfigReloader(configPath string, manager *DynamicConfigManager) *ConfigReloader {
	return &ConfigReloader{
		configPath: configPath,
		manager:    manager,
	}
}

// StartWatching 設定ファイルの監視を開始
func (r *ConfigReloader) StartWatching(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := r.checkAndReload(); err != nil {
					Error("設定ファイルの再読み込みエラー", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
	}()
}

// checkAndReload 設定ファイルの変更をチェックして再読み込み
func (r *ConfigReloader) checkAndReload() error {
	if r.configPath == "" {
		return nil
	}

	// ファイルの更新時刻をチェック
	info, err := os.Stat(r.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // ファイルが存在しない場合は何もしない
		}
		return err
	}

	modTime := info.ModTime()
	if modTime.After(r.lastMod) {
		r.lastMod = modTime

		// 設定ファイルを読み込み
		config, err := LoadConfigFromFile(r.configPath)
		if err != nil {
			return err
		}

		// 動的設定管理に適用
		return r.manager.UpdateConfig(config)
	}

	return nil
}

// LoadConfigFromFile ファイルから設定を読み込み
func LoadConfigFromFile(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// グローバル動的設定管理
var globalDynamicManager *DynamicConfigManager

// InitGlobalDynamicConfig グローバル動的設定管理を初期化
func InitGlobalDynamicConfig(initialConfig *Config) {
	globalDynamicManager = NewDynamicConfigManager(initialConfig)

	// グローバルロガーを監視者として追加
	globalDynamicManager.AddCallback(func(config *Config) error {
		return config.Apply(GetGlobalLogger())
	})
}

// GetGlobalDynamicConfig グローバル動的設定管理を取得
func GetGlobalDynamicConfig() *DynamicConfigManager {
	if globalDynamicManager == nil {
		globalDynamicManager = NewDynamicConfigManager(DefaultConfig())
	}
	return globalDynamicManager
}

// UpdateGlobalConfigDynamic グローバル設定を動的に更新
func UpdateGlobalConfigDynamic(updates map[string]interface{}) error {
	return GetGlobalDynamicConfig().UpdateConfigFields(updates)
}

// GetGlobalConfigDynamic グローバル設定を動的に取得
func GetGlobalConfigDynamic() *Config {
	return GetGlobalDynamicConfig().GetConfig()
}
