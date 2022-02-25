package lux

import (
	"strings"

	"github.com/diy-cloud/lux/middleware"
	"github.com/diy-cloud/lux/router"
)

func (l *Lux) NewRouterGroup(path string, middlewares ...middleware.Set) *router.RouterGroup {
	rg := &router.RouterGroup{
		Path:            path,
		Middlewares:     middlewares,
		Routers:         map[string]map[string]*router.Router{},
		SubRouterGroups: []*router.RouterGroup{},
		Logger:          l.logger,
		Swagger:         l.swagger,
	}
	if strings.HasSuffix(rg.Path, "/") {
		rg.Path = strings.TrimSuffix(rg.Path, "/")
	}
	l.routers = append(l.routers, rg)
	return rg
}
