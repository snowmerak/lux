package logger

import (
	"context"
	"fmt"

	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/logbuffer/logring"
	"github.com/snowmerak/logstream/log/logbuffer/logstream/globalque"
	"github.com/snowmerak/logstream/log/loglevel"
	"github.com/snowmerak/logstream/log/writable"
	"github.com/snowmerak/logstream/log/writable/stdout"
)

type Logger struct{}

const SYSTEM = "SYSTEM"
const MIDDLEWARE = "MIDDLEWARE"

var ctx = context.Background()
var globalQueue = globalque.New(ctx, logring.New, 16)

func Write(topic string, log log.Log) error {
	if err := globalQueue.Write(topic, log); err != nil {
		return fmt.Errorf("log.Write: %w", err)
	}
	return nil
}

func Observe(topic string, writers ...writable.Writable) error {
	if err := globalQueue.ObserveTopic(topic, writers...); err != nil {
		return fmt.Errorf("log.Observe: %w", err)
	}
	return nil
}

func init() {
	Observe(SYSTEM, stdout.New(context.Background(), loglevel.All, nil))
	Observe(MIDDLEWARE, stdout.New(context.Background(), loglevel.All, nil))
}

func (l Logger) WriteLog(topic string, log log.Log) error {
	return Write(topic, log)
}

func (l Logger) Observe(topic string, writers ...writable.Writable) error {
	return Observe(topic, writers...)
}

func (l Logger) Write(topic string, level loglevel.LogLevel, message string) error {
	return Write(topic, log.New(level, message).End())
}

func (l Logger) Printf(format string, args ...interface{}) {
	Write(SYSTEM, log.New(loglevel.Error, fmt.Sprintf(format, args...)).End())
}
