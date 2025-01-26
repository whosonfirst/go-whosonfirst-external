package properties

/*

Basically a faster alternative to this:

COPY (SELECT g.id, g.name,g.geometry  FROM read_parquet('https://data.geocode.earth/wof/dist/parquet/whosonfirst-data-admin-latest.parquet') g, read_parquet('/Users/asc/whosonfirst/whosonfirst-external-duckdb/www/data/foursquare-sfba.parquet') s WHERE g.id=JSON_EXTRACT_STRING("wof:hierarchies", '$[0].neighbourhood_id') OR g.id=JSON_EXTRACT_STRING("wof:hierarchies", '$[0].locality_id') GROUP BY g.id, g.name, g.geometry ) TO 'whosonfirst.parquet' (COMPRESSION ZSTD);

*/

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/whosonfirst/go-reader-http"

	sfom_sql "github.com/sfomuseum/go-database/sql"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-feature/geometry"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
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

	r, err := reader.NewReader(ctx, opts.ReaderURI)

	if err != nil {
		return fmt.Errorf("Failed to create new reader, %w", err)
	}

	if opts.WithSpatialGeom {

		err = sfom_sql.LoadDuckDBExtensions(ctx, db, "SPATIAL")

		if err != nil {
			return fmt.Errorf("Failed to load extensions, %w", err)
		}

		_, err = db.ExecContext(ctx, `CREATE TEMP TABLE whosonfirst ("id" INTEGER NOT NULL, "name" TEXT, "placetype" TEXT, "geometry" GEOMETRY)`)

		if err != nil {
			return fmt.Errorf("Failed to create temp whosonfirst table, %w", err)
		}

	} else {

		_, err = db.ExecContext(ctx, `CREATE TEMP TABLE whosonfirst ("id" INTEGER NOT NULL, "name" TEXT, "placetype" TEXT, "geometry" TEXT)`)

		if err != nil {
			return fmt.Errorf("Failed to create temp whosonfirst table, %w", err)
		}

	}

	queries := []string{
		fmt.Sprintf(`SELECT DISTINCT(IFNULL(JSON_EXTRACT("wof:hierarchies", '$[0].neighbourhood_id'), -1)) FROM read_parquet('%s')`, opts.AreaParquet),
		fmt.Sprintf(`SELECT DISTINCT(IFNULL(JSON_EXTRACT("wof:hierarchies", '$[0].locality_id'), -1)) FROM read_parquet('%s')`, opts.AreaParquet),
	}

	wof_ids := make([]int64, 0)

	for _, q := range queries {

		rows, err := db.QueryContext(ctx, q)

		if err != nil {
			return fmt.Errorf("Failed to query '%s', %w", q, err)
		}

		defer rows.Close()

		for rows.Next() {

			var id int64

			err := rows.Scan(&id)

			if err != nil {
				return fmt.Errorf("Failed to scan row for query '%s', %w", q, err)
			}

			if id > -1 {
				wof_ids = append(wof_ids, id)
			}
		}

		err = rows.Err()

		if err != nil {
			return err
		}
	}

	slog.Debug("Who's On First IDs to fetch", "count", len(wof_ids))
	
	for _, id := range wof_ids {

		logger := slog.Default()
		logger = logger.With("id", id)

		logger.Debug("Fetch ID")
		
		body, err := wof_reader.LoadBytes(ctx, r, id)

		if err != nil {
			logger.Error("Failed to load record", "error", err)
			return fmt.Errorf("Failed to load body for %d, %w", id, err)
		}

		name, err := properties.Name(body)

		if err != nil {
			logger.Error("Failed to derive name", "error", err)			
			return fmt.Errorf("Failed to derive name for %d, %w", id, err)
		}

		pt, err := properties.Placetype(body)

		if err != nil {
			logger.Error("Failed to derive placetype", "error", err)			
			return fmt.Errorf("Failed to derive placetype for %d, %w", id, err)
		}
		
		geom, err := geometry.Geometry(body)

		if err != nil {
			logger.Error("Failed to derive geometry", "error", err)
			return fmt.Errorf("Failed to derive geometry for %d, %w", id, err)
		}

		enc_geom, err := geom.MarshalJSON()

		if err != nil {
			logger.Error("Failed to marshal geometry", "error", err)			
			return fmt.Errorf("Failed to marshal geometry for %d, %w", id, err)
		}

		q := `INSERT INTO whosonfirst (id, name, placetype, geometry) VALUES(?, ?, ?, ?)`

		if opts.WithSpatialGeom {
			q = `INSERT INTO whosonfirst (id, name, placetype, geometry) VALUES(?, ?, ?, ST_GeomFromGeoJSON(?))`
		}

		_, err = db.ExecContext(ctx, q, id, name, pt, string(enc_geom))

		if err != nil {
			logger.Error("Failed to add record", "error", err)
			return fmt.Errorf("Failed to add row for %d, %w", id, err)
		}
	}

	slog.Debug("Copy temporary database", "target", opts.WhosOnFirstParquet)
	
	_, err = db.ExecContext(ctx, fmt.Sprintf(`COPY (SELECT * FROM whosonfirst) TO '%s' (COMPRESSION ZSTD)`, opts.WhosOnFirstParquet))

	if err != nil {
		slog.Error("Failed to copy database", "error", err)
		return fmt.Errorf("Failed to copy table to disk, %w", err)
	}

	return nil
}
