package compile

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/whosonfirst/go-whosonfirst-external/whosonfirst"
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

	tmp_wr, err := os.CreateTemp("", "area.*.csv")

	if err != nil {
		return err
	}

	tmp_path := tmp_wr.Name()
	defer os.Remove(tmp_path)

	compile_opts := &whosonfirst.CompileOptions{
		Geohash:     opts.Geohash,
		ParentIds:   opts.ParentIds,
		AncestorIds: opts.AncestorIds,
		Mode:        opts.Mode,
	}

	err = whosonfirst.Compile(ctx, compile_opts, tmp_wr)

	if err != nil {
		return err
	}

	err = tmp_wr.Close()

	if err != nil {
		return err
	}

	db, err := sql.Open("duckdb", "")

	if err != nil {
		return err
	}

	defer db.Close()

	q := fmt.Sprintf(`COPY (SELECT f.*, w.geohash, w."wof:country", w."wof:parent_id", w."wof:hierarchies" FROM read_parquet('%s') f, read_csv('%s') w WHERE f.fsq_place_id = w."external:id") TO '%s' (COMPRESSION ZSTD)`, opts.ExternalSource, tmp_path, opts.Target)

	_, err = db.ExecContext(ctx, q)

	if err != nil {
		return err
	}

	return nil
}
