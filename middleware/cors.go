package middleware

import (
	"strings"

	"github.com/diy-cloud/lux/context"
)

func SetAllowHeaders(headers ...string) Set {
	return Set{
		Request: nil,
		Response: func(l *context.LuxContext) (*context.LuxContext, error) {
			l.Response.Headers.Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
			return l, nil
		},
	}
}

func SetAllowMethods(methods ...string) Set {
	return Set{
		Request: nil,
		Response: func(l *context.LuxContext) (*context.LuxContext, error) {
			l.Response.Headers.Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
			return l, nil
		},
	}
}

func SetAllowOrigins(origins ...string) Set {
	return Set{
		Request: nil,
		Response: func(l *context.LuxContext) (*context.LuxContext, error) {
			l.Response.Headers.Set("Access-Control-Allow-Origin", strings.Join(origins, ","))
			return l, nil
		},
	}
}

var SetAllowCredentials = Set{
	Request: nil,
	Response: func(l *context.LuxContext) (*context.LuxContext, error) {
		l.Response.Headers.Set("Access-Control-Allow-Credentials", "true")
		return l, nil
	},
}

var SetAllowCORS = Set{
	Request: nil,
	Response: func(l *context.LuxContext) (*context.LuxContext, error) {
		l.Response.Headers.Set("Access-Control-Allow-Headers", "*")
		l.Response.Headers.Set("Access-Control-Allow-Origin", "*")
		l.Response.Headers.Set("Access-Control-Allow-Methods", "*")
		l.Response.Headers.Set("Access-Control-Allow-Credentials", "true")
		return l, nil
	},
}
