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
	fs.StringVar(&namespace, "namespace", "", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "...\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s uri(N) uri(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
