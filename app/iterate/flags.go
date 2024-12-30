package iterate

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var iterator_uri string
var as_geojson bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("iterate")

	fs.StringVar(&iterator_uri, "iterator-uri", "", "A registered whosonfirst/go-whosonfirst-external/iterator.")
	fs.BoolVar(&as_geojson, "as-geojson", false, "Emit records as GeoJSON Features.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Iterate through one or more URIs for an external data source and emit each record as line-separated JSON.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s uri(N) uri(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
