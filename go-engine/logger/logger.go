package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

// Logger provides structured logging
type Logger struct {
	*slog.Logger
	level  string
	format string
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp  time.Time              `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	Component  string                 `json:"component,omitempty"`
	Action     string                 `json:"action,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Error      string                 `json:"error,omitempty"`
	DurationMs int64                  `json:"duration_ms,omitempty"`
}

// New creates a new Logger instance
func New(level, format, filePath string) (*Logger, error) {
	var writer io.Writer = os.Stdout

	// Add file output if specified
	if filePath != "" {
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writer = io.MultiWriter(os.Stdout, file)
	}

	// Set log level
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
		level:  level,
		format: format,
	}, nil
}

// WithComponent returns a logger with a component context
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		Logger: l.Logger.With("component", component),
		level:  l.level,
		format: l.format,
	}
}

// LogAction logs an action with details
func (l *Logger) LogAction(action string, details map[string]interface{}) {
	l.Info(action, "details", details)
}

// LogError logs an error with context
func (l *Logger) LogError(action string, err error, details map[string]interface{}) {
	if details == nil {
		details = make(map[string]interface{})
	}
	details["error"] = err.Error()
	l.Error(action, "details", details)
}

// LogTiming logs an action with duration
func (l *Logger) LogTiming(action string, startTime time.Time, details map[string]interface{}) {
	if details == nil {
		details = make(map[string]interface{})
	}
	details["duration_ms"] = time.Since(startTime).Milliseconds()
	l.Info(action, "details", details)
}

// ToJSON converts log entry to JSON string
func (e *LogEntry) ToJSON() string {
	data, _ := json.Marshal(e)
	return string(data)
}
