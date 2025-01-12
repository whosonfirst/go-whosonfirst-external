package compile

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Verbose        bool
	Sources        []string
	Geohash        string
	ParentIds      []int64
	AncestorIds    []int64
	Mode           string
	Target         string
	ExternalSource string
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		Verbose: verbose,
		Sources: fs.Args(),
	}

	return opts, nil
}
