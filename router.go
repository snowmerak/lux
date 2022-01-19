package lux

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/fasthttp/router"
	"github.com/graphql-go/graphql"
	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/loglevel"
	"github.com/snowmerak/lux/logger"
	"github.com/snowmerak/lux/middleware"
	"github.com/snowmerak/lux/swagger"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	GET     = "GET"
	POST    = "POST"
	HEAD    = "HEAD"
	DELETE  = "DELETE"
	PUT     = "PUT"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

var AllowAllOrigin = []string{"*"}
var DefaultPreflightHeaders = []string{"Origin", "Accept", "Content-Type"}

/*
RouterGroup ...
RouterGroup is a wrapper of router.Group and has some middlewares.
*/
type RouterGroup struct {
	group               *router.Group
	requestMiddlewares  []middleware.Middleware
	responseMiddlewares []middleware.Middleware

	swagger *swagger.Swagger
	path    string
}

/*
Use ...
append middlewares from given middleware.MiddlewareSet to RouterGroup.
*/
func (r *RouterGroup) Use(middlewareset ...middleware.MiddlewareSet) *RouterGroup {
	for _, m := range middlewareset {
		req, res := m.Request, m.Response
		if req != nil {
			r.requestMiddlewares = append(r.requestMiddlewares, req)
		}
		if res != nil {
			r.responseMiddlewares = append(r.responseMiddlewares, res)
		}
	}
	return r
}

/*
Handle ...
register router with given method and path to RouterGroup, and apply handler to it.
*/
func (r *RouterGroup) Handle(method string, path string, handler Handler, swaggerInfo *swagger.Router) {
	if swaggerInfo == nil {
		swaggerInfo = &swagger.Router{}
	}
	name := swagger.Path(r.path)
	if path != "" {
		name = swagger.Path(r.path + "/" + path)
	}
	if r.swagger.Paths[name] == nil {
		r.swagger.Paths[name] = make(map[swagger.Method]swagger.Router)
	}
	r.swagger.Paths[name][swagger.Method(strings.ToLower(method))] = *swaggerInfo
	r.group.Handle(method, path, func(ctx *fasthttp.RequestCtx) {
		luxCtx := &LuxContext{ctx: ctx}
		for _, m := range r.requestMiddlewares {
			if !luxCtx.Ok() {
				return
			}
			luxCtx.ctx = m(luxCtx.ctx)
		}
		if luxCtx.Ok() {
			handler(luxCtx)
		}
		for _, m := range r.responseMiddlewares {
			if !luxCtx.Ok() {
				return
			}
			ctx = m(luxCtx.ctx)
		}
	})
}

/*
Get ...
register router with GET method and given path to RouterGroup, and apply handler to it.
*/
func (r *RouterGroup) Get(path string, handler Handler, swaggerInfo *swagger.Router) {
	r.Handle(GET, path, handler, swaggerInfo)
}

/*
Post ...
register router with POST method and given path to RouterGroup, and apply handler to it.
*/
func (r *RouterGroup) Post(path string, handler Handler, swaggerInfo *swagger.Router) {
	r.Handle(POST, path, handler, swaggerInfo)
}

/*
Head ...
register router with HEAD method and given path to RouterGroup, and apply handler to it.
*/
func (r *RouterGroup) Head(path string, handler Handler, swaggerInfo *swagger.Router) {
	r.Handle(HEAD, path, handler, swaggerInfo)
}

/*
Delete ...
register router with DELETE method and given path to RouterGroup, and apply handler to it.
*/
func (r *RouterGroup) Delete(path string, handler Handler, swaggerInfo *swagger.Router) {
	r.Handle(DELETE, path, handler, swaggerInfo)
}

/*
Put ...
register router with PUT method and given path to RouterGroup, and apply handler to it.
*/
func (r *RouterGroup) Put(path string, handler Handler, swaggerInfo *swagger.Router) {
	r.Handle(PUT, path, handler, swaggerInfo)
}

/*
Patch ...
register router with PATCH method and given path to RouterGroup, and apply handler to it.
*/
func (r *RouterGroup) Patch(path string, handler Handler, swaggerInfo *swagger.Router) {
	r.Handle(PATCH, path, handler, swaggerInfo)
}

/*
Options ...
register router with OPTIONS method and given path to RouterGroup, and apply handler to it.
*/
func (r *RouterGroup) Options(path string, handler Handler, swaggerInfo *swagger.Router) {
	r.Handle(OPTIONS, path, handler, swaggerInfo)
}

/*
Preflight ...
register router with OPTIONS method to default path "/" to RouterGroup,
and apply CORS preflight to it.
*/
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

/*
Statics ...
register router with GET method to given path to RouterGroup, and apply static file server to it.
*/
func (r *RouterGroup) Statics(path string, root string) {
	if path == "" {
		path = "/"
	}
	if !strings.HasSuffix(path, "{filepath:*}") {
		path += "{filepath:*}"
	}
	r.group.ServeFiles(path, root)
}

/*
Embedded ...
register router with GET method to given path to RouterGroup, and apply embedded file server to it.
*/
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
			lc.Log.WriteLog(logger.SYSTEM, log.New(loglevel.Warn, "open file error: "+err.Error()).End())
			lc.ctx.Response.Header.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		stats, err := file.Stat()
		if err != nil {
			lc.Log.WriteLog(logger.SYSTEM, log.New(loglevel.Warn, "read file status: "+err.Error()).End())
			lc.ctx.Response.Header.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		buf := make([]byte, stats.Size())
		_, err = file.Read(buf)
		if err != nil {
			lc.Log.WriteLog(logger.SYSTEM, log.New(loglevel.Warn, "read file error: "+err.Error()).End())
			lc.ctx.Response.Header.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		lc.Reply(GetContentTypeFromExt(filepath.Ext(stats.Name())), buf)
	}, nil)
}

/*
GetContentTypeFromExt ...
get content type from file extension.
*/
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

/*
GetGraph ...
GraphQL with GET method.
*/
func (r *RouterGroup) GetGraph(path string, fields graphql.Fields) {
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemeConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	scheme, err := graphql.NewSchema(schemeConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create new graphql schema, %v", err))
	}

	r.Get(path, func(lc *LuxContext) {
		query := lc.GetURLParam("query")
		buf, err := base64.RawURLEncoding.DecodeString(query)
		if err == nil {
			query = string(buf)
		}
		params := graphql.Params{Schema: scheme, RequestString: query}
		result := graphql.Do(params)
		if len(result.Errors) > 0 {
			lc.Log.WriteLog(logger.SYSTEM, log.New(loglevel.Warn, "graphql error: "+result.Errors[0].Error()).End())
		}
		lc.ReplyJSON(result)
	}, nil)
}

/*
GetTemplateHTML ...
Template HTML with GET method.
*/
func (r *RouterGroup) GetTemplateHTML(path string, tmp string, data interface{}) {
	template, err := template.New("html").Parse(tmp)
	if err != nil {
		panic(fmt.Errorf("failed to create new template, %v", err))
	}
	typ := reflect.TypeOf(data)
	r.Get(path, func(lc *LuxContext) {
		defer func() {
			if err := recover(); err != nil {
				lc.Log.WriteLog(logger.SYSTEM, log.New(loglevel.Error, "template error: "+fmt.Sprint(err)).End())
			}
		}()
		val := reflect.New(typ).Elem()
		for i := 0; i < val.NumField(); i++ {
			val.Field(i).Set(reflect.ValueOf(lc.GetURLParam(val.Type().Field(i).Name)))
		}
		buf := bytes.NewBuffer(nil)
		template.Execute(buf, val.Interface())
		lc.ReplyHTML(buf.Bytes())
	}, nil)
}

/*
PostProtobuf ...
Protobuf function with POST method.
*/
func (r *RouterGroup) PostProtobuf(path string, typ protoreflect.ProtoMessage, handler func(protoreflect.ProtoMessage) (protoreflect.ProtoMessage, error)) {
	protoType := reflect.TypeOf(typ).Elem()
	r.Post(path, func(lc *LuxContext) {
		val := reflect.New(protoType).Elem()
		if err := proto.Unmarshal(lc.GetBody(), val.Addr().Interface().(proto.Message)); err != nil {
			lc.Log.WriteLog(logger.SYSTEM, log.New(loglevel.Error, path+" protobuf error: "+fmt.Sprint(err)).End())
			lc.BadRequest()
			return
		}
		result, err := handler(val.Addr().Interface().(protoreflect.ProtoMessage))
		if err != nil {
			lc.Log.WriteLog(logger.SYSTEM, log.New(loglevel.Error, path+" handler error: "+fmt.Sprint(err)).End())
			lc.BadRequest()
			return
		}
		lc.ReplyProtobuf(result)
	}, nil)
}
