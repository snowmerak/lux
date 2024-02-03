package context

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/julienschmidt/httprouter"
)

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

func NewResponse() *Response {
	return &Response{
		StatusCode: 200,
		Headers:    make(http.Header),
		Body:       []byte{},
	}
}

func (r *Response) WriteHeader(code int) {
	r.StatusCode = code
}

func (r *Response) Header() http.Header {
	return r.Headers
}

func (r *Response) Write(p []byte) (int, error) {
	r.Body = append(r.Body, p...)
	return len(p), nil
}

type LuxContext struct {
	Request        *http.Request
	Response       *Response
	RouteParams    httprouter.Params
	Context        context.Context
	RequestContext context.Context
	Logger         *zerolog.Logger
	JWTConfig      *JWTConfig
}

func (l *LuxContext) IsOk() bool {
	if 400 <= l.Response.StatusCode && l.Response.StatusCode < 600 {
		return false
	}
	return true
}
