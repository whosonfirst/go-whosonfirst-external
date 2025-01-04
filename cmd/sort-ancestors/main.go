package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-external/app/ancestors/sort"
)

func main() {

	ctx := context.Background()
	err := sort.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to sort ancestors, %v", err)
	}
}
