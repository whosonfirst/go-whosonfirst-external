package iterate

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	IteratorURI     string
	IteratorSources []string
	AsGeoJSON       bool
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	sources := fs.Args()

	opts := &RunOptions{
		IteratorURI:     iterator_uri,
		IteratorSources: sources,
		AsGeoJSON:       as_geojson,
	}

	return opts, nil
}
