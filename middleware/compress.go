package middleware

import (
	"compress/gzip"

	"github.com/golang/snappy"
	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/loglevel"
	"github.com/snowmerak/lux/logger"
	"github.com/valyala/fasthttp"
)

func compressSnappy(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	defer func() {
		if err, ok := recover().(error); ok && err != nil {
			logger.Write(logger.MIDDLEWARE, log.New(loglevel.Error, string(ctx.Path())+": "+err.Error()).End())
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
			logger.Write(logger.MIDDLEWARE, log.New(loglevel.Error, string(ctx.Path())+": "+err.Error()).End())
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
			logger.Write(logger.MIDDLEWARE, log.New(loglevel.Error, string(ctx.Path())+": "+err.Error()).End())
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
