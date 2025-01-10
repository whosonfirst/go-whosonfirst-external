package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-external/app/ancestors/sorted/walk"
)

func main() {

	ctx := context.Background()
	err := walk.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to walk sorted, %v", err)
	}
}
