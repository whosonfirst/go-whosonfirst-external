package sort

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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

	csv_writers := make(map[string]*csvdict.Writer)
	writers := make([]io.WriteCloser, 0)

	for _, path_r := range opts.Sources {

		logger := slog.Default()
		logger = logger.With("source", path_r)

		csv_r, err := csvdict.NewReaderFromPath(path_r)

		if err != nil {
			return fmt.Errorf("Failed to create CSV reader for %s, %w", path_r, err)
		}

		// mu := new(sync.RWMutex)

		for row, err := range csv_r.Iterate() {

			if err != nil {
				return fmt.Errorf("CSV reader yielded an error, %w", err)
			}

			// logger.Debug("Data", "row", row)

			country, ok := row["wof:country"]

			if !ok {
				logger.Warn("Row is missing wof:country", "row", row)
				continue
				// return fmt.Errorf("Row is missing wof:country, %v", row)
			}

			if country == "" {
				country = "XY"
			}

			country = strings.ToLower(country)
			fname := fmt.Sprintf("%s.csv", country)

			if opts.WithGeohash {

				geohash, ok := row["geohash"]

				if !ok {
					logger.Warn("Row is missing geohash", "row", row)
					continue
					// return fmt.Errorf("Row is missing geohash")
				}

				fname = fmt.Sprintf("%s-%s.csv", country, geohash[0:opts.GeohashPrecision])
			}

			namespace := "ovtr" // FIX ME...

			root := filepath.Join(opts.Target, fmt.Sprintf("whosonfirst-data-external-%s-%s", namespace, country))
			root = filepath.Join(root, "data")

			path := filepath.Join(root, fname)

			// mu.Lock()
			// defer mu.Unlock()

			csv_wr, ok := csv_writers[path]

			if !ok {

				// slog.Debug("Create new CSV writer", "path", path)

				var new_csv_wr *csvdict.Writer
				var new_csv_err error

				if opts.Target == "-" {
					new_csv_wr, new_csv_err = csvdict.NewWriter(io.Discard)
				} else {

					path_root := filepath.Dir(path)
					err := os.MkdirAll(path_root, 0755)

					if err != nil {
						return fmt.Errorf("Failed to create %s, %w", path_root, err)
					}

					wr, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

					if err != nil {
						return fmt.Errorf("Failed to create writer for %s, %w", path, err)
					}

					writers = append(writers, wr)
					new_csv_wr, new_csv_err = csvdict.NewWriter(wr)
				}

				if new_csv_err != nil {
					return fmt.Errorf("Failed to create new CSV writer for %s, %w", path, new_csv_err)
				}

				csv_wr = new_csv_wr
				csv_writers[path] = csv_wr
			}

			err = csv_wr.WriteRow(row)

			if err != nil {
				return fmt.Errorf("Failed to write row to %s, %w", path, err)
			}

			csv_wr.Flush()
		}
	}

	for _, wr := range writers {
		wr.Close()
	}

	return nil
}
