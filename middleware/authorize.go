package middleware

import "net/http"

type AuthChecker func(authorizationHeader string, tokenCookies ...*http.Cookie) bool

func Auth(authChecker AuthChecker, tokenName ...string) Set {
	return Set{
		Request: func(req *http.Request) (*http.Request, int) {
			authorizationHeader := req.Header.Get("Authorization")
			cookies := []*http.Cookie(nil)
			for _, name := range tokenName {
				cookie, err := req.Cookie(name)
				if err == nil {
					cookies = append(cookies, cookie)
				}
			}
			if authChecker(authorizationHeader, cookies...) {
				return req, http.StatusOK
			}
			return req, http.StatusUnauthorized
		},
		Response: nil,
	}
}
