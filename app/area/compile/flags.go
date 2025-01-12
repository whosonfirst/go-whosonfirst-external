package compile

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var geohash string
var parent_ids multi.MultiInt64
var ancestor_ids multi.MultiInt64
var mode string

var external_source string
var target string

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("flags")

	fs.StringVar(&external_source, "external-source", "", "...")
	fs.StringVar(&target, "target", "", "...")

	fs.StringVar(&geohash, "geohash", "", "An optional geohash to do a prefix-first comparison against.")
	fs.Var(&parent_ids, "parent-id", "One or more \"wof:parent_id\" values to match.")
	fs.Var(&ancestor_ids, "ancestor-id", "Zero or more \"wof:hierarchies\" values to match.")
	fs.StringVar(&mode, "mode", "all", "Indicate whether all or any filter criteria must match. Valid options are: any, all.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "...\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s ...\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
