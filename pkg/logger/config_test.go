package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	assert.Equal(t, INFO, config.Level)
	assert.Equal(t, "stdout", config.Output)
	assert.True(t, config.Caller)
	assert.False(t, config.Pretty)
	assert.Equal(t, "2006-01-02T15:04:05Z07:00", config.TimeFormat)
}

func TestParseConfig(t *testing.T) {
	// 環境変数をクリア
	os.Clearenv()
	
	// デフォルト設定のテスト
	config := ParseConfig()
	assert.Equal(t, INFO, config.Level)
	assert.Equal(t, "stdout", config.Output)
	assert.True(t, config.Caller)
	assert.False(t, config.Pretty)
}

func TestParseConfigWithEnvVars(t *testing.T) {
	// 環境変数をクリア
	os.Clearenv()
	
	// 環境変数を設定
	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("LOG_OUTPUT", "stderr")
	os.Setenv("LOG_CALLER", "false")
	os.Setenv("LOG_PRETTY", "true")
	os.Setenv("LOG_TIME_FORMAT", "2006-01-02 15:04:05")
	
	config := ParseConfig()
	
	assert.Equal(t, DEBUG, config.Level)
	assert.Equal(t, "stderr", config.Output)
	assert.False(t, config.Caller)
	assert.True(t, config.Pretty)
	assert.Equal(t, "2006-01-02 15:04:05", config.TimeFormat)
}

func TestGetOutput(t *testing.T) {
	config := DefaultConfig()
	
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{"stdout", "stdout", "stdout"},
		{"stderr", "stderr", "stderr"},
		{"null", "null", "null"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Output = tt.output
			writer, err := config.GetOutput()
			require.NoError(t, err)
			assert.NotNil(t, writer)
		})
	}
}

func TestGetOutputFile(t *testing.T) {
	config := DefaultConfig()
	config.Output = "/tmp/test.log"
	
	writer, err := config.GetOutput()
	require.NoError(t, err)
	assert.NotNil(t, writer)
	
	// ファイルをクリーンアップ
	os.Remove("/tmp/test.log")
}

func TestApply(t *testing.T) {
	config := DefaultConfig()
	config.Level = WARN
	config.Output = "stderr"
	config.Caller = false
	
	var buf strings.Builder
	logger := NewLogger(INFO, &buf)
	
	err := config.Apply(logger)
	require.NoError(t, err)
	
	assert.Equal(t, WARN, logger.level)
	assert.Equal(t, os.Stderr, logger.output)
	assert.False(t, logger.caller)
}

func TestNewLoggerFromConfig(t *testing.T) {
	config := DefaultConfig()
	config.Level = ERROR
	config.Output = "stderr"
	config.Caller = false
	
	logger, err := NewLoggerFromConfig(config)
	require.NoError(t, err)
	
	assert.Equal(t, ERROR, logger.level)
	assert.Equal(t, os.Stderr, logger.output)
	assert.False(t, logger.caller)
}

func TestInitGlobalLogger(t *testing.T) {
	config := DefaultConfig()
	config.Level = DEBUG
	config.Output = "stderr"
	
	err := InitGlobalLogger(config)
	require.NoError(t, err)
	
	globalLogger := GetGlobalLogger()
	assert.Equal(t, DEBUG, globalLogger.level)
	assert.Equal(t, os.Stderr, globalLogger.output)
}

func TestInitFromEnv(t *testing.T) {
	// 環境変数をクリア
	os.Clearenv()
	
	// 環境変数を設定
	os.Setenv("LOG_LEVEL", "WARN")
	os.Setenv("LOG_OUTPUT", "stderr")
	os.Setenv("LOG_CALLER", "false")
	
	err := InitFromEnv()
	require.NoError(t, err)
	
	globalLogger := GetGlobalLogger()
	assert.Equal(t, WARN, globalLogger.level)
	assert.Equal(t, os.Stderr, globalLogger.output)
	assert.False(t, globalLogger.caller)
}

func TestConfigValidation(t *testing.T) {
	config := DefaultConfig()
	
	// 無効な出力先
	config.Output = "/invalid/path/that/does/not/exist/test.log"
	
	_, err := config.GetOutput()
	assert.Error(t, err)
}

func TestConfigWithInvalidLogLevel(t *testing.T) {
	os.Clearenv()
	os.Setenv("LOG_LEVEL", "INVALID")
	
	config := ParseConfig()
	
	// 無効なログレベルはデフォルトのINFOになる
	assert.Equal(t, INFO, config.Level)
}

func TestConfigWithInvalidCaller(t *testing.T) {
	os.Clearenv()
	os.Setenv("LOG_CALLER", "invalid")
	
	config := ParseConfig()
	
	// 無効な値はデフォルトのtrueになる
	assert.False(t, config.Caller)
}

func TestConfigWithInvalidPretty(t *testing.T) {
	os.Clearenv()
	os.Setenv("LOG_PRETTY", "invalid")
	
	config := ParseConfig()
	
	// 無効な値はデフォルトのfalseになる
	assert.False(t, config.Pretty)
}

func TestConfigEmptyEnvVars(t *testing.T) {
	os.Clearenv()
	os.Setenv("LOG_LEVEL", "")
	os.Setenv("LOG_OUTPUT", "")
	os.Setenv("LOG_CALLER", "")
	os.Setenv("LOG_PRETTY", "")
	os.Setenv("LOG_TIME_FORMAT", "")
	
	config := ParseConfig()
	
	// 空の環境変数はデフォルト値になる
	assert.Equal(t, INFO, config.Level)
	assert.Equal(t, "stdout", config.Output)
	assert.True(t, config.Caller)
	assert.False(t, config.Pretty)
	assert.Equal(t, "2006-01-02T15:04:05Z07:00", config.TimeFormat)
}
