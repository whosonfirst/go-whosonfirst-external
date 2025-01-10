package sort

import (
	"compress/bzip2"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
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

	namespace := opts.Namespace

	csv_writers := make(map[string]*csvdict.Writer)
	writers := make([]io.WriteCloser, 0)

	for _, path_r := range opts.Sources {

		logger := slog.Default()
		logger = logger.With("source", path_r)

		r, err := os.Open(path_r)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", path_r, err)
		}

		defer r.Close()

		bz_r := bzip2.NewReader(r)

		csv_r, err := csvdict.NewReader(bz_r)

		if err != nil {
			return fmt.Errorf("Failed to create CSV reader for %s, %w", path_r, err)
		}

		for row, err := range csv_r.Iterate() {

			if err != nil {
				return fmt.Errorf("CSV reader yielded an error, %w", err)
			}

			to_write := make([]string, 0)

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

			/*
				if country != "GB" {
					continue
				}
			*/

			country = strings.ToLower(country)

			if namespace == "" {

				ns, ok := row["external:namespace"]

				if !ok {
					logger.Warn("Row is missing external:namespace", "row", row)
					continue
				}

				namespace = ns
			}

			hier_paths := make([]string, 0)

			var hierarchies []map[string]int64

			err := json.Unmarshal([]byte(row["wof:hierarchies"]), &hierarchies)

			if err != nil {

				hier_paths = []string{
					fmt.Sprintf("xx/%s-xx-xx.csv", country),
				}

			} else {

				for _, hier := range hierarchies {

					str_region := "xx"
					str_locality := "xx"

					region_id, region_ok := hier["region_id"]
					locality_id, locality_ok := hier["locality_id"]

					if region_ok && region_id != -1 {
						str_region = strconv.FormatInt(region_id, 10)
					}

					if locality_ok && locality_id != -1 {
						str_locality = strconv.FormatInt(locality_id, 10)
					}

					fname := fmt.Sprintf("%s-%s-%s.csv", country, str_region, str_locality)
					rel_path := filepath.Join(str_region, fname)

					if !slices.Contains(hier_paths, rel_path) {
						hier_paths = append(hier_paths, rel_path)
					}
				}
			}

			for _, rel_path := range hier_paths {

				root := filepath.Join(opts.Target, fmt.Sprintf("whosonfirst-external-%s-venue-%s", namespace, country))
				root = filepath.Join(root, "data")

				wr_path := filepath.Join(root, rel_path)
				to_write = append(to_write, wr_path)
			}

			for _, path := range to_write {

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

						wr, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)

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
	}

	for _, wr := range writers {
		wr.Close()
	}

	return nil
}
