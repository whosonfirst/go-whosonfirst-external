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
	ExternalIdKey  string
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		Target:         target,
		ExternalSource: external_source,
		ExternalIdKey:  external_id_key,
		Sources:        fs.Args(),
		Geohash:        geohash,
		ParentIds:      parent_ids,
		AncestorIds:    ancestor_ids,
		Mode:           mode,
		Verbose:        verbose,
	}

	return opts, nil
}
