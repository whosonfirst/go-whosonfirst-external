package whosonfirst

import (
	"context"
	"io"

	"github.com/sfomuseum/go-csvdict/v2"
)

type CompileOptions struct {
	Sources     []string
	Geohash     string
	ParentIds   []int64
	AncestorIds []int64
	Mode        string
}

func Compile(ctx context.Context, opts *CompileOptions, wr io.Writer) error {

	match_opts := &RowHasMatchOptions{
		Geohash:     opts.Geohash,
		ParentIds:   opts.ParentIds,
		AncestorIds: opts.AncestorIds,
		Mode:        opts.Mode,
	}

	var csv_wr *csvdict.Writer

	for _, uri := range opts.Sources {

		for row, err := range Read(ctx, uri, match_opts) {

			if err != nil {
				return err
			}

			if csv_wr == nil {

				new_wr, err := csvdict.NewWriter(wr)

				if err != nil {
					return err
				}

				csv_wr = new_wr
			}

			err = csv_wr.WriteRow(row)

			if err != nil {
				return err
			}

			csv_wr.Flush()
		}

	}

	return nil
}
