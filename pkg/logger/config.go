package logger

import (
	"io"
	"os"
	"strings"
)

// Config ログ設定の構造体
type Config struct {
	Level      LogLevel `json:"level"`
	Output     string   `json:"output"`
	Caller     bool     `json:"caller"`
	Pretty     bool     `json:"pretty"`
	TimeFormat string   `json:"time_format"`
}

// DefaultConfig デフォルト設定を返す
func DefaultConfig() *Config {
	return &Config{
		Level:      INFO,
		Output:     "stdout",
		Caller:     true,
		Pretty:     false,
		TimeFormat: "2006-01-02T15:04:05Z07:00",
	}
}

// ParseConfig 環境変数から設定を解析
func ParseConfig() *Config {
	config := DefaultConfig()

	// ログレベル
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Level = ParseLogLevel(level)
	}

	// 出力先
	if output := os.Getenv("LOG_OUTPUT"); output != "" {
		config.Output = output
	}

	// 呼び出し元情報
	if caller := os.Getenv("LOG_CALLER"); caller != "" {
		config.Caller = strings.ToLower(caller) == "true"
	}

	// プリティプリント
	if pretty := os.Getenv("LOG_PRETTY"); pretty != "" {
		config.Pretty = strings.ToLower(pretty) == "true"
	}

	// 時間フォーマット
	if timeFormat := os.Getenv("LOG_TIME_FORMAT"); timeFormat != "" {
		config.TimeFormat = timeFormat
	}

	return config
}

// GetOutput 設定に基づいて出力先を取得
func (c *Config) GetOutput() (io.Writer, error) {
	switch strings.ToLower(c.Output) {
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	case "null":
		return os.NewFile(0, os.DevNull), nil
	default:
		// ファイル出力
		file, err := os.OpenFile(c.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

// Apply 設定をロガーに適用
func (c *Config) Apply(logger *Logger) error {
	output, err := c.GetOutput()
	if err != nil {
		return err
	}

	logger.SetLevel(c.Level)
	logger.SetOutput(output)
	logger.SetCaller(c.Caller)

	return nil
}

// NewLoggerFromConfig 設定からロガーを作成
func NewLoggerFromConfig(config *Config) (*Logger, error) {
	output, err := config.GetOutput()
	if err != nil {
		return nil, err
	}

	logger := NewLogger(config.Level, output)
	logger.SetCaller(config.Caller)

	return logger, nil
}

// InitGlobalLogger グローバルロガーを初期化
func InitGlobalLogger(config *Config) error {
	logger, err := NewLoggerFromConfig(config)
	if err != nil {
		return err
	}

	SetGlobalLogger(logger)
	return nil
}

// InitFromEnv 環境変数からグローバルロガーを初期化
func InitFromEnv() error {
	config := ParseConfig()
	return InitGlobalLogger(config)
}
