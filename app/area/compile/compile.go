package compile

/*

> go run cmd/compile-area/main.go -external-id-key fsq_place_id -external-source "/usr/local/data/foursquare/parquet/*.parquet" -mode any -ancestor-id 102087579 -ancestor-id 102086959 -ancestor-id 102085387 -target sfba.parquet /usr/local/data/foursquare/whosonfirst/whosonfirst-external-foursquare-venue-us/data/85688637/

*/

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

	db, err := sql.Open("duckdb", "")

	if err != nil {
		return fmt.Errorf("Failed to open database connection, %w", err)
	}

	defer db.Close()

	tmp_wr, err := os.CreateTemp("", "area.*.csv")

	if err != nil {
		return fmt.Errorf("Failed to create tmp file for area data, %w", err)
	}

	tmp_path := tmp_wr.Name()
	defer os.Remove(tmp_path)

	compile_opts := &whosonfirst.CompileOptions{
		Geohash:     opts.Geohash,
		ParentIds:   opts.ParentIds,
		AncestorIds: opts.AncestorIds,
		Mode:        opts.Mode,
		Sources:     opts.Sources,
	}

	err = whosonfirst.Compile(ctx, compile_opts, tmp_wr)

	if err != nil {
		return fmt.Errorf("Failed to compile area data, %w", err)
	}

	err = tmp_wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close tmp file for area data, %w", err)
	}

	q := fmt.Sprintf(`COPY (SELECT e.*, w.geohash, w."wof:country", w."wof:parent_id", w."wof:hierarchies" FROM read_parquet('%s') e, read_csv('%s') w WHERE e.%s = w."external:id") TO '%s' (COMPRESSION ZSTD)`, opts.ExternalSource, tmp_path, opts.ExternalIdKey, opts.Target)

	_, err = db.ExecContext(ctx, q)

	if err != nil {
		return fmt.Errorf("Failed to create area parquet, %w", err)
	}

	return nil
}
