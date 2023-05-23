package router

import (
	"github.com/rs/zerolog"
	"strings"

	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/handler"
	"github.com/snowmerak/lux/middleware"
	"github.com/snowmerak/lux/swagger"
)

type RouterGroup struct {
	Path            string
	Middlewares     []middleware.Set
	Routers         map[string]map[string]*Router
	SubRouterGroups []*RouterGroup
	Logger          *zerolog.Logger
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
	handler = func(ctx *context.LuxContext) error {
		if rs := middleware.ApplyRequests(ctx, r.Middlewares); rs != "" {
			r.Logger.Error().Str("method", ctx.Request.Method).Str("path", ctx.Request.URL.Path).Str("remote", ctx.Request.RemoteAddr).Str("err", rs).Msg("Router group middleware error")
			return nil
		}
		if rs := middleware.ApplyRequests(ctx, middlewares); rs != "" {
			r.Logger.Error().Str("method", ctx.Request.Method).Str("path", ctx.Request.URL.Path).Str("remote", ctx.Request.RemoteAddr).Str("err", rs).Msg("Router middleware error")
			return nil
		}
		if err := handler(ctx); err != nil {
			return err
		}
		if rs := middleware.ApplyResponses(ctx, middlewares); rs != "" {
			r.Logger.Error().Str("method", ctx.Request.Method).Str("path", ctx.Request.URL.Path).Str("remote", ctx.Request.RemoteAddr).Str("err", rs).Msg("Router middleware error")
			return nil
		}
		if rs := middleware.ApplyResponses(ctx, r.Middlewares); rs != "" {
			r.Logger.Error().Str("method", ctx.Request.Method).Str("path", ctx.Request.URL.Path).Str("remote", ctx.Request.RemoteAddr).Str("err", rs).Msg("Router group middleware error")
			return nil
		}
		return nil
	}
	router := &Router{
		Handler:     handler,
		Middlewares: middlewares,
		logger:      r.Logger,
	}
	if _, ok := r.Routers[r.Path+path]; !ok {
		r.Routers[r.Path+path] = map[string]*Router{}
	}
	r.Routers[r.Path+path][method] = router

	return router
}
