package iterate

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var iterator_uri string
var as_geojsonl bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("iterate")

	fs.StringVar(&iterator_uri, "iterator-uri", "", "...")
	fs.BoolVar(&as_geojsonl, "as-geojsonl", false, "...")
	return fs
}
