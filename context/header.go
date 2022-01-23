package context

import "net/http"

func (l *LuxContext) SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	l.Response.Header().Add("Set-Cookie", cookie.String())
}

func (l *LuxContext) SetSecureCookie(name string, value string, maxAge int, path string, domain string) {
	l.SetCookie(name, value, maxAge, path, domain, true, true)
}
