package logger

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected string
	}{
		{"DEBUG", DEBUG, "DEBUG"},
		{"INFO", INFO, "INFO"},
		{"WARN", WARN, "WARN"},
		{"ERROR", ERROR, "ERROR"},
		{"FATAL", FATAL, "FATAL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.String())
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected LogLevel
	}{
		{"DEBUG", "DEBUG", DEBUG},
		{"debug", "debug", DEBUG},
		{"INFO", "INFO", INFO},
		{"info", "info", INFO},
		{"WARN", "WARN", WARN},
		{"warn", "warn", WARN},
		{"ERROR", "ERROR", ERROR},
		{"error", "error", ERROR},
		{"FATAL", "FATAL", FATAL},
		{"fatal", "fatal", FATAL},
		{"UNKNOWN", "UNKNOWN", INFO},
		{"empty", "", INFO},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ParseLogLevel(tt.input))
		})
	}
}

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	assert.Equal(t, INFO, logger.level)
	assert.Equal(t, &buf, logger.output)
	assert.True(t, logger.caller)
	assert.NotNil(t, logger.fields)
}

func TestLoggerWithField(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	loggerWithField := logger.WithField("key", "value")

	assert.NotEqual(t, logger, loggerWithField)
	assert.Equal(t, "value", loggerWithField.fields["key"])
	assert.Empty(t, logger.fields)
}

func TestLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	loggerWithFields := logger.WithFields(fields)

	assert.NotEqual(t, logger, loggerWithFields)
	assert.Equal(t, "value1", loggerWithFields.fields["key1"])
	assert.Equal(t, "value2", loggerWithFields.fields["key2"])
	assert.Empty(t, logger.fields)
}

func TestLoggerWithTraceID(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	loggerWithTraceID := logger.WithTraceID("trace-123")

	assert.NotEqual(t, logger, loggerWithTraceID)
	assert.Equal(t, "trace-123", loggerWithTraceID.traceID)
	assert.Empty(t, logger.traceID)
}

func TestLoggerWithUserID(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	loggerWithUserID := logger.WithUserID("user-123")

	assert.NotEqual(t, logger, loggerWithUserID)
	assert.Equal(t, "user-123", loggerWithUserID.userID)
	assert.Empty(t, logger.userID)
}

func TestLoggerWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	loggerWithRequestID := logger.WithRequestID("req-123")

	assert.NotEqual(t, logger, loggerWithRequestID)
	assert.Equal(t, "req-123", loggerWithRequestID.requestID)
	assert.Empty(t, logger.requestID)
}

func TestLoggerLog(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(DEBUG, &buf)

	// デバッグログのテスト
	logger.Debug("debug message", map[string]interface{}{"key": "value"})

	var entry LogEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "debug message", entry.Message)
	assert.Equal(t, DEBUG, entry.Level)
	assert.Equal(t, "value", entry.Fields["key"])
	assert.NotEmpty(t, entry.Timestamp)
	assert.NotEmpty(t, entry.Caller)
}

func TestLoggerLogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	// デバッグログは出力されない
	logger.Debug("debug message")
	assert.Empty(t, buf.String())

	// 情報ログは出力される
	logger.Info("info message")
	assert.NotEmpty(t, buf.String())
}

func TestLoggerSetLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	logger.SetLevel(ERROR)

	// 情報ログは出力されない
	logger.Info("info message")
	assert.Empty(t, buf.String())

	// エラーログは出力される
	logger.Error("error message")
	assert.NotEmpty(t, buf.String())
}

func TestLoggerSetOutput(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	logger := NewLogger(INFO, &buf1)

	logger.SetOutput(&buf2)
	logger.Info("test message")

	assert.Empty(t, buf1.String())
	assert.NotEmpty(t, buf2.String())
}

func TestLoggerSetCaller(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	logger.SetCaller(false)
	logger.Info("test message")

	var entry LogEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Empty(t, entry.Caller)
}

func TestGlobalLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	SetGlobalLogger(logger)

	// グローバル関数のテスト
	Info("global info message")

	var entry LogEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "global info message", entry.Message)
	assert.Equal(t, INFO, entry.Level)
}

func TestGlobalLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)
	SetGlobalLogger(logger)

	// グローバル関数でのフィールド追加
	WithField("key", "value").Info("message with field")

	var entry LogEntry
	err := json.Unmarshal(buf.Bytes(), &entry)
	require.NoError(t, err)

	assert.Equal(t, "message with field", entry.Message)
	assert.Equal(t, "value", entry.Fields["key"])
}

func TestLogEntryJSON(t *testing.T) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     INFO,
		Message:   "test message",
		Fields: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
		Caller:    "test.go:10",
		TraceID:   "trace-123",
		UserID:    "user-456",
		RequestID: "req-789",
	}

	jsonData, err := json.Marshal(entry)
	require.NoError(t, err)

	var decoded LogEntry
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, entry.Message, decoded.Message)
	assert.Equal(t, entry.Level, decoded.Level)
	assert.Equal(t, entry.Fields["key1"], decoded.Fields["key1"])
	assert.Equal(t, float64(entry.Fields["key2"].(int)), decoded.Fields["key2"])
	assert.Equal(t, entry.Caller, decoded.Caller)
	assert.Equal(t, entry.TraceID, decoded.TraceID)
	assert.Equal(t, entry.UserID, decoded.UserID)
	assert.Equal(t, entry.RequestID, decoded.RequestID)
}

func TestLoggerFatal(t *testing.T) {
	// Fatalはos.Exit(1)を呼ぶので、テストでは実行しない
	// 実際のアプリケーションでは適切に動作する
	t.Skip("Fatal test skipped - would cause os.Exit(1)")
}

func TestLoggerGetCaller(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)

	caller := logger.getCaller()

	// 呼び出し元の情報が含まれていることを確認
	assert.NotEmpty(t, caller)
	assert.Contains(t, caller, ":")
}

func TestLoggerClone(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(INFO, &buf)
	logger.traceID = "trace-123"
	logger.userID = "user-456"
	logger.requestID = "req-789"
	logger.fields["key"] = "value"

	cloned := logger.clone()

	// 同じ値を持つが、異なるインスタンス
	assert.Equal(t, logger.level, cloned.level)
	assert.Equal(t, logger.output, cloned.output)
	assert.Equal(t, logger.caller, cloned.caller)
	assert.Equal(t, logger.traceID, cloned.traceID)
	assert.Equal(t, logger.userID, cloned.userID)
	assert.Equal(t, logger.requestID, cloned.requestID)
	assert.Equal(t, logger.fields["key"], cloned.fields["key"])

	// フィールドの変更が元に影響しない
	cloned.fields["new_key"] = "new_value"
	assert.NotContains(t, logger.fields, "new_key")
}
