package iterate

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"sync/atomic"
	"time"

	"github.com/whosonfirst/go-whosonfirst-external"
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

	iter, err := external.NewIterator(ctx, opts.IteratorURI)

	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)

	t1 := time.Now()
	count := int64(0)

	defer func() {
		slog.Info("Time to iterate records", "count", count, "time", time.Since(t1))
	}()

	for r, err := range iter.Iterate(ctx, opts.IteratorSources...) {

		if err != nil {
			return err
		}

		if opts.AsGeoJSON {

			f, err := external.AsGeoJSONFeature(r)

			if err != nil {
				return err
			}

			err = enc.Encode(f)

			if err != nil {
				return err
			}

		} else {

			err = enc.Encode(r)

			if err != nil {
				return err
			}
		}

		atomic.AddInt64(&count, 1)
	}

	return nil
}
