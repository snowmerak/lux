package handler

import (
	ctx "context"
	"github.com/rs/zerolog"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/julienschmidt/httprouter"
	"github.com/snowmerak/lux/context"
)

type Handler func(*context.LuxContext) error

func Wrap(ctx ctx.Context, logger *zerolog.Logger, jwtCfg *context.JWTConfig, handler Handler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		luxCtx := new(context.LuxContext)
		luxCtx.Context = ctx
		luxCtx.Request = r
		luxCtx.Logger = logger
		luxCtx.JWTConfig = jwtCfg
		ok := false
		luxCtx.Response, ok = w.(*context.Response)
		if !ok {
			logger.Error().Str("method", r.Method).Str("path", r.URL.Path).Msg("Response is not a context.Response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		luxCtx.RouteParams = ps
		if err := handler(luxCtx); err != nil {
			logger.Error().Str("method", r.Method).Str("path", r.URL.Path).Str("remote", r.RemoteAddr).Err(err).Msg("Handler error")
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
