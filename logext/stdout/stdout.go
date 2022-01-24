package stdout

import (
	"fmt"
	"io"
	"os"

	"github.com/Workiva/go-datastructures/queue"
)

type Stdout struct {
	writer io.Writer
	buffer *queue.RingBuffer
	signal chan struct{}
}

func New(bufSize int) *Stdout {
	writer := os.Stdout
	stdout := &Stdout{
		writer: writer,
		buffer: queue.NewRingBuffer(uint64(bufSize)),
		signal: make(chan struct{}, uint64(bufSize)),
	}
	go func() {
		for range stdout.signal {
			v, err := stdout.buffer.Get()
			if err != nil {
				return
			}
			fmt.Fprintln(writer, v.(string))
		}
	}()
	return stdout
}

func (s *Stdout) Dispose() {
	close(s.signal)
	s.buffer.Dispose()
}

func (s *Stdout) Write(p []byte) (n int, err error) {
	s.buffer.Put(string(p))
	s.signal <- struct{}{}
	return len(p), nil
}
