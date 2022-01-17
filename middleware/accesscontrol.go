package middleware

import (
	"strings"

	"github.com/valyala/fasthttp"
)

/*
AllowStaticIPs ...
allow given IPs to access the server
*/
func AllowStaticIPs(ips ...string) MiddlewareSet {
	cache := make(map[string]struct{})
	for _, ip := range ips {
		cache[ip] = struct{}{}
	}
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			if _, ok := cache[ctx.RemoteIP().String()]; ok {
				ctx.Response.SetStatusCode(fasthttp.StatusOK)
				return ctx
			}
			ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
			return ctx
		},
		nil,
	}
}

/*
BlockStaticIPs ...
block given IPs to access the server
*/
func BlockStaticIPs(ips ...string) MiddlewareSet {
	cache := make(map[string]struct{})
	for _, ip := range ips {
		cache[ip] = struct{}{}
	}
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			if _, ok := cache[ctx.RemoteIP().String()]; ok {
				ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
				return ctx
			}
			ctx.Response.SetStatusCode(fasthttp.StatusOK)
			return ctx
		},
		nil,
	}
}

/*
AllowDynamicIPs ...
allow IP address passed checker to access the server
*/
func AllowDynamicIPs(checker func(ip string) bool) MiddlewareSet {
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			if checker(ctx.RemoteIP().String()) {
				ctx.Response.SetStatusCode(fasthttp.StatusOK)
				return ctx
			}
			ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
			return ctx
		},
		nil,
	}
}

/*
BlockDynamicIPs ...
block IP address passed checker to access the server
*/
func BlockDynamicIPs(checker func(ip string) bool) MiddlewareSet {
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			if checker(ctx.RemoteIP().String()) {
				ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
				return ctx
			}
			ctx.Response.SetStatusCode(fasthttp.StatusOK)
			return ctx
		},
		nil,
	}
}

/*
AllowStaticPorts ...
allow given ports to access the server
*/
func AllowStaticPorts(ports ...string) MiddlewareSet {
	cache := make(map[string]struct{})
	for _, port := range ports {
		cache[port] = struct{}{}
	}
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			rip := ctx.RemoteAddr().String()
			port := ""
			if idx := strings.LastIndex(rip, ":"); idx > 0 {
				port = rip[idx+1:]
			}
			if _, ok := cache[port]; ok {
				ctx.Response.SetStatusCode(fasthttp.StatusOK)
				return ctx
			}
			ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
			return ctx
		},
		nil,
	}
}

/*
BlockStaticPorts ...
block given ports to access the server
*/
func BlockStaticPorts(ports ...string) MiddlewareSet {
	cache := make(map[string]struct{})
	for _, port := range ports {
		cache[port] = struct{}{}
	}
	return MiddlewareSet{
		func(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
			rip := ctx.RemoteAddr().String()
			port := ""
			if idx := strings.LastIndex(rip, ":"); idx > 0 {
				port = rip[idx+1:]
			}
			if _, ok := cache[port]; ok {
				ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
				return ctx
			}
			ctx.Response.SetStatusCode(fasthttp.StatusOK)
			return ctx
		},
		nil,
	}
}
