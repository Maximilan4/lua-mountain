package main

import (
	"context"
	"log"
	"lua-mountain/internal/mountain"
	"os"
)



func main() {
	ctx := context.Background()
	if err := mountain.Start(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
