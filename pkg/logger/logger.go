package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel ログレベル
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String ログレベルの文字列表現を返す
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLogLevel 文字列からログレベルを解析
func ParseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// LogEntry ログエントリの構造体
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Caller    string                 `json:"caller,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
}

// Logger ロガーの構造体
type Logger struct {
	level     LogLevel
	output    io.Writer
	fields    map[string]interface{}
	caller    bool
	traceID   string
	userID    string
	requestID string
}

// NewLogger 新しいロガーを作成
func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		output: output,
		fields: make(map[string]interface{}),
		caller: true,
	}
}

// WithField フィールドを追加したロガーを返す
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := l.clone()
	newLogger.fields[key] = value
	return newLogger
}

// WithFields 複数のフィールドを追加したロガーを返す
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := l.clone()
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

// WithTraceID トレースIDを設定したロガーを返す
func (l *Logger) WithTraceID(traceID string) *Logger {
	newLogger := l.clone()
	newLogger.traceID = traceID
	return newLogger
}

// WithUserID ユーザーIDを設定したロガーを返す
func (l *Logger) WithUserID(userID string) *Logger {
	newLogger := l.clone()
	newLogger.userID = userID
	return newLogger
}

// WithRequestID リクエストIDを設定したロガーを返す
func (l *Logger) WithRequestID(requestID string) *Logger {
	newLogger := l.clone()
	newLogger.requestID = requestID
	return newLogger
}

// clone ロガーのコピーを作成
func (l *Logger) clone() *Logger {
	fields := make(map[string]interface{})
	for k, v := range l.fields {
		fields[k] = v
	}
	return &Logger{
		level:     l.level,
		output:    l.output,
		fields:    fields,
		caller:    l.caller,
		traceID:   l.traceID,
		userID:    l.userID,
		requestID: l.requestID,
	}
}

// getCaller 呼び出し元の情報を取得
func (l *Logger) getCaller() string {
	if !l.caller {
		return ""
	}
	
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return ""
	}
	
	// ファイルパスを短縮
	parts := strings.Split(file, "/")
	if len(parts) > 2 {
		file = strings.Join(parts[len(parts)-2:], "/")
	}
	
	return fmt.Sprintf("%s:%d", file, line)
}

// log ログを出力
func (l *Logger) log(level LogLevel, message string, fields ...map[string]interface{}) {
	if level < l.level {
		return
	}
	
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Fields:    make(map[string]interface{}),
		Caller:    l.getCaller(),
		TraceID:   l.traceID,
		UserID:    l.userID,
		RequestID: l.requestID,
	}
	
	// 基本フィールドをコピー
	for k, v := range l.fields {
		entry.Fields[k] = v
	}
	
	// 追加フィールドをマージ
	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			entry.Fields[k] = v
		}
	}
	
	// JSON形式で出力
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// JSON化に失敗した場合はフォールバック
		fmt.Fprintf(l.output, "%s [%s] %s: %v\n", 
			entry.Timestamp.Format(time.RFC3339),
			entry.Level.String(),
			entry.Message,
			entry.Fields)
		return
	}
	
	fmt.Fprintln(l.output, string(jsonData))
}

// Debug デバッグログを出力
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	l.log(DEBUG, message, fields...)
}

// Info 情報ログを出力
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	l.log(INFO, message, fields...)
}

// Warn 警告ログを出力
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	l.log(WARN, message, fields...)
}

// Error エラーログを出力
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	l.log(ERROR, message, fields...)
}

// Fatal 致命的エラーログを出力して終了
func (l *Logger) Fatal(message string, fields ...map[string]interface{}) {
	l.log(FATAL, message, fields...)
	os.Exit(1)
}

// SetLevel ログレベルを設定
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput 出力先を設定
func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

// SetCaller 呼び出し元情報の出力を設定
func (l *Logger) SetCaller(enable bool) {
	l.caller = enable
}

// グローバルロガー
var globalLogger = NewLogger(INFO, os.Stdout)

// SetGlobalLogger グローバルロガーを設定
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// GetGlobalLogger グローバルロガーを取得
func GetGlobalLogger() *Logger {
	return globalLogger
}

// グローバルログ関数
func Debug(message string, fields ...map[string]interface{}) {
	globalLogger.Debug(message, fields...)
}

func Info(message string, fields ...map[string]interface{}) {
	globalLogger.Info(message, fields...)
}

func Warn(message string, fields ...map[string]interface{}) {
	globalLogger.Warn(message, fields...)
}

func Error(message string, fields ...map[string]interface{}) {
	globalLogger.Error(message, fields...)
}

func Fatal(message string, fields ...map[string]interface{}) {
	globalLogger.Fatal(message, fields...)
}

// WithField グローバルロガーにフィールドを追加
func WithField(key string, value interface{}) *Logger {
	return globalLogger.WithField(key, value)
}

// WithFields グローバルロガーに複数フィールドを追加
func WithFields(fields map[string]interface{}) *Logger {
	return globalLogger.WithFields(fields)
}

// WithTraceID グローバルロガーにトレースIDを追加
func WithTraceID(traceID string) *Logger {
	return globalLogger.WithTraceID(traceID)
}

// WithUserID グローバルロガーにユーザーIDを追加
func WithUserID(userID string) *Logger {
	return globalLogger.WithUserID(userID)
}

// WithRequestID グローバルロガーにリクエストIDを追加
func WithRequestID(requestID string) *Logger {
	return globalLogger.WithRequestID(requestID)
}
