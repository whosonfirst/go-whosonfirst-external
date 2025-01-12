package whosonfirst

import (
	"compress/bzip2"
	"context"
	"fmt"
	"io"
	"io/fs"
	"iter"
	_ "log/slog"
	"os"

	"github.com/sfomuseum/go-csvdict/v2"
)

// This should really be an interface but today it is not...

func Read(ctx context.Context, uri string, match_opts *RowHasMatchOptions) iter.Seq2[map[string]string, error] {

	return func(yield func(map[string]string, error) bool) {

		// Something something something interface read from HTTP something something something...

		info, err := os.Stat(uri)

		if err != nil {
			yield(nil, err)
			return
		}

		var read_iter iter.Seq2[map[string]string, error]

		if info.IsDir() {
			read_iter = ReadDir(ctx, uri, match_opts)
		} else {
			read_iter = ReadFile(ctx, uri, match_opts)
		}

		for row, err := range read_iter {
			if !yield(row, err) {
				return
			}
		}
	}
}

func ReadDir(ctx context.Context, uri string, match_opts *RowHasMatchOptions) iter.Seq2[map[string]string, error] {

	return func(yield func(map[string]string, error) bool) {

		to_walk := os.DirFS(uri)

		err := fs.WalkDir(to_walk, ".", func(path string, d fs.DirEntry, err error) error {

			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			r, err := to_walk.Open(path)

			if err != nil {
				return err
			}

			defer r.Close()

			for row, err := range read(ctx, r, match_opts) {
				if !yield(row, err) {
					return nil
				}
			}

			return nil
		})

		if err != nil {
			yield(nil, err)
			return
		}
	}
}

func ReadFile(ctx context.Context, uri string, match_opts *RowHasMatchOptions) iter.Seq2[map[string]string, error] {

	return func(yield func(map[string]string, error) bool) {

		r, err := os.Open(uri)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to open %s for reading, %w", uri, err))
			return
		}

		defer r.Close()

		for row, err := range read(ctx, r, match_opts) {
			if !yield(row, err) {
				return
			}
		}
	}
}

func read(ctx context.Context, r io.Reader, match_opts *RowHasMatchOptions) iter.Seq2[map[string]string, error] {

	return func(yield func(map[string]string, error) bool) {

		br := bzip2.NewReader(r)

		csv_r, err := csvdict.NewReader(br)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to create CSV reader, %w", err))
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
