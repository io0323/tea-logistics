package logger

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogConfigManager(t *testing.T) {
	// 一時ディレクトリを作成
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "log_config.json")

	manager := NewLogConfigManager(configPath)

	t.Run("初期設定", func(t *testing.T) {
		config := manager.GetConfig()
		assert.Equal(t, INFO, config.Level)
		assert.Equal(t, "stdout", config.Output)
		assert.True(t, config.Caller)
	})

	t.Run("設定の保存と読み込み", func(t *testing.T) {
		// 設定を変更
		config := &Config{
			Level:      DEBUG,
			Output:     "stderr",
			Caller:     false,
			Pretty:     true,
			TimeFormat: "2006-01-02 15:04:05",
		}
		manager.SetConfig(config)

		// 設定を保存
		err := manager.SaveConfig()
		require.NoError(t, err)

		// 新しいマネージャーで設定を読み込み
		newManager := NewLogConfigManager(configPath)
		err = newManager.LoadConfig()
		require.NoError(t, err)

		loadedConfig := newManager.GetConfig()
		assert.Equal(t, DEBUG, loadedConfig.Level)
		assert.Equal(t, "stderr", loadedConfig.Output)
		assert.False(t, loadedConfig.Caller)
		assert.True(t, loadedConfig.Pretty)
		assert.Equal(t, "2006-01-02 15:04:05", loadedConfig.TimeFormat)
	})

	t.Run("設定の部分更新", func(t *testing.T) {
		updates := map[string]interface{}{
			"level": "ERROR",
			"caller": false,
		}

		err := manager.UpdateConfig(updates)
		require.NoError(t, err)

		config := manager.GetConfig()
		assert.Equal(t, ERROR, config.Level)
		assert.False(t, config.Caller)
	})

	t.Run("設定の妥当性検証", func(t *testing.T) {
		// 有効な設定
		err := manager.ValidateConfig()
		assert.NoError(t, err)

		// 無効な出力先
		invalidConfig := &Config{
			Level:  INFO,
			Output: "/invalid/path/that/does/not/exist/test.log",
		}
		manager.SetConfig(invalidConfig)

		err = manager.ValidateConfig()
		assert.Error(t, err)
	})

	t.Run("設定のエクスポートとインポート", func(t *testing.T) {
		config := &Config{
			Level:      WARN,
			Output:     "stdout",
			Caller:     true,
			Pretty:     false,
			TimeFormat: "2006-01-02T15:04:05Z07:00",
		}
		manager.SetConfig(config)

		// エクスポート
		data, err := manager.ExportConfig()
		require.NoError(t, err)
		assert.NotEmpty(t, data)

		// インポート
		newManager := NewLogConfigManager("")
		err = newManager.ImportConfig(data)
		require.NoError(t, err)

		importedConfig := newManager.GetConfig()
		assert.Equal(t, WARN, importedConfig.Level)
		assert.Equal(t, "stdout", importedConfig.Output)
		assert.True(t, importedConfig.Caller)
		assert.False(t, importedConfig.Pretty)
	})

	t.Run("設定の概要取得", func(t *testing.T) {
		summary := manager.GetConfigSummary()
		assert.Contains(t, summary, "level")
		assert.Contains(t, summary, "output")
		assert.Contains(t, summary, "caller")
		assert.Contains(t, summary, "pretty")
		assert.Contains(t, summary, "time_format")
	})
}

func TestEnvVarManager(t *testing.T) {
	manager := NewEnvVarManager("TEST_LOG")

	t.Run("文字列の取得と設定", func(t *testing.T) {
		// デフォルト値のテスト
		value := manager.GetString("OUTPUT", "default")
		assert.Equal(t, "default", value)

		// 環境変数の設定
		err := manager.SetString("OUTPUT", "test_value")
		require.NoError(t, err)

		// 環境変数の取得
		value = manager.GetString("OUTPUT", "default")
		assert.Equal(t, "test_value", value)

		// クリーンアップ
		manager.Unset("OUTPUT")
	})

	t.Run("整数の取得と設定", func(t *testing.T) {
		// デフォルト値のテスト
		value := manager.GetInt("MAX_SIZE", 1024)
		assert.Equal(t, 1024, value)

		// 環境変数の設定
		err := manager.SetInt("MAX_SIZE", 2048)
		require.NoError(t, err)

		// 環境変数の取得
		value = manager.GetInt("MAX_SIZE", 1024)
		assert.Equal(t, 2048, value)

		// クリーンアップ
		manager.Unset("MAX_SIZE")
	})

	t.Run("ブール値の取得と設定", func(t *testing.T) {
		// デフォルト値のテスト
		value := manager.GetBool("ENABLED", true)
		assert.True(t, value)

		// 環境変数の設定
		err := manager.SetBool("ENABLED", false)
		require.NoError(t, err)

		// 環境変数の取得
		value = manager.GetBool("ENABLED", true)
		assert.False(t, value)

		// クリーンアップ
		manager.Unset("ENABLED")
	})

	t.Run("ログレベルの取得と設定", func(t *testing.T) {
		// デフォルト値のテスト
		value := manager.GetLogLevel("LEVEL", INFO)
		assert.Equal(t, INFO, value)

		// 環境変数の設定
		err := manager.SetLogLevel("LEVEL", DEBUG)
		require.NoError(t, err)

		// 環境変数の取得
		value = manager.GetLogLevel("LEVEL", INFO)
		assert.Equal(t, DEBUG, value)

		// クリーンアップ
		manager.Unset("LEVEL")
	})

	t.Run("環境変数の存在チェック", func(t *testing.T) {
		// 存在しない環境変数
		exists := manager.Exists("NONEXISTENT")
		assert.False(t, exists)

		// 環境変数を設定
		manager.SetString("EXISTENT", "value")
		exists = manager.Exists("EXISTENT")
		assert.True(t, exists)

		// クリーンアップ
		manager.Unset("EXISTENT")
	})

	t.Run("全ての環境変数の取得", func(t *testing.T) {
		// テスト用の環境変数を設定
		manager.SetString("VAR1", "value1")
		manager.SetString("VAR2", "value2")
		manager.SetInt("VAR3", 123)

		all := manager.GetAll()
		assert.Contains(t, all, "VAR1")
		assert.Contains(t, all, "VAR2")
		assert.Contains(t, all, "VAR3")
		assert.Equal(t, "value1", all["VAR1"])
		assert.Equal(t, "value2", all["VAR2"])
		assert.Equal(t, "123", all["VAR3"])

		// クリーンアップ
		manager.Unset("VAR1")
		manager.Unset("VAR2")
		manager.Unset("VAR3")
	})

	t.Run("環境変数から設定を構築", func(t *testing.T) {
		// 環境変数を設定
		manager.SetLogLevel("LEVEL", WARN)
		manager.SetString("OUTPUT", "stderr")
		manager.SetBool("CALLER", false)
		manager.SetBool("PRETTY", true)
		manager.SetString("TIME_FORMAT", "2006-01-02 15:04:05")

		config := manager.GetConfigFromEnv()
		assert.Equal(t, WARN, config.Level)
		assert.Equal(t, "stderr", config.Output)
		assert.False(t, config.Caller)
		assert.True(t, config.Pretty)
		assert.Equal(t, "2006-01-02 15:04:05", config.TimeFormat)

		// クリーンアップ
		manager.Unset("LEVEL")
		manager.Unset("OUTPUT")
		manager.Unset("CALLER")
		manager.Unset("PRETTY")
		manager.Unset("TIME_FORMAT")
	})

	t.Run("設定を環境変数に適用", func(t *testing.T) {
		config := &Config{
			Level:      ERROR,
			Output:     "stdout",
			Caller:     true,
			Pretty:     false,
			TimeFormat: "2006-01-02T15:04:05Z07:00",
		}

		err := manager.ApplyConfigToEnv(config)
		require.NoError(t, err)

		// 環境変数を確認
		assert.Equal(t, "ERROR", manager.GetString("LEVEL", ""))
		assert.Equal(t, "stdout", manager.GetString("OUTPUT", ""))
		assert.Equal(t, "true", manager.GetString("CALLER", ""))
		assert.Equal(t, "false", manager.GetString("PRETTY", ""))
		assert.Equal(t, "2006-01-02T15:04:05Z07:00", manager.GetString("TIME_FORMAT", ""))

		// クリーンアップ
		manager.Unset("LEVEL")
		manager.Unset("OUTPUT")
		manager.Unset("CALLER")
		manager.Unset("PRETTY")
		manager.Unset("TIME_FORMAT")
	})

	t.Run("環境変数の妥当性検証", func(t *testing.T) {
		// 有効な環境変数
		manager.SetString("LEVEL", "INFO")
		manager.SetString("OUTPUT", "stdout")
		manager.SetString("TIME_FORMAT", "2006-01-02T15:04:05Z07:00")

		errors := manager.ValidateEnv()
		assert.Empty(t, errors)

		// 無効なログレベル
		manager.SetString("LEVEL", "INVALID")
		errors = manager.ValidateEnv()
		assert.NotEmpty(t, errors)

		// クリーンアップ
		manager.Unset("LEVEL")
		manager.Unset("OUTPUT")
		manager.Unset("TIME_FORMAT")
	})
}

func TestDynamicConfigManager(t *testing.T) {
	initialConfig := &Config{
		Level:  INFO,
		Output: "stdout",
		Caller: true,
	}

	manager := NewDynamicConfigManager(initialConfig)
	defer manager.Close()

	t.Run("初期設定の取得", func(t *testing.T) {
		config := manager.GetConfig()
		assert.Equal(t, INFO, config.Level)
		assert.Equal(t, "stdout", config.Output)
		assert.True(t, config.Caller)
	})

	t.Run("設定の更新", func(t *testing.T) {
		newConfig := &Config{
			Level:  DEBUG,
			Output: "stderr",
			Caller: false,
		}

		err := manager.UpdateConfig(newConfig)
		require.NoError(t, err)

		config := manager.GetConfig()
		assert.Equal(t, DEBUG, config.Level)
		assert.Equal(t, "stderr", config.Output)
		assert.False(t, config.Caller)
	})

	t.Run("設定の部分更新", func(t *testing.T) {
		updates := map[string]interface{}{
			"level": "WARN",
			"caller": true,
		}

		err := manager.UpdateConfigFields(updates)
		require.NoError(t, err)

		config := manager.GetConfig()
		assert.Equal(t, WARN, config.Level)
		assert.True(t, config.Caller)
	})

	t.Run("設定変更監視者", func(t *testing.T) {
		var receivedConfig *Config
		var callbackCalled bool

		callback := func(config *Config) error {
			receivedConfig = config
			callbackCalled = true
			return nil
		}

		manager.AddCallback(callback)

		// 設定を更新
		newConfig := &Config{
			Level:  ERROR,
			Output: "stdout",
			Caller: false,
		}

		err := manager.UpdateConfig(newConfig)
		require.NoError(t, err)

		// コールバックが呼ばれるまで少し待つ
		time.Sleep(100 * time.Millisecond)

		assert.True(t, callbackCalled)
		assert.Equal(t, ERROR, receivedConfig.Level)
	})

	t.Run("無効な設定の更新", func(t *testing.T) {
		invalidConfig := &Config{
			Level:  INFO,
			Output: "/invalid/path/that/does/not/exist/test.log",
		}

		err := manager.UpdateConfig(invalidConfig)
		assert.Error(t, err)
	})

	t.Run("不明な設定キー", func(t *testing.T) {
		updates := map[string]interface{}{
			"unknown_key": "value",
		}

		err := manager.UpdateConfigFields(updates)
		assert.Error(t, err)
	})
}

func TestGlobalConfigManagement(t *testing.T) {
	t.Run("グローバル設定管理の初期化", func(t *testing.T) {
		err := InitGlobalConfigManager("")
		require.NoError(t, err)

		manager := GetGlobalConfigManager()
		assert.NotNil(t, manager)
	})

	t.Run("グローバル環境変数管理", func(t *testing.T) {
		InitGlobalEnvManager("TEST")

		manager := GetGlobalEnvManager()
		assert.NotNil(t, manager)

		// テスト用の環境変数を設定
		manager.SetString("LEVEL", "WARN")
		manager.SetString("OUTPUT", "stderr")

		config := LoadConfigFromGlobalEnv()
		assert.Equal(t, WARN, config.Level)
		assert.Equal(t, "stderr", config.Output)

		// クリーンアップ
		manager.Unset("LEVEL")
		manager.Unset("OUTPUT")
	})

	t.Run("グローバル動的設定管理", func(t *testing.T) {
		config := &Config{
			Level:  INFO,
			Output: "stdout",
			Caller: true,
		}

		InitGlobalDynamicConfig(config)

		manager := GetGlobalDynamicConfig()
		assert.NotNil(t, manager)

		// 設定を動的に更新
		updates := map[string]interface{}{
			"level": "ERROR",
		}

		err := UpdateGlobalConfigDynamic(updates)
		require.NoError(t, err)

		updatedConfig := GetGlobalConfigDynamic()
		assert.Equal(t, ERROR, updatedConfig.Level)
	})
}
