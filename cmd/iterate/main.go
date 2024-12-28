package main

import (
	"context"
	"log"

	_ "github.com/whosonfirst/go-whosonfirst-external/foursquare"
	_ "github.com/whosonfirst/go-whosonfirst-external/overture"

	"github.com/whosonfirst/go-whosonfirst-external/app/iterate"
)

func main() {

	ctx := context.Background()
	err := iterate.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to iterate records, %v", err)
	}
}
