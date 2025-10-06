package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LogConfigManager ログ設定管理
type LogConfigManager struct {
	configPath string
	config     *Config
}

// NewLogConfigManager 新しいログ設定管理を作成
func NewLogConfigManager(configPath string) *LogConfigManager {
	return &LogConfigManager{
		configPath: configPath,
		config:     DefaultConfig(),
	}
}

// LoadConfig 設定ファイルから設定を読み込み
func (m *LogConfigManager) LoadConfig() error {
	if m.configPath == "" {
		return nil // 設定ファイルが指定されていない場合はデフォルト設定を使用
	}

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 設定ファイルが存在しない場合はデフォルト設定を作成
			return m.SaveConfig()
		}
		return fmt.Errorf("設定ファイルの読み込みに失敗しました: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("設定ファイルの解析に失敗しました: %v", err)
	}

	m.config = &config
	return nil
}

// SaveConfig 設定をファイルに保存
func (m *LogConfigManager) SaveConfig() error {
	if m.configPath == "" {
		return nil // 設定ファイルが指定されていない場合は何もしない
	}

	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗しました: %v", err)
	}

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("設定のシリアライズに失敗しました: %v", err)
	}

	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("設定ファイルの保存に失敗しました: %v", err)
	}

	return nil
}

// GetConfig 現在の設定を取得
func (m *LogConfigManager) GetConfig() *Config {
	return m.config
}

// SetConfig 設定を更新
func (m *LogConfigManager) SetConfig(config *Config) {
	m.config = config
}

// UpdateConfig 設定の一部を更新
func (m *LogConfigManager) UpdateConfig(updates map[string]interface{}) error {
	if m.config == nil {
		m.config = DefaultConfig()
	}

	for key, value := range updates {
		switch key {
		case "level":
			if levelStr, ok := value.(string); ok {
				m.config.Level = ParseLogLevel(levelStr)
			}
		case "output":
			if output, ok := value.(string); ok {
				m.config.Output = output
			}
		case "caller":
			if caller, ok := value.(bool); ok {
				m.config.Caller = caller
			}
		case "pretty":
			if pretty, ok := value.(bool); ok {
				m.config.Pretty = pretty
			}
		case "time_format":
			if timeFormat, ok := value.(string); ok {
				m.config.TimeFormat = timeFormat
			}
		default:
			return fmt.Errorf("不明な設定キー: %s", key)
		}
	}

	return nil
}

// ValidateConfig 設定の妥当性を検証
func (m *LogConfigManager) ValidateConfig() error {
	if m.config == nil {
		return fmt.Errorf("設定が初期化されていません")
	}

	// 出力先の検証
	if _, err := m.config.GetOutput(); err != nil {
		return fmt.Errorf("出力先の設定が無効です: %v", err)
	}

	// 時間フォーマットの検証（簡易版）
	if m.config.TimeFormat != "" {
		// 基本的な文字列長チェック
		if len(m.config.TimeFormat) < 5 {
			return fmt.Errorf("時間フォーマットが短すぎます")
		}
	}

	return nil
}

// ApplyToLogger 設定をロガーに適用
func (m *LogConfigManager) ApplyToLogger(logger *Logger) error {
	if err := m.ValidateConfig(); err != nil {
		return err
	}

	return m.config.Apply(logger)
}

// ResetToDefault 設定をデフォルトにリセット
func (m *LogConfigManager) ResetToDefault() {
	m.config = DefaultConfig()
}

// ExportConfig 設定をJSON形式でエクスポート
func (m *LogConfigManager) ExportConfig() ([]byte, error) {
	if m.config == nil {
		return nil, fmt.Errorf("設定が初期化されていません")
	}

	return json.MarshalIndent(m.config, "", "  ")
}

// ImportConfig JSON形式の設定をインポート
func (m *LogConfigManager) ImportConfig(data []byte) error {
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("設定の解析に失敗しました: %v", err)
	}

	m.config = &config
	return nil
}

// GetConfigSummary 設定の概要を取得
func (m *LogConfigManager) GetConfigSummary() map[string]interface{} {
	if m.config == nil {
		return map[string]interface{}{
			"status": "not_initialized",
		}
	}

	return map[string]interface{}{
		"level":       m.config.Level.String(),
		"output":      m.config.Output,
		"caller":      m.config.Caller,
		"pretty":      m.config.Pretty,
		"time_format": m.config.TimeFormat,
		"config_path": m.configPath,
	}
}

// グローバル設定管理
var globalConfigManager *LogConfigManager

// InitGlobalConfigManager グローバル設定管理を初期化
func InitGlobalConfigManager(configPath string) error {
	globalConfigManager = NewLogConfigManager(configPath)
	return globalConfigManager.LoadConfig()
}

// GetGlobalConfigManager グローバル設定管理を取得
func GetGlobalConfigManager() *LogConfigManager {
	if globalConfigManager == nil {
		globalConfigManager = NewLogConfigManager("")
	}
	return globalConfigManager
}

// UpdateGlobalConfig グローバル設定を更新
func UpdateGlobalConfig(updates map[string]interface{}) error {
	manager := GetGlobalConfigManager()
	if err := manager.UpdateConfig(updates); err != nil {
		return err
	}

	// 設定をファイルに保存
	if err := manager.SaveConfig(); err != nil {
		return err
	}

	// グローバルロガーに適用
	return manager.ApplyToLogger(GetGlobalLogger())
}

// GetGlobalConfigSummary グローバル設定の概要を取得
func GetGlobalConfigSummary() map[string]interface{} {
	return GetGlobalConfigManager().GetConfigSummary()
}

// ReloadGlobalConfig グローバル設定を再読み込み
func ReloadGlobalConfig() error {
	manager := GetGlobalConfigManager()
	if err := manager.LoadConfig(); err != nil {
		return err
	}

	return manager.ApplyToLogger(GetGlobalLogger())
}

// SaveGlobalConfig グローバル設定を保存
func SaveGlobalConfig() error {
	return GetGlobalConfigManager().SaveConfig()
}
