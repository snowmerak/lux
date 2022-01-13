package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"

	"github.com/golang/snappy"
	"github.com/valyala/fasthttp"
)

//LimitBody: Limit the size of the body. If the body is bigger than the limit, set the status code to 413 and return an error.
//unit: bytes
func LimitBody(size int) Middleware {
	return func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
		if len(ctx.Request.Body()) > size {
			ctx.Response.SetStatusCode(fasthttp.StatusRequestEntityTooLarge)
		}
		return ctx
	}
}

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
			log.Println(string(ctx.Path()) + ": " + err.Error())
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
					log.Println(string(ctx.Path()) + ": " + err.Error())
				}
				break
			}
			ctx.Request.AppendBody(buf[:p])
		}
		if err != nil {
			log.Println(string(ctx.Path()) + ": " + err.Error())
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			return ctx
		}
		ctx.Request.Header.DelBytes([]byte("Content-Encoding"))
	}
	return ctx
}

func DecompressBrotli(ctx *fasthttp.RequestCtx)
