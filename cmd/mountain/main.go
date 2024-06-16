package main

import (
	"context"
	"log"
	"os"

	"lua-mountain/internal/mountain"
)

var (
	version = "dev"
)

func main() {
	ctx := context.Background()
	mountain.Engine.Version = version
	if err := mountain.Start(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
