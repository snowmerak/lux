package middleware

import (
	"strings"

	"github.com/valyala/fasthttp"
)

/*
CORS ...
CORS middleware
allow all origins
allow all methods
allow all headers
allow credentials
*/
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

/*
AllowHeaders ...
allow given headers
*/
func AllowHeaders(header ...string) MiddlewareSet {
	return MiddlewareSet{
		nil,
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			ctx.Response.Header.Set("Access-Control-Allow-Headers", strings.Join(header, ","))
			return ctx
		},
	}
}

/*
AllowMethods ...
allow given methods
*/
func AllowMethods(method ...string) MiddlewareSet {
	return MiddlewareSet{
		nil,
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			ctx.Response.Header.Set("Access-Control-Allow-Methods", strings.Join(method, ","))
			return ctx
		},
	}
}

/*
AllowOrigins ...
allow given origins
*/
func AllowOrigins(origin ...string) MiddlewareSet {
	return MiddlewareSet{
		nil,
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			ctx.Response.Header.Set("Access-Control-Allow-Origin", strings.Join(origin, ","))
			return ctx
		},
	}
}

/*
AllowCredentials ...
allow credentials
*/
func AllowCredentials() MiddlewareSet {
	return MiddlewareSet{
		nil,
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			return ctx
		},
	}
}
