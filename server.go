package lux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/fasthttp/router"
	"github.com/snowmerak/lux/logger"
	"github.com/snowmerak/lux/middleware"
	"github.com/snowmerak/lux/swagger"
	"github.com/valyala/fasthttp"
)

/*
Lux ...
server is a fasthttp.Server
router is a router.Router
and has some global middlewares
*/
type Lux struct {
	server              *fasthttp.Server
	router              *router.Router
	requestMiddlewares  []middleware.Middleware
	responseMiddlewares []middleware.Middleware
	Swagger             Swagger
}

/*
NewServer ...
Create NewServer with default settings.
*/
func NewServer() *Lux {
	return &Lux{
		server: &fasthttp.Server{
			Logger: logger.Logger{},
		},
		router: router.New(),
		Swagger: Swagger{
			&swagger.Swagger{
				SwaggerVersion:      "2.0",
				Paths:               make(map[swagger.Path]map[swagger.Method]swagger.Router),
				SecurityDefinitions: make(map[string]swagger.SecurityDefinition),
				Definitions:         make(map[string]swagger.Definition),
			},
		},
	}
}

/*
RouterGroup ...
returns a RouterGroup of paths.
*/
func (l *Lux) RouterGroup(path ...string) *RouterGroup {
	name := "/" + strings.Join(path, "/")
	l.Swagger.inner.Paths[swagger.Path(name)] = make(map[swagger.Method]swagger.Router)
	group := l.router.Group(name)
	return &RouterGroup{
		group:               group,
		requestMiddlewares:  []middleware.Middleware{},
		responseMiddlewares: []middleware.Middleware{},

		swagger: l.Swagger.inner,
		path:    name,
	}
}

//SetRequestBodySize ...
//SetRequestBodySize sets the maximum request body size.
//unit: bytes
func (l *Lux) SetRequestBodySize(size int) {
	l.server.MaxRequestBodySize = size
}

/*
SetReadTimeout ...
SetReadTimeout sets the maximum duration for reading the entire
*/
func (l *Lux) SetReadTimeout(timeout time.Duration) {
	l.server.ReadTimeout = timeout
}

/*
SetWriteTimeout ...
SetWriteTimeout sets the maximum duration before timing out
*/
func (l *Lux) SetWriteTimeout(timeout time.Duration) {
	l.server.WriteTimeout = timeout
}

/*
SetConnPerIPLimit ...
SetConnPerIPLimit sets the maximum number of concurrent connections per IP.
*/
func (l *Lux) SetConnPerIPLimit(limit int) {
	l.server.MaxConnsPerIP = limit
}

/*
SetMaxRequestPerConn ...
SetMaxRequestPerConn sets the maximum number of requests per connection.
*/
func (l *Lux) SetMaxRequestsPerConn(limit int) {
	l.server.MaxRequestsPerConn = limit
}

/*
Use ...
Use adds a middleware to the server.
*/
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

/*
ListenAndServe ...
ListenAndServe starts the server over http.
*/
func (l *Lux) ListenAndServe(addr string) error {
	l.wrapHandler()
	printInfo(addr)
	return l.server.ListenAndServe(addr)
}

/*
ListenAndServeTLS ...
ListenAndServeTLS starts the server over https with cert file and key file.
*/
func (l *Lux) ListenAndServeTLS(addr, certFile, keyFile string) error {
	l.wrapHandler()
	printInfo(addr)
	return l.server.ListenAndServeTLS(addr, certFile, keyFile)
}

/*
ListenAndServeAutoTLSEmbed ...
ListenAndServeAutoTLSEmbed starts the server over https with embedded cert data and key data.
*/
func (l *Lux) ListenAndServeTLSEmbed(addr string, certData, keyData []byte) error {
	l.wrapHandler()
	printInfo(addr)
	return l.server.ListenAndServeTLSEmbed(addr, certData, keyData)
}

/*
ListenAndServeUNIX ...
ListenAndServeUNIX starts the server over unix socket.
*/
func (l *Lux) ListenAndServeUNIX(addr string, mode os.FileMode) error {
	l.wrapHandler()
	printInfo(addr)
	return l.server.ListenAndServeUNIX(addr, mode)
}

/*
ListenAndServeAutoTLS ...
ListenAndServeAutoTLS starts the server over https with cert file and key file automatic generated.
*/
func (l *Lux) ListenAndServeAutoTLS(addr string) error {
	ln, err := certmagic.Listen([]string{addr})
	if err != nil {
		return err
	}
	l.wrapHandler()
	printInfo(addr)
	return l.server.Serve(ln)
}

func (l *Lux) ShowSwagger(path string) {
	if path == "" {
		path = "/"
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if !strings.HasSuffix(path, "{filepath:*}") {
		path += "{filepath:*}"
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(l.Swagger.inner); err != nil {
		panic(err)
	}

	f, err := os.Create("./swagger/dist/swagger.json")
	if err != nil {
		panic(err)
	}

	f.Write(buf.Bytes())
	f.Close()

	l.router.ServeFiles(path, "./swagger/dist")
}
