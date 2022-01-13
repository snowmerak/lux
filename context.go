package lux

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/snowmerak/logstream/log"
	"github.com/snowmerak/logstream/log/loglevel"
	"github.com/snowmerak/lux/logger"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
)

type LuxContext struct {
	ctx *fasthttp.RequestCtx
	Log Logger
}

func (l *LuxContext) Ok() bool {
	status := l.ctx.Response.StatusCode()
	return 200 <= status && status < 400
}

func (l *LuxContext) Redirect(url string) {
	l.ctx.Redirect(url, fasthttp.StatusMovedPermanently)
}

func (l *LuxContext) ReplyPlainText(data string) {
	l.ctx.SetContentType("text/plain")
	l.ctx.SetStatusCode(fasthttp.StatusOK)
	l.ctx.Response.Header.Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
	l.ctx.Response.SetBodyString(data)
}

func (l *LuxContext) ReplyString(data string) {
	l.ReplyPlainText(data)
}

func (l *LuxContext) ReplyJSON(data interface{}) {
	l.ctx.SetContentType("application/json")
	l.ctx.SetStatusCode(fasthttp.StatusOK)
	encoder := json.NewEncoder(l.ctx.Response.BodyWriter())
	if err := encoder.Encode(data); err != nil {
		l.Log.Write(logger.SYSTEM, log.New(loglevel.Error, err.Error()).End())
		l.ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

func (l *LuxContext) Reply(contentType string, data []byte) {
	l.ctx.SetContentType(contentType)
	l.ctx.SetStatusCode(fasthttp.StatusOK)
	l.ctx.Response.Header.Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
	l.ctx.Response.SetBody(data)
}

func (l *LuxContext) ReplyProtobuf(data proto.Message) {
	buf, err := proto.Marshal(data)
	if err != nil {
		l.Log.Write(logger.SYSTEM, log.New(loglevel.Error, err.Error()).End())
		l.ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
	l.Reply("application/protobuf", buf)
}

func (l *LuxContext) ReplyBinary(data []byte) {
	l.Reply("application/octet-stream", data)
}

func (l *LuxContext) ReplyWebP(data []byte) {
	l.Reply("image/webp", data)
}

func (l *LuxContext) ReplyWebM(data []byte) {
	l.Reply("video/webm", data)
}

func (l *LuxContext) ReplyCSV(data []byte) {
	l.Reply("text/csv", data)
}

func (l *LuxContext) ReplyHTML(data []byte) {
	l.Reply("text/html", data)
}

func (l *LuxContext) ReplyXML(data []byte) {
	l.Reply("text/xml", data)
}

func (l *LuxContext) ReplyExcel(data []byte) {
	l.Reply("application/vnd.ms-excel", data)
}

func (l *LuxContext) ReplyWord(data []byte) {
	l.Reply("application/msword", data)
}

func (l *LuxContext) ReplyPdf(data []byte) {
	l.Reply("application/pdf", data)
}

func (l *LuxContext) ReplyJPEG(data []byte) {
	l.Reply("image/jpeg", data)
}

func (l *LuxContext) ReplyPNG(data []byte) {
	l.Reply("image/png", data)
}

func (l *LuxContext) ReplyFile(path string) {
	l.ctx.SendFile(path)
}

func (l *LuxContext) GetFile(name, path string) error {
	f, err := l.ctx.FormFile(name)
	if err != nil {
		return fmt.Errorf("lux.GetFile: %w", err)
	}
	fmt.Println(filepath.Join(path, f.Filename))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("lux.GetFile: %w", err)
		}
	}
	if err := fasthttp.SaveMultipartFile(f, filepath.Join(path, f.Filename)); err != nil {
		return fmt.Errorf("lux.GetFile: %w", err)
	}
	return nil
}

func (l *LuxContext) GetForm(name string) string {
	return string(l.ctx.FormValue(name))
}

func (l *LuxContext) GetPostArgs(name string) string {
	return string(l.ctx.PostArgs().Peek(name))
}

func (l *LuxContext) GetParam(name string) string {
	data, ok := l.ctx.UserValue(name).(string)
	if !ok {
		return ""
	}
	return data
}
