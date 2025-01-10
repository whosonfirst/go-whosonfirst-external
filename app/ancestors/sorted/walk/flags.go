package walk

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var verbose bool
var geohash string
var parent_ids multi.MultiInt64
var ancestor_ids multi.MultiInt64
var mode string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("flags")

	fs.StringVar(&geohash, "geohash", "", "An optional geohash to do a prefix-first comparison against.")
	fs.Var(&parent_ids, "parent-id", "One or more \"wof:parent_id\" values to match.")
	fs.Var(&ancestor_ids, "ancestor-id", "Zero or more \"wof:hierarchies\" values to match.")
	fs.StringVar(&mode, "mode", mode_all, "Indicate whether all or any filter criteria must match. Valid options are: any, all.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Walk one or more whosonfirst-external-* repositories, with optional filtering, emitting CSV-encoded rows to STDOUT.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
