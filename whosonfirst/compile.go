package whosonfirst

import (
	"context"
	"io"
	"log/slog"
	"sync"

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

	mu := new(sync.RWMutex)

	done_ch := make(chan bool)
	err_ch := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, uri := range opts.Sources {

		go func(uri string) {

			defer func() {
				slog.Info("Finished compiling", "uri", uri)
				done_ch <- true
			}()

			for row, err := range Read(ctx, uri, match_opts) {

				if err != nil {
					err_ch <- err
					return
				}

				select {
				case <-ctx.Done():
					return
				default:
					// pass
				}

				mu.Lock()
				// defer mu.Unlock()

				if csv_wr == nil {

					new_wr, err := csvdict.NewWriter(wr)

					if err != nil {
						mu.Unlock()
						err_ch <- err
						return
					}

					csv_wr = new_wr
				}

				err = csv_wr.WriteRow(row)

				if err != nil {
					mu.Unlock()
					err_ch <- err
					return
				}

				csv_wr.Flush()
				mu.Unlock()
			}

		}(uri)
	}

	remaining := len(opts.Sources)

	for remaining > 0 {
		select {
		case <-done_ch:
			remaining -= 1
		case err := <-err_ch:
			return err
		}
	}

	return nil
}
