package lux

import (
	"context"

	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/loglevel"
	"github.com/snowmerak/logstream/log/writable"
	"github.com/snowmerak/logstream/log/writable/stdout"
	"github.com/snowmerak/lux/logger"
	"github.com/snowmerak/lux/middleware"
)

type Logger struct{}

func init() {
	logger.Observe(logger.SYSTEM, stdout.New(context.Background(), loglevel.All, nil))
	logger.Observe(middleware.MIDDLEWARE, stdout.New(context.Background(), loglevel.All, nil))
}

func (l Logger) Write(topic string, log log.Log) error {
	return logger.Write(topic, log)
}

func (l Logger) Observe(topic string, writers ...writable.Writable) error {
	return logger.Observe(topic, writers...)
}
