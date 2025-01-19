package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-external/app/area/whosonfirst/properties"
)

func main() {

	ctx := context.Background()
	err := properties.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to derive who's on first properties for area, %v", err)
	}
}
