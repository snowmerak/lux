package lux

import (
	"fmt"
	"os"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/fasthttp/router"
	"github.com/snowmerak/lux/logger"
	"github.com/valyala/fasthttp"
)

type Lux struct {
	server *fasthttp.Server
	router *router.Router
}

func NewServer() *Lux {
	return &Lux{
		server: &fasthttp.Server{
			Logger: logger.Logger{},
		},
		router: router.New(),
	}
}

//SetRequestBodySize sets the maximum request body size.
//unit: bytes
func (l *Lux) SetRequestBodySize(size int) {
	l.server.MaxRequestBodySize = size
}

func (l *Lux) SetReadTimeout(timeout time.Duration) {
	l.server.ReadTimeout = timeout
}

func (l *Lux) SetWriteTimeout(timeout time.Duration) {
	l.server.WriteTimeout = timeout
}

func (l *Lux) SetConnPerIPLimit(limit int) {
	l.server.MaxConnsPerIP = limit
}

func (l *Lux) SetMaxRequestsPerConn(limit int) {
	l.server.MaxRequestsPerConn = limit
}

func printInfo(addr string) {
	fmt.Print(banner)
	fmt.Printf("Lux is running at %s\n", addr)
}

func (l *Lux) ListenAndServe(addr string) error {
	l.server.Handler = l.router.Handler
	printInfo(addr)
	return l.server.ListenAndServe(addr)
}

func (l *Lux) ListenAndServeTLS(addr, certFile, keyFile string) error {
	l.server.Handler = l.router.Handler
	printInfo(addr)
	return l.server.ListenAndServeTLS(addr, certFile, keyFile)
}

func (l *Lux) ListenAndServeTLSEmbed(addr string, certData, keyData []byte) error {
	l.server.Handler = l.router.Handler
	printInfo(addr)
	return l.server.ListenAndServeTLSEmbed(addr, certData, keyData)
}

func (l *Lux) ListenAndServeUNIX(addr string, mode os.FileMode) error {
	l.server.Handler = l.router.Handler
	printInfo(addr)
	return l.server.ListenAndServeUNIX(addr, mode)
}

func (l *Lux) ListenAndServeAutoTLS(addr string) error {
	ln, err := certmagic.Listen([]string{addr})
	if err != nil {
		return err
	}
	l.server.Handler = l.router.Handler
	printInfo(addr)
	return l.server.Serve(ln)
}
