package logger

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type StdLogger struct {
	level LogLevel
	log   *log.Logger
}

func NewStdLogger(level LogLevel, output io.Writer) *StdLogger {
	return &StdLogger{
		level: level,
		log:   log.New(output, "", 0),
	}
}

func (sl *StdLogger) Debug(msg string, fields ...any) {
	if DEBUG >= sl.level {
		sl.log.Printf(formatLog(sl.level, "DEBUG", msg, fields...))
	}
}

func (sl *StdLogger) Info(msg string, fields ...any) {
	if INFO >= sl.level {
		sl.log.Printf(formatLog(sl.level, "INFO", msg, fields...))
	}
}

func (sl *StdLogger) Warn(msg string, fields ...any) {
	if WARN >= sl.level {
		sl.log.Printf(formatLog(sl.level, "WARN", msg, fields...))
	}
}

func (sl *StdLogger) Error(msg string, fields ...any) {
	if ERROR >= sl.level {
		sl.log.Printf(formatLog(sl.level, "ERROR", msg, fields...))
	}
}

func formatLog(level LogLevel, callLevel string, msg string, fields ...any) string {
	var logMsg strings.Builder

	appendTimestamp(&logMsg)
	if level == DEBUG {
		appendDebugInfo(&logMsg)
	}
	appendMessage(&logMsg, callLevel, msg)
	appendFields(&logMsg, fields...)

	return logMsg.String()
}

func appendDebugInfo(msg *strings.Builder) {
	// Skipping levels up the call stack until we get the original call
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file, line = "unknown", 0
	} else {
		file = truncatePathToDepth(file, 2)
	}
	msg.WriteString(fmt.Sprintf("(%s:%d) ", file, line))
}

func truncatePathToDepth(path string, numParents int) string {
	path = filepath.Clean(path)
	dirCount := 0

	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == filepath.Separator {
			dirCount++
		}
		if dirCount == numParents+1 {
			return path[i+1:]
		}
	}
	// Return full path as numParents is greater than actual number of parents
	return path
}

func appendTimestamp(msg *strings.Builder) {
	isoTimestamp := time.Now().Format(time.RFC3339)
	msg.WriteString(fmt.Sprintf("%v ", isoTimestamp))
}

func appendMessage(msg *strings.Builder, callLevel string, s string) {
	msg.WriteString(fmt.Sprintf("[%s] %s", callLevel, s))
}

func appendFields(msg *strings.Builder, fields ...any) {
	if len(fields) < 1 {
		return
	}
	if len(fields)%2 != 0 {
		msg.WriteString(" {LOG_ERROR: fields are key-value pairs; must be even number}")
		return
	}
	invalidKey := false
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			key = "INVALID_KEY"
			invalidKey = true
		}
		msg.WriteString(fmt.Sprintf("\n  %s: '%v'", key, fields[i+1]))
	}
	if invalidKey {
		msg.WriteString(" {LOG_ERROR: all keys must be strings}")
	}
}
