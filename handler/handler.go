package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/logext"
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
