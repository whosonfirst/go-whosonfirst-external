package properties

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var area_parquet string
var whosonfirst_parquet string
var reader_uri string
var with_spatial_geom bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("flags")

	fs.StringVar(&area_parquet, "area-parquet", "", "The URI for the \"area\" parquet file (produced by by the `compile-area` tool) from which Who's On First properties will be derived.")
	fs.StringVar(&whosonfirst_parquet, "whosonfirst-parquet", "", "The URI for the parquet file where Who's On First properties will be written to.")
	fs.StringVar(&reader_uri, "reader-uri", "https://data.whosonfirst.org", "A registered whosonfirst/go-reader.Reader URI.")
	fs.BoolVar(&with_spatial_geom, "with-spatial-geom", false, "Store geometry property as spatial GEOMETRY type (rather than TEXT.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Derive Who's On First properties for an \"area\" parquet file (produced by by the `compile-area` tool).\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
