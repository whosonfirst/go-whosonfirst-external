package assign

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	SpatialDatabaseURI string
	PropertiesReaderURI string
	IteratorURI string
	IteratorSources []string
	Workers int
	StartAfter int64
	Verbose bool
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	sources := fs.Args()

	opts := &RunOptions{
		SpatialDatabaseURI: spatial_database_uri,
		PropertiesReaderURI: properties_reader_uri,
		IteratorURI: iterator_uri,
		IteratorSources: sources,
		Workers: workers,
		StartAfter: start_after,
		Verbose: verbose,
	}

	return opts, nil
}
