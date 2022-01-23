# Lux

A web library collection based on net/http.

## hello, world!

### HTTP

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootRouterGroup := app.NewRouterGroup("/")
	rootRouterGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("Hello World!")
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### HTTP TLS

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootRouterGroup := app.NewRouterGroup("/")
	rootRouterGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("Hello World!")
	})

	if err := app.ListenAndServe1TLS(":8080", "cert.pem", "key.pem"); err != nil {
		panic(err)
	}
}
```

### HTTP2

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootRouterGroup := app.NewRouterGroup("/")
	rootRouterGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("Hello World!")
	})

	if err := app.ListenAndServe2(":8080"); err != nil {
		panic(err)
	}
}
```

### HTTP2 TLS

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootRouterGroup := app.NewRouterGroup("/")
	rootRouterGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("Hello World!")
	})

	if err := app.ListenAndServe2TLS(":8080", "cert.pem", "key.pem"); err != nil {
		panic(err)
	}
}
```

## Server

### set logger

```go
package main

import (
	"os"

	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/logext"
)

func main() {
	app := lux.New()

	app.SetLogger(logext.New(os.Stderr))

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

By calling `logext.New(writer io.Writer)`, lux server will write log to `writer`.

### set max header bytes

```go
package main

import (
	"github.com/snowmerak/lux"
)

func main() {
	app := lux.New()

	app.SetMaxHeaderBytes(1024)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### set idle timeout

```go
package main

import (
	"time"

	"github.com/snowmerak/lux"
)

func main() {
	app := lux.New()

	app.SetIdleTimeout(time.Second * 10)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

If `KeepAlive` is enabled, wait connection during IdleTimeout.

### set read and write timeout

```go
package main

import (
	"time"

	"github.com/snowmerak/lux"
)

func main() {
	app := lux.New()

	app.SetReadTimeout(time.Second * 5)
	app.SetWriteTimeout(time.Second * 5)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

## reply

### string

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("reply string")
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### binary

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyBinary([]byte("Hello World!"))
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### file

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyFile("./public/sample.txt")
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### JSON

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyJSON(map[string]string{"hello": "world"})
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### ETC

```go
func (l *LuxContext) Reply(contentType string, body []byte) error

func (l *LuxContext) ReplyProtobuf(data protoreflect.ProtoMessage) error

func (l *LuxContext) ReplyWebP(data []byte) error

func (l *LuxContext) ReplyJPEG(data []byte) error

func (l *LuxContext) ReplyPNG(data []byte) error

func (l *LuxContext) ReplyWebM(data []byte) error

func (l *LuxContext) ReplyMP4(data []byte) error

func (l *LuxContext) ReplyCSV(data []byte) error

func (l *LuxContext) ReplyHTML(data []byte) error

func (l *LuxContext) ReplyXML(data []byte) error

func (l *LuxContext) ReplyExcel(data []byte) error

func (l *LuxContext) ReplyWord(data []byte) error

func (l *LuxContext) ReplyPDF(data []byte) error

func (l *LuxContext) ReplyPowerpoint(data []byte) error

func (l *LuxContext) ReplyZip(data []byte) error 

func (l *LuxContext) ReplyTar(data []byte) error

func (l *LuxContext) ReplyGZIP(data []byte) error

func (l *LuxContext) Reply7Z(data []byte) error
```

## set status

### set status

```go
package main

import (
	"net/http"

	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.POST("/", func(lc *context.LuxContext) error {
		lc.SetStatus(http.StatusNoContent)
		return nil
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### OK

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		lc.SetOK()
		return nil
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### bad request

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		lc.SetBadRequest()
		return nil
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### not found

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		lc.SetNotFound()
		return nil
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### ETC

```go
func (l *LuxContext) SetAccepted()

func (l *LuxContext) SetNoContent()

func (l *LuxContext) SetResetContent()

func (l *LuxContext) SetFound()

func (l *LuxContext) SetUnauthorized()

func (l *LuxContext) SetForbidden()

func (l *LuxContext) SetInternalServerError()

func (l *LuxContext) SetNotImplemented()

func (l *LuxContext) SetServiceUnavailable()

func (l *LuxContext) SetConflict()

func (l *LuxContext) SetUnsupportedMediaType()
```

## set cookie

### full option cookie

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		lc.SetCookie("key", "value", 0, "/", "", false, false)
		return nil
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### secure cookie

```go
package main

import (
	"time"

	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		lc.SetSecureCookie("key", "value", int(time.Hour.Seconds()), "/", "")
		return nil
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

## get

### get form file

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.POST("/", func(lc *context.LuxContext) error {
		file, header, err := lc.GetFormFile("filename")
		if err != nil {
			return err
		}
		defer file.Close()
		return lc.SaveFile(file, "./"+header.Filename)
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### get multipart file

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.POST("/", func(lc *context.LuxContext) error {
		headers, err := lc.GetMultipartFile("filename", 10<<20)
		if err != nil {
			return err
		}
		return lc.SaveMultipartFile(headers, "./"+headers[0].Filename)
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### get form data

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.POST("/", func(lc *context.LuxContext) error {
		value := lc.GetFormData("key")
		return lc.ReplyString(value)
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

The data from post form.

### get url query

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		value := lc.GetURLQuery("key")
		return lc.ReplyString(value)
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

`http://localhost:8080/?key=value`

### get path variable

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/:key", func(lc *context.LuxContext) error {
		value := lc.GetPathVariable("key")
		return lc.ReplyString(value)
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

`http://localhost:8080/value`

### get body

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.POST("/", func(lc *context.LuxContext) error {
		body, err := lc.GetBody()
		if err != nil {
			return err
		}
		return lc.ReplyBinary(body)
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

The data from http body.

### get body reader

```go
package main

import (
	"io"

	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.POST("/", func(lc *context.LuxContext) error {
		reader := lc.GetBodyReader()
        defer reader.Close()
		_, err := io.Copy(lc.Response, reader)
		return err
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### get cookie

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.POST("/", func(lc *context.LuxContext) error {
		value := lc.GetCookie("key")
		return lc.ReplyString(value)
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### get remote address, ip, port

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.POST("/", func(lc *context.LuxContext) error {
		addr := lc.GetRemoteAddress()
		ip := lc.GetRemoteIP()
		port := lc.GetRemotePort()
		return lc.ReplyJSON(map[string]string{
			"address": addr,
			"ip":      ip,
			"port":    port,
		})
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

## logext

### set stderr

```go
package main

import (
	"os"

	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/logext"
)

func main() {
	app := lux.New()

	logger := logext.New(os.Stderr)
	app.SetLogger(logger)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### set file

```go
package main

import (
	"os"

	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/logext"
)

func main() {
	app := lux.New()

	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	logger := logext.New(file)
	app.SetLogger(logger)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### change default logger's buffer size

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/logext"
	"github.com/snowmerak/lux/logext/stdout"
)

func main() {
	app := lux.New()

	bufferSize := 16
	logger := logext.New(stdout.New(bufferSize))
	app.SetLogger(logger)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

## middleware

### allow static ip

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New(
		middleware.AllowStaticIPs(
			"localhost",
			"127.0.0.1",
			"[::1]",
		),
	)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

allowing `localhost`, `127.0.0.1`, `[::1]`.

### allow dynamic ip

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	ipMap := map[string]bool{
		"localhost": true,
		"127.0.0.1": true,
		"[::1]":     true,
	}

	app := lux.New(
		middleware.AllowDynamicIPs(
			func(remoteIP string) bool {
				return ipMap[remoteIP]
			},
		),
	)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### block static ip

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New(
		middleware.BlockStaticIPs(
			"localhost",
			"127.0.0.1",
			"[::1]",
		),
	)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

blocking `localhost`, `127.0.0.1`, `[::1]`.

### block dynamic ip

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	ipMap := map[string]bool{
		"localhost": true,
		"127.0.0.1": true,
		"[::1]":     true,
	}

	app := lux.New(
		middleware.BlockDynamicIPs(
			func(remoteIP string) bool {
				return ipMap[remoteIP]
			},
		),
	)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### allow and block ports

```go
func AllowStaticPorts(ports ...string) Set

func BlockStaticPorts(ports ...string) Set 

func AllowDynamicPorts(checker func(remotePort string) bool) Set

func BlockDynamicPorts(checker func(remotePort string) bool) Set
```

### authorize

```go
package main

import (
	"net/http"

	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/", middleware.Auth(func(authorizaionHeader string, tokenCookie *http.Cookie) bool {
		return true
	}))
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("authorized")
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### compress snappy

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("hello!")
	}, middleware.CompressSnappy)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

If `Accept-Encoding` do not has `snappy`, ignore this middleware.

### compress gzip

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("hello!")
	}, middleware.CompressGzip)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

If `Accept-Encoding` do not has `gzip`, ignore this middleware.

### compress brotli

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/")
	rootGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("hello!")
	}, middleware.CompressBrotli)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

If `Accept-Encoding` do not has `brotli`, ignore this middleware.

### allow headers

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New(
		middleware.SetAllowHeaders("*"),
	)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### allow methods

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New(
		middleware.SetAllowMethods("*"),
	)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### allow origins

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New(
		middleware.SetAllowOrigins("*"),
	)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### allow credentials

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New(
		middleware.SetAllowCredentials,
	)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### allow cors

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New(
		middleware.SetAllowCORS,
	)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

## router

### http methods

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/context"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/", middleware.SetAllowCORS)

	rootGroup.GET("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("GET request")
	})

	rootGroup.POST("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("POST request")
	})

	rootGroup.PATCH("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("PATCH request")
	})

	rootGroup.PUT("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("PUT request")
	})

	rootGroup.DELETE("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("DELETE request")
	})

	rootGroup.OPTIONS("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("OPTIONS request")
	})
	
	rootGroup.HEAD("/", func(lc *context.LuxContext) error {
		return lc.ReplyString("HEAD request") // will be ignored
	})

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### preflight

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/", middleware.SetAllowCORS)

	rootGroup.Preflight(
		[]string{"*"},
		[]string{"*"},
		[]string{"*"},
	)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### statics

```go
package main

import (
	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/", middleware.SetAllowCORS)

	rootGroup.Statics("/public", "./public")

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```

### embed

```go
package main

import (
	"embed"
	"io/fs"

	"github.com/snowmerak/lux"
	"github.com/snowmerak/lux/middleware"
)

//go:embed public
var public embed.FS

func main() {
	app := lux.New()

	rootGroup := app.NewRouterGroup("/", middleware.SetAllowCORS)

	publicFS, err := fs.Sub(public, "public")
	if err != nil {
		panic(err)
	}
	rootGroup.Embedded("/", publicFS)

	if err := app.ListenAndServe1(":8080"); err != nil {
		panic(err)
	}
}
```
