package middleware

import "github.com/valyala/fasthttp"

type Middleware func(*fasthttp.RequestCtx) *fasthttp.RequestCtx

type MiddlewareSet struct {
	Request  Middleware
	Response Middleware
}

func New(request, response Middleware) MiddlewareSet {
	return MiddlewareSet{
		request,
		response,
	}
}
