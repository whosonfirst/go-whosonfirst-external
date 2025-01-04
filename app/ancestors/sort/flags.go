package sort

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var target string
var with_geohash bool
var geohash_precision int
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("assign")

	fs.StringVar(&target, "target", "-", "If target is '-' then all data will be written to /dev/null (or equivalent).")
	fs.BoolVar(&with_geohash, "with-geohash", true, "...")
	fs.IntVar(&geohash_precision, "geohash-precision", 3, "...")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "...\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s uri(N) uri(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
