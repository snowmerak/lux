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

### simple get graphql

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
	root.GetGraph("", graphql.Fields{
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

### simple get template

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
	root.GetTemplateHTML("", "Hello, {{.Name}}!", struct{ Name string }{Name: "World"})
	if err := app.ListenAndServe(":8080"); err != nil {
		panic(err)
	}
}
```

### simple post protobuf

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
	"github.com/snowmerak/lux/test/model/capsule"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func main() {
	app := lux.NewServer()
	root := app.RouterGroup("")
	root.Use(middleware.CORS(), middleware.CompressBrotli(), middleware.AllowStaticIPs("127.0.0.1", "localhost"))
	root.Preflight(lux.AllowAllOrigin, []string{"GET"}, lux.DefaultPreflightHeaders)
	root.Get("", func(lc *lux.LuxContext) {
		lc.ReplyString("Hello World")
	})
	root.PostProtobuf("", new(capsule.Capsule), func(pm protoreflect.ProtoMessage) (protoreflect.ProtoMessage, error) {
		capsule := pm.(*capsule.Capsule)
		capsule.ID = "snowmerak"
		capsule.Data = []byte("Hello World")
		return capsule, nil
	})
	if err := app.ListenAndServe(":8080"); err != nil {
		panic(err)
	}
}
```