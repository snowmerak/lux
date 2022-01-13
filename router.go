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

type Router struct {
	requestMiddlewares  []middleware.Middleware
	responseMiddlewares []middleware.Middleware
}

func (r *Router) UseRequest(m ...middleware.Middleware) *Router {
	r.requestMiddlewares = append(r.requestMiddlewares, m...)
	return r
}

func (r *Router) UseResponse(m ...middleware.Middleware) *Router {
	r.responseMiddlewares = append(r.responseMiddlewares, m...)
	return r
}

func (r *RouterGroup) Handle(method string, path string, handler Handler) *Router {
	router := &Router{}
	r.group.Handle(method, path, func(ctx *fasthttp.RequestCtx) {
		luxCtx := &LuxContext{ctx: ctx}
		for _, m := range r.requestMiddlewares {
			if !luxCtx.Ok() {
				return
			}
			luxCtx.ctx = m(luxCtx.ctx)
		}
		for _, m := range router.requestMiddlewares {
			if !luxCtx.Ok() {
				return
			}
			luxCtx.ctx = m(luxCtx.ctx)
		}
		if luxCtx.Ok() {
			handler(luxCtx)
		}
		for _, m := range router.responseMiddlewares {
			if !luxCtx.Ok() {
				return
			}
			luxCtx.ctx = m(luxCtx.ctx)
		}
		for _, m := range r.responseMiddleware {
			if !luxCtx.Ok() {
				return
			}
			ctx = m(luxCtx.ctx)
		}
	})
	return router
}

func (r *RouterGroup) Get(path string, handler Handler) *Router {
	return r.Handle(GET, path, handler)
}

func (r *RouterGroup) Post(path string, handler Handler) *Router {
	return r.Handle(POST, path, handler)
}

func (r *RouterGroup) Head(path string, handler Handler) *Router {
	return r.Handle(HEAD, path, handler)
}

func (r *RouterGroup) Delete(path string, handler Handler) *Router {
	return r.Handle(DELETE, path, handler)
}

func (r *RouterGroup) Put(path string, handler Handler) *Router {
	return r.Handle(PUT, path, handler)
}

func (r *RouterGroup) Patch(path string, handler Handler) *Router {
	return r.Handle(PATCH, path, handler)
}

func (r *RouterGroup) Options(path string, handler Handler) *Router {
	return r.Handle(OPTIONS, path, handler)
}
