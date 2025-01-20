package properties

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Verbose            bool
	AreaParquet        string
	WhosOnFirstParquet string
	ReaderURI          string
	WithSpatialGeom    bool
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		AreaParquet:        area_parquet,
		WhosOnFirstParquet: whosonfirst_parquet,
		ReaderURI:          reader_uri,
		WithSpatialGeom:    with_spatial_geom,
		Verbose:            verbose,
	}

	return opts, nil
}
