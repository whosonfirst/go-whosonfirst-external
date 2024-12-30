package iterate

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var iterator_uri string
var as_geojson bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("iterate")

	fs.StringVar(&iterator_uri, "iterator-uri", "", "A registered whosonfirst/go-whosonfirst-external/iterator.")
	fs.BoolVar(&as_geojson, "as-geojson", false, "Emit records as GeoJSON Features.")
	return fs
}
