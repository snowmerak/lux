package lux

import "github.com/valyala/fasthttp"

var globalMiddleware = []func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx{}

func PushGlobalMiddleware(m func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx) {
	globalMiddleware = append(globalMiddleware, m)
}
