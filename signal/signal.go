package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Terminate() <-chan struct{} {
	sigs := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- struct{}{}
	}()

	return done
}
