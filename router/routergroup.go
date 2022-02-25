package router

import (
	"strings"

	"github.com/diy-cloud/lux/context"
	"github.com/diy-cloud/lux/handler"
	"github.com/diy-cloud/lux/logext"
	"github.com/diy-cloud/lux/middleware"
	"github.com/diy-cloud/lux/swagger"
)

type RouterGroup struct {
	Path            string
	Middlewares     []middleware.Set
	Routers         map[string]map[string]*Router
	SubRouterGroups []*RouterGroup
	Logger          *logext.Logger
	Swagger         *swagger.Swagger
}

func (r *RouterGroup) UseMiddlewares(middlewares ...middleware.Set) {
	r.Middlewares = append(r.Middlewares, middlewares...)
}

func (r *RouterGroup) AddRouter(method, path string, handler handler.Handler, swaggerRouter *swagger.Router, middlewares ...middleware.Set) *Router {
	if swaggerRouter != nil {
		p := r.Path + path
		m := strings.ToLower(method)
		if r.Swagger.Paths == nil {
			r.Swagger.Paths = map[swagger.Path]map[swagger.Method]swagger.Router{}
		}
		if _, ok := r.Swagger.Paths[swagger.Path(p)]; !ok {
			r.Swagger.Paths[swagger.Path(p)] = map[swagger.Method]swagger.Router{}
		}
		r.Swagger.Paths[swagger.Path(p)][swagger.Method(m)] = *swaggerRouter
	}
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
