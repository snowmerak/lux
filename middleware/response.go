package middleware

import (
	"strings"

	"github.com/golang/snappy"
	"github.com/valyala/fasthttp"
)

type Middleware func(*fasthttp.RequestCtx) *fasthttp.RequestCtx

func CORS(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "*")
	return ctx
}

func AllowHeaders(header ...string) Middleware {
	return func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
		ctx.Response.Header.Set("Access-Control-Allow-Headers", strings.Join(header, ","))
		return ctx
	}
}

func AllowMethods(method ...string) Middleware {
	return func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
		ctx.Response.Header.Set("Access-Control-Allow-Methods", strings.Join(method, ","))
		return ctx
	}
}

func AllowOrigins(origin ...string) Middleware {
	return func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", strings.Join(origin, ","))
		return ctx
	}
}

func AllowCredentials(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	return ctx
}

func CompressSnappy(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	if string(ctx.Request.Header.Peek("Content-Encoding")) != "snappy" {
		body := snappy.Encode(nil, ctx.Request.Body())
		ctx.Request.SetBody(body)
		ctx.Request.Header.Set("Content-Encoding", "snappy")
	}
	return ctx
}
