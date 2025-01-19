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
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		AreaParquet:        area_parquet,
		WhosOnFirstParquet: whosonfirst_parquet,
		ReaderURI:          reader_uri,
		Verbose:            verbose,
	}

	return opts, nil
}
