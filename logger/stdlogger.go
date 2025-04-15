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
		sl.log.Print(formatLog(sl.level, "DEBUG", msg, fields...))
	}
}

func (sl *StdLogger) Info(msg string, fields ...any) {
	if INFO >= sl.level {
		sl.log.Print(formatLog(sl.level, "INFO", msg, fields...))
	}
}

func (sl *StdLogger) Warn(msg string, fields ...any) {
	if WARN >= sl.level {
		sl.log.Print(formatLog(sl.level, "WARN", msg, fields...))
	}
}

func (sl *StdLogger) Error(msg string, fields ...any) {
	if ERROR >= sl.level {
		sl.log.Print(formatLog(sl.level, "ERROR", msg, fields...))
	}
}

func (sl *StdLogger) Fatal(msg string, fields ...any) {
	if FATAL >= sl.level {
		sl.log.Fatal(formatLog(sl.level, "FATAL", msg, fields...))
	}
}

func (sl *StdLogger) Printf(msg string, fields ...any) {
  fmt.Printf(msg, fields...)
}

func (sl *StdLogger) Verbose() bool {
  return sl.level >= DEBUG
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
	fmt.Fprintf(msg, "(%s:%d) ", file, line)
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
	fmt.Fprintf(msg, "%v ", isoTimestamp)
}

func appendMessage(msg *strings.Builder, callLevel string, s string) {
	fmt.Fprintf(msg, "[%s] %s", callLevel, s)
}

func appendFields(msg *strings.Builder, fields ...any) {
	if len(fields) < 1 {
		return
	}
	if len(fields)%2 != 0 {
		msg.WriteString("\n{LOG_ERROR: fields must be key-value pairs")
		return
	}
	invalidKey := false
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			key = "INVALID_KEY"
			invalidKey = true
		}
		fmt.Fprintf(msg, "\n%s: '%v'", key, fields[i+1])
	}
	if invalidKey {
		msg.WriteString("\n{LOG_ERROR: all keys must be strings}")
	}
}
