package iterate

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var iterator_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("iterate")

	fs.StringVar(&iterator_uri, "iterator-uri", "", "...")

	return fs
}
