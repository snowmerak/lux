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

func CompressGzip(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
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

func CompressBrotli(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
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
