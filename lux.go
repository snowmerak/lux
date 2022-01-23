package lux

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/handler"
	"github.com/snowmerak/lux/logext"
	"github.com/snowmerak/lux/logext/stdout"
	"github.com/snowmerak/lux/middleware"
	"github.com/snowmerak/lux/router"
	"golang.org/x/net/http2"
)

type Lux struct {
	routers       []*router.RouterGroup
	logger        *logext.Logger
	server        *http.Server
	middlewares   []middleware.Set
	buildedRouter *httprouter.Router
}

func New(middlewares ...middleware.Set) *Lux {
	return &Lux{
		routers:       []*router.RouterGroup{},
		logger:        logext.New(stdout.New(8)),
		server:        new(http.Server),
		middlewares:   middlewares,
		buildedRouter: httprouter.New(),
	}
}

func (l *Lux) SetLogger(logger *logext.Logger) {
	l.logger = logger
}

func (l *Lux) SetReadHeaderTimeout(duration time.Duration) {
	l.server.ReadHeaderTimeout = duration
}

func (l *Lux) SetReadTimeout(duration time.Duration) {
	l.server.ReadTimeout = duration
}

func (l *Lux) SetWriteTimeout(duration time.Duration) {
	l.server.WriteTimeout = duration
}

func (l *Lux) SetIdleTimeout(duration time.Duration) {
	l.server.IdleTimeout = duration
}

func (l *Lux) SetMaxHeaderBytes(n int) {
	l.server.MaxHeaderBytes = n
}

func (l *Lux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	luxCtx := new(context.LuxContext)
	luxCtx.Request = r
	luxCtx.Response = context.NewResponse()
	defer func() {
		for key, values := range luxCtx.Response.Headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(luxCtx.Response.StatusCode)
		w.Write(luxCtx.Response.Body)
	}()
	if rs := middleware.ApplyRequests(luxCtx, l.middlewares); rs != "" {
		l.logger.Warnf(rs)
		return
	}
	l.buildedRouter.ServeHTTP(luxCtx.Response, luxCtx.Request)
	if !luxCtx.IsOk() {
		return
	}
	if rs := middleware.ApplyResponses(luxCtx, l.middlewares); rs != "" {
		l.logger.Warnf(rs)
		return
	}
}

func (l *Lux) buildServer(addr string) {
	l.server.Addr = addr
	l.server.Handler = l
	l.buildedRouter = new(httprouter.Router)
	for _, routerGroup := range l.routers {
		for path, routerMap := range routerGroup.Routers {
			for method, router := range routerMap {
				l.buildedRouter.Handle(method, path, handler.Wrap(l.logger, router.Handler))
			}
		}
	}
	l.routers = nil
}

func (l *Lux) ListenAndServeHTTP(addr string) error {
	l.buildServer(addr)
	if err := l.server.ListenAndServe(); err != nil {
		l.logger.Fatalf("ListenAndServeHTTP: %s", err)
		return err
	}
	return nil
}

func (l *Lux) ListenAndServeHTTPS(addr string, certFile string, keyFile string) error {
	l.buildServer(addr)
	if err := l.server.ListenAndServeTLS(certFile, keyFile); err != nil {
		l.logger.Fatalf("ListenAndServeHTTPS: %s", err)
		return err
	}
	return nil
}

func (l *Lux) ListenAndServeHTTP2(addr string) error {
	l.buildServer(addr)
	if err := http2.ConfigureServer(l.server, nil); err != nil {
		l.logger.Fatalf("ListenAndServeHTTP2: %s", err)
		return err
	}
	if err := l.server.ListenAndServe(); err != nil {
		l.logger.Fatalf("ListenAndServeHTTP2: %s", err)
		return err
	}
	return nil
}
