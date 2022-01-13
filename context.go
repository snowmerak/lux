package lux

import "github.com/valyala/fasthttp"

type LuxContext struct {
	ctx *fasthttp.RequestCtx
	Log Logger
}
