package walk

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Verbose     bool
	Sources     []string
	GeoHash     string
	ParentIds   []int64
	AncestorIds []int64
	Mode string
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	sources := fs.Args()

	opts := &RunOptions{
		Verbose:     verbose,
		Sources:     sources,
		GeoHash:     geohash,
		ParentIds:   parent_ids,
		AncestorIds: ancestor_ids,
		Mode: mode,
	}

	return opts, nil
}
