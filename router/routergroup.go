package router

import (
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/handler"
	"github.com/snowmerak/lux/logext"
	"github.com/snowmerak/lux/middleware"
)

type RouterGroup struct {
	Path            string
	Middlewares     []middleware.Set
	Routers         map[string]map[string]*Router
	SubRouterGroups []*RouterGroup
	Logger          *logext.Logger
}

func (r *RouterGroup) UseMiddlewares(middlewares ...middleware.Set) {
	r.Middlewares = append(r.Middlewares, middlewares...)
}

func (r *RouterGroup) AddRouter(method, path string, handler handler.Handler, middlewares ...middleware.Set) *Router {
	router := &Router{
		Handler:     handler,
		Middlewares: middlewares,
		logger:      r.Logger,
	}
	handler = func(ctx *context.LuxContext) error {
		if rs := middleware.ApplyRequests(ctx, r.Middlewares); rs != "" {
			r.Logger.Warnf("Router %s %s: %s from %s", ctx.Request.Method, ctx.Request.URL.Path, rs, ctx.Request.RemoteAddr)
			return nil
		}
		if rs := middleware.ApplyRequests(ctx, router.Middlewares); rs != "" {
			r.Logger.Warnf("Router %s %s: %s from %s", ctx.Request.Method, ctx.Request.URL.Path, rs, ctx.Request.RemoteAddr)
			return nil
		}
		if err := handler(ctx); err != nil {
			return err
		}
		if rs := middleware.ApplyResponses(ctx, router.Middlewares); rs != "" {
			r.Logger.Warnf("Router %s %s: %s from %s", ctx.Request.Method, ctx.Request.URL.Path, rs, ctx.Request.RemoteAddr)
			return nil
		}
		if rs := middleware.ApplyResponses(ctx, r.Middlewares); rs != "" {
			r.Logger.Warnf("Router %s %s: %s from %s", ctx.Request.Method, ctx.Request.URL.Path, rs, ctx.Request.RemoteAddr)
			return nil
		}
		return nil
	}
	if _, ok := r.Routers[r.Path+path]; !ok {
		r.Routers[r.Path+path] = map[string]*Router{}
	}
	r.Routers[r.Path+path][method] = router
	return router
}
