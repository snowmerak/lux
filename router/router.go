package router

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/handler"
	"github.com/snowmerak/lux/logext"
	"github.com/snowmerak/lux/middleware"
	"github.com/snowmerak/lux/swagger"
	"github.com/snowmerak/lux/util"
)

type Router struct {
	Handler     handler.Handler
	Middlewares []middleware.Set
	Method      string

	logger *logext.Logger
}

func (r *Router) UseMiddlewares(middlewares ...middleware.Set) {
	r.Middlewares = append(r.Middlewares, middlewares...)
}

func (r *RouterGroup) GET(path string, handler handler.Handler, swaggerRouter *swagger.Router, middlewares ...middleware.Set) *Router {
	return r.AddRouter("GET", path, handler, swaggerRouter, middlewares...)
}

func (r *RouterGroup) POST(path string, handler handler.Handler, swaggerRouter *swagger.Router, middlewares ...middleware.Set) *Router {
	return r.AddRouter("POST", path, handler, swaggerRouter, middlewares...)
}

func (r *RouterGroup) PATCH(path string, handler handler.Handler, swaggerRouter *swagger.Router, middlewares ...middleware.Set) *Router {
	return r.AddRouter("PATCH", path, handler, swaggerRouter, middlewares...)
}

func (r *RouterGroup) PUT(path string, handler handler.Handler, swaggerRouter *swagger.Router, middlewares ...middleware.Set) *Router {
	return r.AddRouter("PUT", path, handler, swaggerRouter, middlewares...)
}

func (r *RouterGroup) DELETE(path string, handler handler.Handler, swaggerRouter *swagger.Router, middlewares ...middleware.Set) *Router {
	return r.AddRouter("DELETE", path, handler, swaggerRouter, middlewares...)
}

func (r *RouterGroup) OPTIONS(path string, handler handler.Handler, swaggerRouter *swagger.Router, middlewares ...middleware.Set) *Router {
	return r.AddRouter("OPTIONS", path, handler, swaggerRouter, middlewares...)
}

func (r *RouterGroup) HEAD(path string, handler handler.Handler, swaggerRouter *swagger.Router, middlewares ...middleware.Set) *Router {
	return r.AddRouter("HEAD", path, handler, swaggerRouter, middlewares...)
}

func (r *RouterGroup) Preflight(allowOrigins, allowMethods, allowHeaders []string, swaggerRouter *swagger.Router, middlewares ...middleware.Set) *Router {
	return r.OPTIONS("/", func(l *context.LuxContext) error {
		origin := l.Request.Header.Get("Origin")
		for _, o := range allowOrigins {
			if o == "*" || o == origin {
				l.Response.Header().Set("Access-Control-Allow-Origin", o)
				break
			}
		}
		method := l.Request.Header.Get("Access-Control-Request-Method")
		for _, m := range allowMethods {
			if m == "*" || m == method {
				l.Response.Header().Set("Access-Control-Allow-Methods", m)
				break
			}
		}
		headers := strings.Split(l.Request.Header.Get("Access-Control-Request-Headers"), ",")
		for _, h := range headers {
			for _, a := range allowHeaders {
				if a == "*" || a == h {
					l.Response.Header().Set("Access-Control-Allow-Headers", a)
					break
				}
			}
		}
		l.Response.Headers.Set("Access-Control-Allow-Credentials", "true")
		l.Response.Headers.Set("Access-Control-Max-Age", "86400")
		l.Response.StatusCode = 200
		return nil
	}, swaggerRouter, middlewares...)
}

func (r *RouterGroup) Statics(path string, folderPath string, middlewares ...middleware.Set) *Router {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if !strings.HasSuffix(path, "*filepath") {
		path += "*filepath"
	}
	return r.AddRouter("GET", path, func(l *context.LuxContext) error {
		path := l.GetPathVariable("filepath")
		paths := append([]string{folderPath}, strings.Split(path, "/")...)
		f, err := os.Open(filepath.Join(paths...))
		if err != nil {
			l.SetBadRequest()
			return err
		}
		defer f.Close()
		if os.IsNotExist(err) {
			l.SetBadRequest()
			return err
		}
		data, err := io.ReadAll(f)
		if err != nil {
			l.SetInternalServerError()
			return err
		}
		l.Response.Headers.Set("Content-Type", util.GetContentTypeFromExt(filepath.Ext(path)))
		l.Response.Body = data
		l.SetOK()
		return nil
	}, nil)
}

func (r *RouterGroup) Embedded(path string, embed fs.FS, middlewares ...middleware.Set) *Router {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if !strings.HasSuffix(path, "*filepath") {
		path += "*filepath"
	}
	return r.AddRouter("GET", path, func(l *context.LuxContext) error {
		path := l.GetPathVariable("filepath")
		path = strings.TrimPrefix(path, "/")
		file, err := embed.Open(path)
		if err != nil {
			l.SetBadRequest()
			return err
		}
		data, err := io.ReadAll(file)
		if err != nil {
			l.SetInternalServerError()
			return err
		}
		l.Response.Headers.Set("Content-Type", util.GetContentTypeFromExt(filepath.Ext(path)))
		l.Response.Body = data
		l.SetOK()
		return nil
	}, nil)
}
