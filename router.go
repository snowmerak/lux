package lux

import (
	"strings"

	"github.com/fasthttp/router"
	"github.com/snowmerak/lux/middleware"
	"github.com/valyala/fasthttp"
)

const GET = "GET"
const POST = "POST"
const HEAD = "HEAD"
const DELETE = "DELETE"
const PUT = "PUT"
const PATCH = "PATCH"
const OPTIONS = "OPTIONS"

type RouterGroup struct {
	group              *router.Group
	requestMiddlewares []middleware.Middleware
	responseMiddleware []middleware.Middleware
}

func (l *Lux) RouterGroup(path ...string) *RouterGroup {
	group := l.router.Group("/" + strings.Join(path, "/"))
	return &RouterGroup{
		group:              group,
		requestMiddlewares: []middleware.Middleware{},
	}
}

func (r *RouterGroup) UseRequest(m ...middleware.Middleware) *RouterGroup {
	r.requestMiddlewares = append(r.requestMiddlewares, m...)
	return r
}

func (r *RouterGroup) UseResponse(m ...middleware.Middleware) *RouterGroup {
	r.responseMiddleware = append(r.responseMiddleware, m...)
	return r
}

func (r *RouterGroup) Handle(method string, path string, handler Handler) {
	r.group.Handle(method, path, func(ctx *fasthttp.RequestCtx) {
		for _, m := range r.requestMiddlewares {
			ctx = m(ctx)
		}
		luxCtx := &LuxContext{ctx: ctx}
		if luxCtx.Ok() {
			handler(luxCtx)
			luxCtx.Ok()
			for _, m := range r.responseMiddleware {
				ctx = m(luxCtx.ctx)
			}
		}
	})
}

func (r *RouterGroup) Get(path string, handler Handler) {
	r.Handle(GET, path, handler)
}

func (r *RouterGroup) Post(path string, handler Handler) {
	r.Handle(POST, path, handler)
}

func (r *RouterGroup) Head(path string, handler Handler) {
	r.Handle(HEAD, path, handler)
}

func (r *RouterGroup) Delete(path string, handler Handler) {
	r.Handle(DELETE, path, handler)
}

func (r *RouterGroup) Put(path string, handler Handler) {
	r.Handle(PUT, path, handler)
}

func (r *RouterGroup) Patch(path string, handler Handler) {
	r.Handle(PATCH, path, handler)
}

func (r *RouterGroup) Options(path string, handler Handler) {
	r.Handle(OPTIONS, path, handler)
}
