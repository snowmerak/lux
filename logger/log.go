package logger

import (
	"context"
	"fmt"

	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/logbuffer/logring"
	"github.com/snowmerak/logstream/log/logbuffer/logstream/globalque"
	"github.com/snowmerak/logstream/log/writable"
)

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
