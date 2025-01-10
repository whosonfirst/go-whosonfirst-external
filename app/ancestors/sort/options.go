package sort

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	Sources   []string
	Target    string
	Namespace string
	Verbose   bool
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	sources := fs.Args()

	opts := &RunOptions{
		Target:    target,
		Sources:   sources,
		Namespace: namespace,
		Verbose:   verbose,
	}

	return opts, nil
}
