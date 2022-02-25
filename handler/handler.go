package handler

import (
	"net/http"

	"github.com/diy-cloud/lux/context"
	"github.com/diy-cloud/lux/logext"
	"github.com/gobwas/ws"
	"github.com/julienschmidt/httprouter"
)

type Handler func(*context.LuxContext) error

func Wrap(logger *logext.Logger, handler Handler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		luxCtx := new(context.LuxContext)
		luxCtx.Request = r
		ok := false
		luxCtx.Response, ok = w.(*context.Response)
		if !ok {
			logger.Errorf("Router %s %s: Response is not a context.Response", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		luxCtx.RouteParams = ps
		if err := handler(luxCtx); err != nil {
			logger.Errorf("Router %s %s: %s from %s", r.Method, r.URL.Path, err, r.RemoteAddr)
		}
	}
}

type WSHandler func(*context.WSContext) error

func WSWrap(wsHandler WSHandler) Handler {
	return func(luxCtx *context.LuxContext) error {
		conn, _, _, err := ws.UpgradeHTTP(luxCtx.Request, luxCtx.Response)
		if err != nil {
			return err
		}
		defer conn.Close()
		wsCtx := new(context.WSContext)
		wsCtx.Conn = conn
		if err := wsHandler(wsCtx); err != nil {
			return err
		}
		return nil
	}
}
