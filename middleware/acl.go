package middleware

import (
	"github.com/snowmerak/lux/context"
	"net/http"

	"github.com/snowmerak/lux/util"
)

func AllowStaticIPs(ips ...string) Set {
	return Set{
		Request: func(ctx *context.LuxContext) (*http.Request, int) {
			req := ctx.Request
			remoteIP := util.GetIP(req.RemoteAddr)
			for _, ip := range ips {
				if remoteIP == ip {
					return req, http.StatusOK
				}
			}
			return req, http.StatusForbidden
		},
		Response: nil,
	}
}

func BlockStaticIPs(ips ...string) Set {
	return Set{
		Request: func(ctx *context.LuxContext) (*http.Request, int) {
			req := ctx.Request
			remoteIP := util.GetIP(req.RemoteAddr)
			for _, ip := range ips {
				if remoteIP == ip {
					return req, http.StatusForbidden
				}
			}
			return req, http.StatusOK
		},
		Response: nil,
	}
}

func AllowDynamicIPs(checker func(remoteIP string) bool) Set {
	return Set{
		Request: func(ctx *context.LuxContext) (*http.Request, int) {
			req := ctx.Request
			remoteIP := util.GetIP(req.RemoteAddr)
			if checker(remoteIP) {
				return req, http.StatusOK
			}
			return req, http.StatusForbidden
		},
		Response: nil,
	}
}

func BlockDynamicIPs(checker func(remoteIP string) bool) Set {
	return Set{
		Request: func(ctx *context.LuxContext) (*http.Request, int) {
			req := ctx.Request
			remoteIP := util.GetIP(req.RemoteAddr)
			if checker(remoteIP) {
				return req, http.StatusForbidden
			}
			return req, http.StatusOK
		},
		Response: nil,
	}
}

func AllowStaticPorts(ports ...string) Set {
	return Set{
		Request: func(ctx *context.LuxContext) (*http.Request, int) {
			req := ctx.Request
			remotePort := util.GetPort(req.RemoteAddr)
			for _, port := range ports {
				if remotePort == port {
					return req, http.StatusOK
				}
			}
			return req, http.StatusForbidden
		},
		Response: nil,
	}
}

func BlockStaticPorts(ports ...string) Set {
	return Set{
		Request: func(ctx *context.LuxContext) (*http.Request, int) {
			req := ctx.Request
			remotePort := util.GetPort(req.RemoteAddr)
			for _, port := range ports {
				if remotePort == port {
					return req, http.StatusForbidden
				}
			}
			return req, http.StatusOK
		},
		Response: nil,
	}
}

func AllowDynamicPorts(checker func(remotePort string) bool) Set {
	return Set{
		Request: func(ctx *context.LuxContext) (*http.Request, int) {
			req := ctx.Request
			remotePort := util.GetPort(req.RemoteAddr)
			if checker(remotePort) {
				return req, http.StatusOK
			}
			return req, http.StatusForbidden
		},
		Response: nil,
	}
}

func BlockDynamicPorts(checker func(remotePort string) bool) Set {
	return Set{
		Request: func(ctx *context.LuxContext) (*http.Request, int) {
			req := ctx.Request
			remotePort := util.GetPort(req.RemoteAddr)
			if checker(remotePort) {
				return req, http.StatusForbidden
			}
			return req, http.StatusOK
		},
		Response: nil,
	}
}
