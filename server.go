package lux

import (
	"time"

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

func (l *Lux) ListenAndServe(addr string) error {
	l.server.Handler = l.router.Handler
	return l.server.ListenAndServe(addr)
}
