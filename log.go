package lux

import (
	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/writable"
	"github.com/snowmerak/lux/logger"
)

type Logger struct{}

func (l Logger) Write(topic string, log log.Log) error {
	return logger.Write(topic, log)
}

func (l Logger) Observe(topic string, writers ...writable.Writable) error {
	return logger.Observe(topic, writers...)
}
