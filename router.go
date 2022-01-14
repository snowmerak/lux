package lux

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fasthttp/router"
	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/loglevel"
	"github.com/snowmerak/lux/logger"
	"github.com/snowmerak/lux/middleware"
	"github.com/valyala/fasthttp"
)

const GET = "GET"
const POST = "POST"
const HEAD = "HEAD"
const DELETE = "DELETE"
const PUT = "PUT"
const PATCH = "PATCH"
const OPTIONS = "OPTIONS"

var AllowAllOrigin = []string{"*"}
var DefaultPreflightHeaders = []string{"Origin", "Accept", "Content-Type"}

type RouterGroup struct {
	group              *router.Group
	requestMiddlewares []middleware.Middleware
	responseMiddleware []middleware.Middleware
}

func (l *Lux) RouterGroup(path ...string) *RouterGroup {
	group := l.router.Group("/" + strings.Join(path, "/"))
	return &RouterGroup{
		group:              group,
		requestMiddlewares: []middleware.Middleware{},
	}
}

func (r *RouterGroup) UseRequest(m ...middleware.Middleware) *RouterGroup {
	r.requestMiddlewares = append(r.requestMiddlewares, m...)
	return r
}

func (r *RouterGroup) UseResponse(m ...middleware.Middleware) *RouterGroup {
	r.responseMiddleware = append(r.responseMiddleware, m...)
	return r
}

type Router struct {
	requestMiddlewares  []middleware.Middleware
	responseMiddlewares []middleware.Middleware
}

func (r *Router) UseRequest(m ...middleware.Middleware) *Router {
	r.requestMiddlewares = append(r.requestMiddlewares, m...)
	return r
}

func (r *Router) UseResponse(m ...middleware.Middleware) *Router {
	r.responseMiddlewares = append(r.responseMiddlewares, m...)
	return r
}

func (r *RouterGroup) Handle(method string, path string, handler Handler) *Router {
	router := &Router{}
	r.group.Handle(method, path, func(ctx *fasthttp.RequestCtx) {
		luxCtx := &LuxContext{ctx: ctx}
		for _, m := range r.requestMiddlewares {
			if !luxCtx.Ok() {
				return
			}
			luxCtx.ctx = m(luxCtx.ctx)
		}
		for _, m := range router.requestMiddlewares {
			if !luxCtx.Ok() {
				return
			}
			luxCtx.ctx = m(luxCtx.ctx)
		}
		if luxCtx.Ok() {
			handler(luxCtx)
		}
		for _, m := range router.responseMiddlewares {
			if !luxCtx.Ok() {
				return
			}
			luxCtx.ctx = m(luxCtx.ctx)
		}
		for _, m := range r.responseMiddleware {
			if !luxCtx.Ok() {
				return
			}
			ctx = m(luxCtx.ctx)
		}
	})
	return router
}

func (r *RouterGroup) Get(path string, handler Handler) *Router {
	return r.Handle(GET, path, handler)
}

func (r *RouterGroup) Post(path string, handler Handler) *Router {
	return r.Handle(POST, path, handler)
}

func (r *RouterGroup) Head(path string, handler Handler) *Router {
	return r.Handle(HEAD, path, handler)
}

func (r *RouterGroup) Delete(path string, handler Handler) *Router {
	return r.Handle(DELETE, path, handler)
}

func (r *RouterGroup) Put(path string, handler Handler) *Router {
	return r.Handle(PUT, path, handler)
}

func (r *RouterGroup) Patch(path string, handler Handler) *Router {
	return r.Handle(PATCH, path, handler)
}

func (r *RouterGroup) Options(path string, handler Handler) *Router {
	return r.Handle(OPTIONS, path, handler)
}

func (r *RouterGroup) Preflight(allowOrigins []string, allowMethods []string, allowHeaders []string) {
	r.group.OPTIONS("", func(ctx *fasthttp.RequestCtx) {
		origin := string(ctx.Request.Header.Peek("Origin"))
		for _, o := range allowOrigins {
			if o == "*" {
				ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
				break
			}
			if o == origin {
				ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		method := string(ctx.Request.Header.Peek("Access-Control-Request-Method"))
		for _, m := range allowMethods {
			if m == method {
				ctx.Response.Header.Set("Access-Control-Allow-Methods", method)
				break
			}
		}
		headers := strings.Split(string(ctx.Request.Header.Peek("Access-Control-Request-Headers")), ",")
		for _, h := range headers {
			for _, a := range allowHeaders {
				if a == h {
					ctx.Response.Header.Add("Access-Control-Allow-Headers", h)
					break
				}
			}
		}
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Set("Access-Control-Max-Age", "1728000")
		ctx.Response.SetStatusCode(200)
	})
}

func (r *RouterGroup) Statics(path string, root string) {
	if path == "" {
		path = "/"
	}
	if !strings.HasSuffix(path, "{filepath:*}") {
		path += "{filepath:*}"
	}
	r.group.ServeFiles(path, root)
}

func (r *RouterGroup) Embedded(path string, embedded fs.FS) {
	if path == "" {
		path = "/"
	}
	if !strings.HasSuffix(path, "{filepath:*}") {
		path += "{filepath:*}"
	}
	r.Get(path, func(lc *LuxContext) {
		path, ok := lc.ctx.UserValue("filepath").(string)
		if !ok {
			lc.ctx.Response.Header.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		file, err := embedded.Open(path)
		if err != nil {
			lc.Log.Write(logger.SYSTEM, log.New(loglevel.Warn, "open file error: "+err.Error()).End())
			lc.ctx.Response.Header.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		stats, err := file.Stat()
		if err != nil {
			lc.Log.Write(logger.SYSTEM, log.New(loglevel.Warn, "read file status: "+err.Error()).End())
			lc.ctx.Response.Header.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		buf := make([]byte, stats.Size())
		_, err = file.Read(buf)
		if err != nil {
			lc.Log.Write(logger.SYSTEM, log.New(loglevel.Warn, "read file error: "+err.Error()).End())
			lc.ctx.Response.Header.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		lc.Reply(GetContentTypeFromExt(filepath.Ext(stats.Name())), buf)
	})
}

func GetContentTypeFromExt(ext string) string {
	contentType := "text/plain"
	switch ext {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".json":
		contentType = "application/json"
	case ".png":
		contentType = "image/png"
	case ".jpg":
		contentType = "image/jpeg"
	case ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	case ".svg":
		contentType = "image/svg+xml"
	case ".ico":
		contentType = "image/x-icon"
	case ".woff":
		contentType = "application/font-woff"
	case ".woff2":
		contentType = "application/font-woff2"
	case ".ttf":
		contentType = "application/font-ttf"
	case ".otf":
		contentType = "application/font-otf"
	case ".eot":
		contentType = "application/vnd.ms-fontobject"
	case ".mp4":
		contentType = "video/mp4"
	case ".webm":
		contentType = "video/webm"
	case ".ogv":
		contentType = "video/ogg"
	case ".mp3":
		contentType = "audio/mpeg"
	case ".wav":
		contentType = "audio/wav"
	case ".ogg":
		contentType = "audio/ogg"
	case ".flac":
		contentType = "audio/flac"
	case ".wma":
		contentType = "audio/x-ms-wma"
	case ".aac":
		contentType = "audio/aac"
	case ".m4a":
		contentType = "audio/m4a"
	case ".mpg":
		contentType = "video/mpeg"
	case ".mpeg":
		contentType = "video/mpeg"
	case ".avi":
		contentType = "video/x-msvideo"
	case ".mov":
		contentType = "video/quicktime"
	case ".zip":
		contentType = "application/zip"
	case ".rar":
		contentType = "application/x-rar-compressed"
	case ".7z":
		contentType = "application/x-7z-compressed"
	case ".tar":
		contentType = "application/x-tar"
	case ".gz":
		contentType = "application/gzip"
	case ".bz2":
		contentType = "application/x-bzip2"
	case ".doc":
		contentType = "application/msword"
	case ".docx":
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		contentType = "application/vnd.ms-excel"
	case ".xlsx":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		contentType = "application/vnd.ms-powerpoint"
	case ".pptx":
		contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".pdf":
		contentType = "application/pdf"
	case ".txt":
		contentType = "text/plain"
	case ".rtf":
		contentType = "application/rtf"
	case ".xml":
		contentType = "text/xml"
	case ".xsl":
		contentType = "text/xsl"
	case ".csv":
		contentType = "text/csv"
	case ".tsv":
		contentType = "text/tab-separated-values"
	}
	return contentType
}
