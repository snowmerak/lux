package main

import (
	"context"
	"log"

	"github.com/snowmerak/lux/v4/pkg/lux"
)

func main() {
	app := lux.New()

	log.Println("Lux Knowledge MCP Server starting...")
	if err := app.Run(context.Background(), ""); err != nil {
		log.Fatal(err)
	}
}
