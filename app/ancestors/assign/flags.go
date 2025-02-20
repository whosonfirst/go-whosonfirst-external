package assign

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var spatial_database_uri string
var properties_reader_uri string

var iterator_uri string
var workers int
var start_after int64

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("assign")

	fs.StringVar(&spatial_database_uri, "spatial-database-uri", "", "A registered whosonfirst/go-whosonfirst-spatial/database/SpatialDatabase URI to use for perforning reverse geocoding tasks.")

	fs.StringVar(&properties_reader_uri, "properties-reader-uri", "{spatial-database-uri}", "A registered whosonfirst/go-reader.Reader URI for reading properties from parent records. If '{spatial-database-uri}' the spatial database instance will be used to read those properties.")

	fs.StringVar(&iterator_uri, "iterator-uri", "", "A registered whosonfirst/go-whosonfirst-external/iterator.Iterator URI.")
	fs.IntVar(&workers, "workers", 5, "The maximum number of workers to process reverse geocoding tasks.")

	fs.Int64Var(&start_after, "start-after", 0, "If > 0 then delay processing for 'start_after' number of records.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Iterate through one or more URIs for an external data source and reverse-geocode each record emitting the record ID, Who's On First parent ID and Who's On First ancestry as CSV data to STDOUT.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s uri(N) uri(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
