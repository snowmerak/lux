package stdout

import (
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
	stdout := &Stdout{
		writer: os.Stdout,
		buffer: queue.NewRingBuffer(uint64(bufSize)),
		signal: make(chan struct{}, uint64(bufSize)),
	}
	go func() {
		for range stdout.signal {
			v, err := stdout.buffer.Get()
			if err != nil {
				return
			}
			stdout.writer.Write(append(v.([]byte), '\n'))
		}
	}()
	return stdout
}

func (s *Stdout) Dispose() {
	close(s.signal)
	s.buffer.Dispose()
}

func (s *Stdout) Write(p []byte) (n int, err error) {
	s.buffer.Put(p)
	s.signal <- struct{}{}
	return len(p), nil
}
