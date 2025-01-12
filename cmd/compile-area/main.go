package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-external/app/area/compile"
)

func main() {

	ctx := context.Background()
	err := compile.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to compile area, %v", err)
	}
}
