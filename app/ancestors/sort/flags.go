package sort

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var target string
var namespace string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("assign")

	fs.StringVar(&target, "target", "-", "If target is '-' then all data will be written to /dev/null (or equivalent).")
	fs.StringVar(&namespace, "namespace", "", "The namespace of the external source.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Sort CSV data containing external sources and their ancestors (produced by assign-ancestors) in to nested {REGION_ID}/{COUNTY_CODE}-{REGION_ID}-{LOCALITY_ID} CSV files\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s uri(N) uri(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
