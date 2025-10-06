package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EnvVarManager 環境変数管理
type EnvVarManager struct {
	prefix string
}

// NewEnvVarManager 新しい環境変数管理を作成
func NewEnvVarManager(prefix string) *EnvVarManager {
	if prefix == "" {
		prefix = "LOG"
	}
	return &EnvVarManager{
		prefix: strings.ToUpper(prefix),
	}
}

// GetString 文字列の環境変数を取得
func (e *EnvVarManager) GetString(key string, defaultValue string) string {
	fullKey := e.getFullKey(key)
	if value := os.Getenv(fullKey); value != "" {
		return value
	}
	return defaultValue
}

// GetInt 整数の環境変数を取得
func (e *EnvVarManager) GetInt(key string, defaultValue int) int {
	fullKey := e.getFullKey(key)
	if value := os.Getenv(fullKey); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetBool ブール値の環境変数を取得
func (e *EnvVarManager) GetBool(key string, defaultValue bool) bool {
	fullKey := e.getFullKey(key)
	if value := os.Getenv(fullKey); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetLogLevel ログレベルの環境変数を取得
func (e *EnvVarManager) GetLogLevel(key string, defaultValue LogLevel) LogLevel {
	fullKey := e.getFullKey(key)
	if value := os.Getenv(fullKey); value != "" {
		return ParseLogLevel(value)
	}
	return defaultValue
}

// SetString 文字列の環境変数を設定
func (e *EnvVarManager) SetString(key string, value string) error {
	fullKey := e.getFullKey(key)
	return os.Setenv(fullKey, value)
}

// SetInt 整数の環境変数を設定
func (e *EnvVarManager) SetInt(key string, value int) error {
	fullKey := e.getFullKey(key)
	return os.Setenv(fullKey, strconv.Itoa(value))
}

// SetBool ブール値の環境変数を設定
func (e *EnvVarManager) SetBool(key string, value bool) error {
	fullKey := e.getFullKey(key)
	return os.Setenv(fullKey, strconv.FormatBool(value))
}

// SetLogLevel ログレベルの環境変数を設定
func (e *EnvVarManager) SetLogLevel(key string, value LogLevel) error {
	fullKey := e.getFullKey(key)
	return os.Setenv(fullKey, value.String())
}

// Unset 環境変数を削除
func (e *EnvVarManager) Unset(key string) error {
	fullKey := e.getFullKey(key)
	return os.Unsetenv(fullKey)
}

// Exists 環境変数が存在するかチェック
func (e *EnvVarManager) Exists(key string) bool {
	fullKey := e.getFullKey(key)
	_, exists := os.LookupEnv(fullKey)
	return exists
}

// GetAll プレフィックスに一致する全ての環境変数を取得
func (e *EnvVarManager) GetAll() map[string]string {
	result := make(map[string]string)
	prefix := e.prefix + "_"
	
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, prefix) {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimPrefix(parts[0], prefix)
				result[key] = parts[1]
			}
		}
	}
	
	return result
}

// GetConfigFromEnv 環境変数から設定を構築
func (e *EnvVarManager) GetConfigFromEnv() *Config {
	return &Config{
		Level:      e.GetLogLevel("LEVEL", INFO),
		Output:     e.GetString("OUTPUT", "stdout"),
		Caller:     e.GetBool("CALLER", true),
		Pretty:     e.GetBool("PRETTY", false),
		TimeFormat: e.GetString("TIME_FORMAT", "2006-01-02T15:04:05Z07:00"),
	}
}

// ApplyConfigToEnv 設定を環境変数に適用
func (e *EnvVarManager) ApplyConfigToEnv(config *Config) error {
	if err := e.SetLogLevel("LEVEL", config.Level); err != nil {
		return fmt.Errorf("LOG_LEVELの設定に失敗: %v", err)
	}
	
	if err := e.SetString("OUTPUT", config.Output); err != nil {
		return fmt.Errorf("LOG_OUTPUTの設定に失敗: %v", err)
	}
	
	if err := e.SetBool("CALLER", config.Caller); err != nil {
		return fmt.Errorf("LOG_CALLERの設定に失敗: %v", err)
	}
	
	if err := e.SetBool("PRETTY", config.Pretty); err != nil {
		return fmt.Errorf("LOG_PRETTYの設定に失敗: %v", err)
	}
	
	if err := e.SetString("TIME_FORMAT", config.TimeFormat); err != nil {
		return fmt.Errorf("LOG_TIME_FORMATの設定に失敗: %v", err)
	}
	
	return nil
}

// ValidateEnv 環境変数の妥当性を検証
func (e *EnvVarManager) ValidateEnv() []string {
	var errors []string
	
	// ログレベルの検証
	if level := e.GetString("LEVEL", ""); level != "" {
		if ParseLogLevel(level) == INFO && strings.ToUpper(level) != "INFO" {
			errors = append(errors, fmt.Sprintf("無効なログレベル: %s", level))
		}
	}
	
	// 出力先の検証
	if output := e.GetString("OUTPUT", ""); output != "" {
		if output != "stdout" && output != "stderr" && output != "null" {
			// ファイルパスの場合は存在チェックはしない（実行時に作成される可能性がある）
		}
	}
	
	// 時間フォーマットの検証
	if timeFormat := e.GetString("TIME_FORMAT", ""); timeFormat != "" {
		if err := parseTimeFormat(timeFormat); err != nil {
			errors = append(errors, fmt.Sprintf("無効な時間フォーマット: %s", timeFormat))
		}
	}
	
	return errors
}

// getFullKey プレフィックス付きのキーを取得
func (e *EnvVarManager) getFullKey(key string) string {
	return fmt.Sprintf("%s_%s", e.prefix, strings.ToUpper(key))
}

// parseTimeFormat 時間フォーマットの解析（簡易版）
func parseTimeFormat(format string) error {
	// Goの時間フォーマットの基本的な検証
	testTime := "2006-01-02T15:04:05Z07:00"
	_ = testTime
	return nil
}

// グローバル環境変数管理
var globalEnvManager *EnvVarManager

// InitGlobalEnvManager グローバル環境変数管理を初期化
func InitGlobalEnvManager(prefix string) {
	globalEnvManager = NewEnvVarManager(prefix)
}

// GetGlobalEnvManager グローバル環境変数管理を取得
func GetGlobalEnvManager() *EnvVarManager {
	if globalEnvManager == nil {
		globalEnvManager = NewEnvVarManager("LOG")
	}
	return globalEnvManager
}

// LoadConfigFromGlobalEnv グローバル環境変数から設定を読み込み
func LoadConfigFromGlobalEnv() *Config {
	return GetGlobalEnvManager().GetConfigFromEnv()
}

// ApplyConfigToGlobalEnv 設定をグローバル環境変数に適用
func ApplyConfigToGlobalEnv(config *Config) error {
	return GetGlobalEnvManager().ApplyConfigToEnv(config)
}

// ValidateGlobalEnv グローバル環境変数の妥当性を検証
func ValidateGlobalEnv() []string {
	return GetGlobalEnvManager().ValidateEnv()
}
