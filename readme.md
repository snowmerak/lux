# Lux

## example

### simple http

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.NewServer()
	root := app.RouterGroup("")
	root.Use(middleware.CORS(), middleware.CompressBrotli())
	root.Preflight(lux.AllowAllOrigin, []string{"GET"}, lux.DefaultPreflightHeaders)
	root.Get("{name?}", func(lc *lux.LuxContext) {
		greeting := "Hello!"
		name := lc.GetParam("name")
		if name != "" {
			greeting = "Hello, " + name + "!"
		}
		lc.ReplyString(greeting)
	})
	if err := app.ListenAndServe(":8080"); err != nil {
		panic(err)
	}
}
```

### simple https

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.NewServer()
	root := app.RouterGroup("")
	root.Use(middleware.CORS(), middleware.CompressBrotli())
	root.Preflight(lux.AllowAllOrigin, []string{"GET"}, lux.DefaultPreflightHeaders)
	root.Get("{name?}", func(lc *lux.LuxContext) {
		greeting := "Hello!"
		name := lc.GetParam("name")
		if name != "" {
			greeting = "Hello, " + name + "!"
		}
		lc.ReplyString(greeting)
	})
	if err := app.ListenAndServeTLS(":8080", "minica.pem", "minica-key.pem"); err != nil {
		panic(err)
	}
}
```

### simple graphql get

```go
package main

import (
	"github.com/graphql-go/graphql"
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.NewServer()
	root := app.RouterGroup("")
	root.Use(middleware.CORS(), middleware.CompressBrotli())
	root.Preflight(lux.AllowAllOrigin, []string{"GET"}, lux.DefaultPreflightHeaders)
	root.GraphGet("", graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	})
	if err := app.ListenAndServe(":8080"); err != nil {
		panic(err)
	}
}
```