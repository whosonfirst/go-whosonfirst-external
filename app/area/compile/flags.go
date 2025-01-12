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
var external_id_key string
var target string

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("flags")

	fs.StringVar(&external_source, "external-source", "", "The string to pass to the DuckDB 'read_parquet' command for reading an external data source.")
	fs.StringVar(&external_id_key, "external-id-key", "", "The name of the unique identifier key for an external data source. The output is a new GeoParquet file.xs")
	fs.StringVar(&target, "target", "", "The path where the final GeoParquet file should be written.")

	fs.StringVar(&geohash, "geohash", "", "An optional geohash to do a prefix-first comparison against.")
	fs.Var(&parent_ids, "parent-id", "One or more \"wof:parent_id\" values to match.")
	fs.Var(&ancestor_ids, "ancestor-id", "Zero or more \"wof:hierarchies\" values to match.")
	fs.StringVar(&mode, "mode", "all", "Indicate whether all or any filter criteria must match. Valid options are: any, all.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Merge GeoParquet data for an external data source with one or more Who's On First ancestry sources. Ancestry ources can either be individual whosonfirst-external-* bzip2-compressed CSV files or folders containing one or more bzip2-compressed CSV files.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s uri(N) uri(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
