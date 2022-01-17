package middleware

import (
	"strings"

	"github.com/valyala/fasthttp"
)

type Middleware func(*fasthttp.RequestCtx) *fasthttp.RequestCtx
type MiddlewareSet struct {
	Request  Middleware
	Response Middleware
}

func CORS() MiddlewareSet {
	return MiddlewareSet{
		nil,
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
			ctx.Response.Header.Set("Access-Control-Allow-Methods", "*")
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			return ctx
		},
	}
}

func AllowHeaders(header ...string) MiddlewareSet {
	return MiddlewareSet{
		nil,
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			ctx.Response.Header.Set("Access-Control-Allow-Headers", strings.Join(header, ","))
			return ctx
		},
	}
}

func AllowMethods(method ...string) MiddlewareSet {
	return MiddlewareSet{
		nil,
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			ctx.Response.Header.Set("Access-Control-Allow-Methods", strings.Join(method, ","))
			return ctx
		},
	}
}

func AllowOrigins(origin ...string) MiddlewareSet {
	return MiddlewareSet{
		nil,
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			ctx.Response.Header.Set("Access-Control-Allow-Origin", strings.Join(origin, ","))
			return ctx
		},
	}
}

func AllowCredentials() MiddlewareSet {
	return MiddlewareSet{
		nil,
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			return ctx
		},
	}
}
