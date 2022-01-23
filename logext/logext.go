package logext

import (
	"fmt"
	"io"
	"time"
)

type Logger struct {
	writer io.Writer
}

func New(writer io.Writer) *Logger {
	return &Logger{
		writer: writer,
	}
}

func (l *Logger) Write(msg string) (int, error) {
	return l.writer.Write([]byte(msg))
}

func (l *Logger) Writef(format string, args ...interface{}) (int, error) {
	now := time.Now().Format(time.RFC3339) + " "
	return fmt.Fprintf(l.writer, fmt.Sprintf(now+format, args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) (int, error) {
	level := "[DEBUG] "
	return l.Writef(level+format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) (int, error) {
	level := "[INFO] "
	return l.Writef(level+format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) (int, error) {
	level := "[WARN] "
	return l.Writef(level+format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) (int, error) {
	level := "[ERROR] "
	return l.Writef(level+format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) (int, error) {
	level := "[FATAL] "
	return l.Writef(level+format, args...)
}
