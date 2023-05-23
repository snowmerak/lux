package middleware

import (
	"github.com/snowmerak/lux/context"
	"net/http"
)

type AuthChecker func(lc *context.LuxContext, authorizationHeader string, tokenCookies ...*http.Cookie) bool

func Auth(authChecker AuthChecker, tokenName ...string) Set {
	return Set{
		Request: func(ctx *context.LuxContext) (*http.Request, int) {
			req := ctx.Request
			authorizationHeader := req.Header.Get("Authorization")
			cookies := []*http.Cookie(nil)
			for _, name := range tokenName {
				cookie, err := req.Cookie(name)
				if err == nil {
					cookies = append(cookies, cookie)
				}
			}
			if authChecker(ctx, authorizationHeader, cookies...) {
				return req, http.StatusOK
			}
			return req, http.StatusUnauthorized
		},
		Response: nil,
	}
}
