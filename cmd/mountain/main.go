package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"lua-mountain/internal/mountain"
)

var (
	version = "dev"
)

func main() {
	ctx := context.Background()
	fmt.Println(version)
	if err := mountain.Start(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
