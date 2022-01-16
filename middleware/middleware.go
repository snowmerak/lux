package middleware

import (
	"compress/gzip"

	"strings"

	"github.com/golang/snappy"
	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/loglevel"
	"github.com/snowmerak/lux/logger"
	"github.com/valyala/fasthttp"
)

const MIDDLEWARE = "middleware"

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

func compressSnappy(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	defer func() {
		if err, ok := recover().(error); ok && err != nil {
			logger.Write(MIDDLEWARE, log.New(loglevel.Error, string(ctx.Path())+": "+err.Error()).End())
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		}
	}()
	if string(ctx.Request.Header.Peek("Content-Encoding")) != "snappy" {
		body := snappy.Encode(nil, ctx.Request.Body())
		ctx.Request.SetBody(body)
		ctx.Request.Header.Set("Content-Encoding", "snappy")
	}
	return ctx
}

func CompressSnappy() MiddlewareSet {
	return MiddlewareSet{
		nil,
		compressSnappy,
	}
}

func compressGzip(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	if string(ctx.Request.Header.Peek("Content-Encoding")) != "gzip" {
		body := ctx.Response.Body()
		ctx.Response.ResetBody()
		writer := gzip.NewWriter(ctx.Response.BodyWriter())
		defer writer.Close()
		if _, err := writer.Write(body); err != nil {
			logger.Write(MIDDLEWARE, log.New(loglevel.Error, string(ctx.Path())+": "+err.Error()).End())
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			return ctx
		}
		ctx.Request.Header.Set("Content-Encoding", "gzip")
	}
	return ctx
}

func CompressGzip() MiddlewareSet {
	return MiddlewareSet{
		nil,
		compressGzip,
	}
}

func compressBrotli(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	if string(ctx.Request.Header.Peek("Content-Encoding")) != "br" {
		body := ctx.Response.Body()
		ctx.Response.ResetBody()
		writer := ctx.Response.BodyWriter()
		if _, err := writer.Write(body); err != nil {
			logger.Write(MIDDLEWARE, log.New(loglevel.Error, string(ctx.Path())+": "+err.Error()).End())
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			return ctx
		}
		ctx.Request.Header.Set("Content-Encoding", "br")
	}
	return ctx
}

func CompressBrotli() MiddlewareSet {
	return MiddlewareSet{
		nil,
		compressBrotli,
	}
}

func Authenticate(tokenChecker func(authorization []byte, tokenCookie []byte) error) MiddlewareSet {
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			token, cookie := ctx.Request.Header.Peek("Authorization"), ctx.Request.Header.Cookie("token")
			if token == nil {
				ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
				return ctx
			}
			if err := tokenChecker(token, cookie); err != nil {
				ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
				return ctx
			}
			return ctx
		},
		nil,
	}
}

func SetStaticAllowListIP(ips ...string) MiddlewareSet {
	cache := make(map[string]struct{})
	for _, ip := range ips {
		cache[ip] = struct{}{}
	}
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			if _, ok := cache[ctx.RemoteIP().String()]; ok {
				ctx.Response.SetStatusCode(fasthttp.StatusOK)
				return ctx
			}
			ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
			return ctx
		},
		nil,
	}
}

func SetStaticBlockListIP(ips ...string) MiddlewareSet {
	cache := make(map[string]struct{})
	for _, ip := range ips {
		cache[ip] = struct{}{}
	}
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			if _, ok := cache[ctx.RemoteIP().String()]; ok {
				ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
				return ctx
			}
			ctx.Response.SetStatusCode(fasthttp.StatusOK)
			return ctx
		},
		nil,
	}
}

func SetDynamicAllowListIP(checker func(ip string) bool) MiddlewareSet {
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			if checker(ctx.RemoteIP().String()) {
				ctx.Response.SetStatusCode(fasthttp.StatusOK)
				return ctx
			}
			ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
			return ctx
		},
		nil,
	}
}

func SetDynamicBlockListIP(checker func(ip string) bool) MiddlewareSet {
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			if checker(ctx.RemoteIP().String()) {
				ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
				return ctx
			}
			ctx.Response.SetStatusCode(fasthttp.StatusOK)
			return ctx
		},
		nil,
	}
}
