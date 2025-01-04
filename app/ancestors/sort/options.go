package sort

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Sources          []string
	Target           string
	WithGeohash      bool
	GeohashPrecision int
	Verbose          bool
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	sources := fs.Args()

	opts := &RunOptions{
		Target:           target,
		Sources:          sources,
		WithGeohash:      with_geohash,
		GeohashPrecision: geohash_precision,
		Verbose:          verbose,
	}

	return opts, nil
}
