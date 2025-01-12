package whosonfirst

import (
	"compress/bzip2"
	"context"
	"fmt"
	"iter"
	"os"

	"github.com/sfomuseum/go-csvdict/v2"
)

func Read(ctx context.Context, uri string, match_opts *RowHasMatchOptions) iter.Seq2[map[string]string, error] {

	return func(yield func(map[string]string, error) bool) {

		r, err := os.Open(uri)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to open %s for reading, %w", uri, err))
			return
		}

		defer r.Close()

		br := bzip2.NewReader(r)

		csv_r, err := csvdict.NewReader(br)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to create CSV reader for %s, %w", uri, err))
			return
		}

		for row, err := range csv_r.Iterate() {

			if err != nil {
				yield(nil, err)
				return
			}

			has_match, err := RowHasMatch(ctx, row, match_opts)

			if err != nil {
				yield(nil, err)
				return
			}

			if !has_match {
				continue
			}

			if !yield(row, nil) {
				return
			}
		}
	}
}
