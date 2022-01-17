package lux

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/loglevel"
	"github.com/snowmerak/lux/logger"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/*
LuxContext ...
LuxContext has fasthttp.RequestCtx for fasthttp context and logger.Logger for log.
*/
type LuxContext struct {
	ctx *fasthttp.RequestCtx
	Log logger.Logger
}

/*
Ok ...
Ok method returns true when l.ctx.Response.StatusCode is between 200 and 399.
over then 399, it returns false.
*/
func (l *LuxContext) Ok() bool {
	status := l.ctx.Response.StatusCode()
	return 200 <= status && status < 400
}

/*
Redirect ...
Redirect method redirects to the given url.
*/
func (l *LuxContext) Redirect(url string) {
	l.ctx.Redirect(url, fasthttp.StatusMovedPermanently)
}

/*
ReplyPlainText ...
ReplyPlainText method replies plain text to client with text/plain.
*/
func (l *LuxContext) ReplyPlainText(data string) {
	l.ctx.SetContentType("text/plain")
	l.ctx.SetStatusCode(fasthttp.StatusOK)
	l.ctx.Response.Header.Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
	l.ctx.Response.SetBodyString(data)
}

/*
ReplyString ...
ReplyString method is alias of ReplyPlainText.
*/
func (l *LuxContext) ReplyString(data string) {
	l.ReplyPlainText(data)
}

/*
ReplyJSON ...
ReplyJSON method replies json of given data struct to client with application/json.
*/
func (l *LuxContext) ReplyJSON(data interface{}) {
	l.ctx.SetContentType("application/json")
	l.ctx.SetStatusCode(fasthttp.StatusOK)
	encoder := json.NewEncoder(l.ctx.Response.BodyWriter())
	if err := encoder.Encode(data); err != nil {
		l.Log.WriteLog(logger.SYSTEM, log.New(loglevel.Error, err.Error()).End())
		l.ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

/*
Reply ...
Reply method replies given byte buffer with given content type.
*/
func (l *LuxContext) Reply(contentType string, data []byte) {
	l.ctx.SetContentType(contentType)
	l.ctx.SetStatusCode(fasthttp.StatusOK)
	l.ctx.Response.Header.Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
	l.ctx.Response.SetBody(data)
}

/*
ReplyProtobuf ...
ReplyProtobuf method replies given protobuf of data struct message with application/protobuf.
*/
func (l *LuxContext) ReplyProtobuf(data protoreflect.ProtoMessage) {
	buf, err := proto.Marshal(data)
	if err != nil {
		l.Log.WriteLog(logger.SYSTEM, log.New(loglevel.Error, err.Error()).End())
		l.ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
	l.Reply("application/protobuf", buf)
}

/*
ReplyBinary ...
ReplyBinary method replies given binary data with application/octet-stream.
*/
func (l *LuxContext) ReplyBinary(data []byte) {
	l.Reply("application/octet-stream", data)
}

/*
ReplyWebP ...
ReplyWebP method replies given webp data with image/webp.
*/
func (l *LuxContext) ReplyWebP(data []byte) {
	l.Reply("image/webp", data)
}

/*
ReplyWebM ...
ReplyWebM method replies given webm data with video/webm.
*/
func (l *LuxContext) ReplyWebM(data []byte) {
	l.Reply("video/webm", data)
}

/*
ReplyCSV ...
ReplyCSV method replies given csv data with text/csv.
*/
func (l *LuxContext) ReplyCSV(data []byte) {
	l.Reply("text/csv", data)
}

/*
ReplyHTML ...
ReplyHTML method replies given html data with text/html.
*/
func (l *LuxContext) ReplyHTML(data []byte) {
	l.Reply("text/html", data)
}

/*
ReplyXML ...
ReplyXML method replies given xml data with text/xml.
*/
func (l *LuxContext) ReplyXML(data []byte) {
	l.Reply("text/xml", data)
}

/*
ReplyExcel ...
ReplyExcel method replies given excel data with application/vnd.ms-excel.
*/
func (l *LuxContext) ReplyExcel(data []byte) {
	l.Reply("application/vnd.ms-excel", data)
}

/*
ReplyWord ...
ReplyWord method replies given word data with application/msword.
*/
func (l *LuxContext) ReplyWord(data []byte) {
	l.Reply("application/msword", data)
}

/*
ReplyPdf ...
ReplyPdf method replies given pdf data with application/pdf.
*/
func (l *LuxContext) ReplyPdf(data []byte) {
	l.Reply("application/pdf", data)
}

/*
ReplyJPEG ...
ReplyJPEG method replies given jpeg data with image/jpeg.
*/
func (l *LuxContext) ReplyJPEG(data []byte) {
	l.Reply("image/jpeg", data)
}

/*
ReplyPNG ...
ReplyPNG method replies given png data with image/png.
*/
func (l *LuxContext) ReplyPNG(data []byte) {
	l.Reply("image/png", data)
}

/*
ReplyFile ...
ReplyFile method replies given file data through fasthttp.Context.SendFile
*/
func (l *LuxContext) ReplyFile(path string) {
	l.ctx.SendFile(path)
}

/*
SetStatus ...
SetStatus method sets status code to l.ctx.Response.
*/
func (l *LuxContext) SetStatus(status int) {
	l.ctx.SetStatusCode(status)
}

/*
BadRequest ...
BadRequest method sets status code of context to 400.
*/
func (l *LuxContext) BadRequest() {
	l.SetStatus(fasthttp.StatusBadRequest)
}

/*
NotFound ...
NotFound method sets status code of context to 404.
*/
func (l *LuxContext) NotFound() {
	l.SetStatus(fasthttp.StatusNotFound)
}

/*
InternalServerError ...
InternalServerError method sets status code of context to 500.
*/
func (l *LuxContext) InternalServerError() {
	l.SetStatus(fasthttp.StatusInternalServerError)
}

/*
NotImplemented ...
NotImplemented method sets status code of context to 501.
*/
func (l *LuxContext) NotImplemented() {
	l.SetStatus(fasthttp.StatusNotImplemented)
}

/*
Unauthorized ...
Unauthorized method sets status code of context to 401.
*/
func (l *LuxContext) Unauthorized() {
	l.SetStatus(fasthttp.StatusUnauthorized)
}

/*
Forbidden ...
Forbidden method sets status code of context to 403.
*/
func (l *LuxContext) Forbidden() {
	l.SetStatus(fasthttp.StatusForbidden)
}

/*
GetFile ...
GetFile method returns multipart.FileHeader of given name from context.
*/
func (l *LuxContext) GetFile(name string) (*multipart.FileHeader, error) {
	file, err := l.ctx.FormFile(name)
	if err != nil {
		return nil, fmt.Errorf("lux.GetFile: %w", err)
	}
	return file, nil
}

/*
SaveFile ...
SaveFile method saves file of given name from context to given path.
*/
func (l *LuxContext) SaveFile(name, path string) error {
	f, err := l.ctx.FormFile(name)
	if err != nil {
		return fmt.Errorf("lux.SaveFile: %w", err)
	}
	fmt.Println(filepath.Join(path, f.Filename))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("lux.SaveFile: %w", err)
		}
	}
	if err := fasthttp.SaveMultipartFile(f, filepath.Join(path, f.Filename)); err != nil {
		return fmt.Errorf("lux.SaveFile: %w", err)
	}
	return nil
}

/*
GetFiles ...
GetFiles method returns multipart.FileHeaders of given names from context.
*/
func (l *LuxContext) GetFiles(name string) ([]*multipart.FileHeader, error) {
	parts, err := l.ctx.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("lux.GetFiles: %w", err)
	}
	return parts.File[name], nil
}

/*
SaveFiles ...
SaveFiles method saves files of given names from context to given path.
*/
func (l *LuxContext) SaveFiles(name string, path string) error {
	files, err := l.GetFiles(name)
	if err != nil {
		return fmt.Errorf("lux.SaveFiles: %w", err)
	}
	for _, file := range files {
		if err := fasthttp.SaveMultipartFile(file, filepath.Join(path, file.Filename)); err != nil {
			return fmt.Errorf("lux.SaveFiles: %w", err)
		}
	}
	return nil
}

/*
GetForm ...
GetForm method returns a value of given name from context.FromValue
*/
func (l *LuxContext) GetForm(name string) string {
	return string(l.ctx.FormValue(name))
}

/*
GetBody ...
GetBody method returns body of context.
*/
func (l *LuxContext) GetBody() []byte {
	return l.ctx.Request.Body()
}

/*
GetPostArgs ...
GetPostArgs method returns a value of given name from context's post arguments.
*/
func (l *LuxContext) GetPostArgs(name string) string {
	return string(l.ctx.PostArgs().Peek(name))
}

/*
GetParam ...
GetParam method returns a value of given name from context's URL arguments.
*/
func (l *LuxContext) GetParam(name string) string {
	return string(l.ctx.QueryArgs().Peek(name))
}

/*
GetRouteValue ...
GetRouteValue method returns a value of given name from httprouter's route values.
*/
func (l *LuxContext) GetRouteValue(name string) string {
	data, ok := l.ctx.UserValue(name).(string)
	if !ok {
		return ""
	}
	return data
}

/*
UpdateWebSocket ...
UpdateWebSocket method upgrades GET request to web socket of given handler.
*/
func (l *LuxContext) UpgradeWebSocket(upgrader websocket.FastHTTPUpgrader, handler func(*LuxContext, *websocket.Conn)) error {
	err := upgrader.Upgrade(l.ctx, func(c *websocket.Conn) {
		handler(l, c)
	})
	if err != nil {
		return fmt.Errorf("lux.UpgradeWebSocket: %w", err)
	}
	return nil
}

/*
SetCookie ...
SetCookie method sets cookie to context.
*/
func (l *LuxContext) SetCookie(key, value string, expireAt time.Time, httpOnly, secure bool) {
	ck := fasthttp.AcquireCookie()
	ck.SetKey(key)
	ck.SetValue(value)
	ck.SetHTTPOnly(httpOnly)
	ck.SetSecure(secure)
	ck.SetExpire(expireAt)
	l.ctx.Response.Header.SetCookie(ck)
}

/*
GetCookies ...
GetCookies method returns cookies of context.
*/
func (l *LuxContext) GetCookie(key string) string {
	return string(l.ctx.Request.Header.Cookie(key))
}

/*
GetIP ...
GetIP method returns remote IP of context.
*/
func (l *LuxContext) GetIP() string {
	return l.ctx.RemoteIP().String()
}

/*
GetPort ...
GetPort method returns remote port of context.
*/
func (l *LuxContext) GetPort() string {
	rip := l.ctx.RemoteAddr().String()
	if idx := strings.LastIndex(rip, ":"); idx > 0 {
		return rip[idx+1:]
	}
	return ""
}
