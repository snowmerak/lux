package middleware

import (
	"fmt"
	"net/http"

	"github.com/snowmerak/lux/context"
)

type Set struct {
	Request  func(*context.LuxContext) (*context.LuxContext, int)
	Response func(*context.LuxContext) (*context.LuxContext, error)
}

func ApplyRequests(ctx *context.LuxContext, middlewares []Set) string {
	for _, m := range middlewares {
		if m.Request == nil {
			continue
		}
		_, code := m.Request(ctx)
		if 400 <= code && code < 600 {
			ctx.Response.WriteHeader(code)
			return fmt.Sprintf("Middleware Request Reading %s: %s from %s", ctx.Request.URL.Path, http.StatusText(code), ctx.Request.RemoteAddr)
		}
	}
	return ""
}

func ApplyResponses(ctx *context.LuxContext, middlewares []Set) string {
	for _, m := range middlewares {
		if m.Response == nil {
			continue
		}
		_, err := m.Response(ctx)
		if err != nil {
			return fmt.Sprintf("Middleware Response Writing %s: %s from %s", ctx.Request.URL.Path, err.Error(), ctx.Request.RemoteAddr)
		}
	}
	return ""
}
