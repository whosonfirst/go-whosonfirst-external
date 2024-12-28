package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/whosonfirst/go-reader-database-sql"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-pmtiles"
	
	"github.com/whosonfirst/go-whosonfirst-external/app/ancestors/assign"
)

func main() {

	ctx := context.Background()
	err := assign.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to assign ancestors, %v", err)
	}
}
