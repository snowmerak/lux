package middleware

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/golang/snappy"
	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/loglevel"
	"github.com/snowmerak/lux/logger"
	"github.com/valyala/fasthttp"
)

const MIDDLEWARE = "middleware"

func DecompressSnappy(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	if string(ctx.Request.Header.Peek("Content-Encoding")) == "snappy" {
		body, err := snappy.Decode(nil, ctx.Request.Body())
		if err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			return ctx
		}
		ctx.Request.Header.DelBytes([]byte("Content-Encoding"))
		ctx.Request.SetBody(body)
	}
	return ctx
}

func DecompressGzip(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	if string(ctx.Request.Header.Peek("Content-Encoding")) == "gzip" {
		reader, err := gzip.NewReader(bytes.NewReader(ctx.Request.Body()))
		if err != nil {
			logger.Write(MIDDLEWARE, log.New(loglevel.Error, string(ctx.Path())+": "+err.Error()).End())
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			return ctx
		}
		defer reader.Close()
		ctx.Request.ResetBody()
		buf := make([]byte, 1024)
		for {
			p, err := reader.Read(buf)
			if err != nil {
				if err != io.EOF {
					logger.Write(MIDDLEWARE, log.New(loglevel.Error, string(ctx.Path())+": "+err.Error()).End())
				}
				break
			}
			ctx.Request.AppendBody(buf[:p])
		}
		ctx.Request.Header.DelBytes([]byte("Content-Encoding"))
	}
	return ctx
}

func DecompressBrotli(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	if string(ctx.Request.Header.Peek("Content-Encoding")) == "br" {
		ctx.Request.Header.DelBytes([]byte("Content-Encoding"))
		reader := brotli.NewReader(bytes.NewReader(ctx.Request.Body()))
		ctx.Request.ResetBody()
		buf := make([]byte, 1024)
		for {
			p, err := reader.Read(buf)
			if err != nil {
				if err != io.EOF {
					logger.Write(MIDDLEWARE, log.New(loglevel.Error, string(ctx.Path())+": "+err.Error()).End())
				}
				break
			}
			ctx.Request.AppendBody(buf[:p])
		}
	}
	return ctx
}

func Authenticate(ctx *fasthttp.RequestCtx, tokenChecker func(authorization []byte, tokenCookie []byte) error) Middleware {
	return func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
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
	}
}
