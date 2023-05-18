package middleware

import "net/http"

func Auth(tokenChecker func(authorizationHeader string, tokenCookie *http.Cookie) bool) Set {
	return Set{
		Request: func(req *http.Request) (*http.Request, int) {
			authorizationHeader := req.Header.Get("Authorization")
			tokenCookie, _ := req.Cookie("token")
			if tokenChecker(authorizationHeader, tokenCookie) {
				return req, http.StatusOK
			}
			return req, http.StatusUnauthorized
		},
		Response: nil,
	}
}
