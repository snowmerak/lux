# Lux

## example

### simple

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.NewServer()
	root := app.RouterGroup("")
	root.UseResponse(middleware.CORS, middleware.CompressBrotli)
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