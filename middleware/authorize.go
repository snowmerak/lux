package middleware

import "github.com/valyala/fasthttp"

/*
Auth ...
authorize middleware
tokenChecker is consume authorization header and token cookie
*/
func Auth(tokenChecker func(authorizationHeader []byte, tokenCookie []byte) error) MiddlewareSet {
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
