package compile

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("flags")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "...\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s ...\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
