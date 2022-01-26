package lux

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/julienschmidt/httprouter"
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/handler"
	"github.com/snowmerak/lux/logext"
	"github.com/snowmerak/lux/logext/stdout"
	"github.com/snowmerak/lux/middleware"
	"github.com/snowmerak/lux/router"
	"github.com/snowmerak/lux/swagger"
	"github.com/snowmerak/lux/util"
	"golang.org/x/net/http2"
)

type Lux struct {
	routers       []*router.RouterGroup
	logger        *logext.Logger
	server        *http.Server
	middlewares   []middleware.Set
	buildedRouter *httprouter.Router
	swagger       *swagger.Swagger
}

func New(swaggerInfo *swagger.Info, middlewares ...middleware.Set) *Lux {
	swg := new(swagger.Swagger)
	if swaggerInfo != nil {
		swg.Info = *swaggerInfo
	}
	swg.SwaggerVersion = "2.0"
	return &Lux{
		routers:       []*router.RouterGroup{},
		logger:        logext.New(stdout.New(8)),
		server:        new(http.Server),
		middlewares:   middlewares,
		buildedRouter: httprouter.New(),
		swagger:       swg,
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

func (l *Lux) SetInfoEmail(email string) {
	l.swagger.Info.Contact.Email = email
}

func (l *Lux) SetInfoLicense(name, link string) {
	l.swagger.Info.License.Name = name
	l.swagger.Info.License.URL = link
}

func (l *Lux) ShowSwagger(path string, middlewares ...middleware.Set) {
	swaggerjson, err := json.Marshal(l.swagger)
	if err != nil {
		panic(err)
	}
	swagger.Dist["swagger.json"] = swaggerjson

	rg := l.NewRouterGroup(path, middlewares...)
	rg.GET("/*filepath", func(lc *context.LuxContext) error {
		filename := lc.GetPathVariable("filepath")
		filename = strings.TrimPrefix(filename, "/")
		if filename == "" {
			filename = "index.html"
		}
		if _, ok := swagger.Dist[filename]; !ok {
			lc.SetBadRequest()
			return nil
		}
		lc.Response.Headers.Set("Content-Type", util.GetContentTypeFromExt(filepath.Ext(filename)))
		lc.Response.Body = swagger.Dist[filename]
		lc.SetOK()
		return nil
	}, nil)

	l.logger.Infof("Swagger is available at %s/", path)
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
	l.logger.Infof("Server is ready to serve")
	l.logger.Infof("listen and serve on %s\n", addr)
}

func (l *Lux) ListenAndServe1(addr string) error {
	l.buildServer(addr)
	if err := l.server.ListenAndServe(); err != nil {
		l.logger.Fatalf("ListenAndServeHTTP: %s", err)
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe1TLS(addr string, certFile string, keyFile string) error {
	l.buildServer(addr)
	if err := l.server.ListenAndServeTLS(certFile, keyFile); err != nil {
		l.logger.Fatalf("ListenAndServeHTTPS: %s", err)
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe1AutoTLS(addr []string) error {
	if len(addr) == 0 {
		addr = []string{"localhost:443"}
	}
	l.buildServer(addr[0])
	if err := certmagic.HTTPS(addr, l.buildedRouter); err != nil {
		l.logger.Fatalf("ListenAndServeAutoHTTPS: %s", err)
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe2(addr string) error {
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

func (l *Lux) ListenAndServe2TLS(addr string, certFile string, keyFile string) error {
	l.buildServer(addr)
	if err := http2.ConfigureServer(l.server, nil); err != nil {
		l.logger.Fatalf("ListenAndServeHTTPS2: %s", err)
		return err
	}
	if err := l.server.ListenAndServeTLS(certFile, keyFile); err != nil {
		l.logger.Fatalf("ListenAndServeHTTPS2: %s", err)
		return err
	}
	return nil
}

func (l *Lux) ListenAndServe2AutoTLS(addr []string) error {
	if len(addr) == 0 {
		addr = []string{"localhost:443"}
	}
	l.buildServer(addr[0])
	if err := http2.ConfigureServer(l.server, nil); err != nil {
		l.logger.Fatalf("ListenAndServeAutoHTTPS2: %s", err)
		return err
	}
	if err := certmagic.HTTPS(addr, l.buildedRouter); err != nil {
		l.logger.Fatalf("ListenAndServeAutoHTTPS2: %s", err)
		return err
	}
	return nil
}
