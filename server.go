package lux

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/fasthttp/router"
	"github.com/snowmerak/lux/logger"
	"github.com/snowmerak/lux/middleware"
	"github.com/valyala/fasthttp"
)

type Lux struct {
	server              *fasthttp.Server
	router              *router.Router
	requestMiddlewares  []middleware.Middleware
	responseMiddlewares []middleware.Middleware
}

func NewServer() *Lux {
	return &Lux{
		server: &fasthttp.Server{
			Logger: logger.Logger{},
		},
		router: router.New(),
	}
}

func (l *Lux) RouterGroup(path ...string) *RouterGroup {
	group := l.router.Group("/" + strings.Join(path, "/"))
	return &RouterGroup{
		group:               group,
		requestMiddlewares:  []middleware.Middleware{},
		responseMiddlewares: []middleware.Middleware{},
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

func (l *Lux) Use(middlewareset ...middleware.MiddlewareSet) {
	for _, m := range middlewareset {
		req, res := m.Request, m.Response
		if req != nil {
			l.requestMiddlewares = append(l.requestMiddlewares, req)
		}
		if res != nil {
			l.responseMiddlewares = append(l.responseMiddlewares, res)
		}
	}
}

func printInfo(addr string) {
	fmt.Print(banner)
	fmt.Printf("Lux is running at %s\n", addr)
}

func (l *Lux) wrapHandler() {
	handler := func(ctx *fasthttp.RequestCtx) {
		luxCtx := &LuxContext{
			ctx: ctx,
		}
		for _, m := range l.requestMiddlewares {
			if !luxCtx.Ok() {
				return
			}
			m(luxCtx.ctx)
		}
		l.router.Handler(ctx)
		for _, m := range l.responseMiddlewares {
			if !luxCtx.Ok() {
				return
			}
			m(luxCtx.ctx)
		}
	}
	l.server.Handler = handler
}

func (l *Lux) ListenAndServe(addr string) error {
	l.wrapHandler()
	printInfo(addr)
	return l.server.ListenAndServe(addr)
}

func (l *Lux) ListenAndServeTLS(addr, certFile, keyFile string) error {
	l.wrapHandler()
	printInfo(addr)
	return l.server.ListenAndServeTLS(addr, certFile, keyFile)
}

func (l *Lux) ListenAndServeTLSEmbed(addr string, certData, keyData []byte) error {
	l.wrapHandler()
	printInfo(addr)
	return l.server.ListenAndServeTLSEmbed(addr, certData, keyData)
}

func (l *Lux) ListenAndServeUNIX(addr string, mode os.FileMode) error {
	l.wrapHandler()
	printInfo(addr)
	return l.server.ListenAndServeUNIX(addr, mode)
}

func (l *Lux) ListenAndServeAutoTLS(addr string) error {
	ln, err := certmagic.Listen([]string{addr})
	if err != nil {
		return err
	}
	l.wrapHandler()
	printInfo(addr)
	return l.server.Serve(ln)
}
