package walk

import (
	"compress/bzip2"
	"context"
	"flag"
	"io/fs"
	"iter"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/sfomuseum/go-csvdict/v2"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	for _, path := range opts.Sources {

		to_walk := os.DirFS(path)

		for row, err := range walkFS(ctx, to_walk) {

			if err != nil {
				return err
			}

			slog.Info("row", "row", row)
		}
	}

	return nil
}

func walkFS(ctx context.Context, to_walk fs.FS) iter.Seq2[map[string]string, error] {

	return func(yield func(map[string]string, error) bool) {

		err := fs.WalkDir(to_walk, ".", func(path string, d fs.DirEntry, err error) error {

			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			ext := filepath.Ext(path)

			if ext != ".bz2" {
				return nil
			}

			r, err := to_walk.Open(path)

			if err != nil {
				return err
			}

			defer r.Close()

			br := bzip2.NewReader(r)

			csv_r, err := csvdict.NewReader(br)

			if err != nil {
				return err
			}

			for row, err := range csv_r.Iterate() {

				if err != nil {
					return err
				}

				if !yield(row, nil) {
					break
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
