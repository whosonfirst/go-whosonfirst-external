package overture

import (
	"context"
	"database/sql"
	"fmt"
	"iter"
	"log/slog"
	"net/url"

	_ "github.com/marcboeker/go-duckdb"

	_ "github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/whosonfirst/go-whosonfirst-external"
)

type OvertureIterator struct {
	external.Iterator
	db *sql.DB
}

func init() {
	ctx := context.Background()
	err := external.RegisterIterator(ctx, "overture", NewOvertureIterator)
	if err != nil {
		panic(err)
	}
}

func NewOvertureIterator(ctx context.Context, uri string) (external.Iterator, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	if u.Host != "parquet" {
		return nil, fmt.Errorf("Unsupported data type")
	}

	if u.Path != "/places" {
		return nil, fmt.Errorf("Unsupported place type")
	}

	engine := "duckdb"
	dsn := ""

	db, err := sql.Open(engine, dsn)

	if err != nil {
		return nil, err
	}

	extensions := []string{
		"SPATIAL",
	}

	// START OF put me in sfomuseum/go-database

	for _, ext := range extensions {

		commands := []string{
			fmt.Sprintf("INSTALL %s", ext),
			fmt.Sprintf("LOAD %s", ext),
		}

		for _, cmd := range commands {

			_, err := db.ExecContext(ctx, cmd)

			if err != nil {
				return nil, err
			}
		}
	}

	// END OF put me in sfomuseum/go-database

	it := &OvertureIterator{
		db: db,
	}

	return it, nil
}

func (it *OvertureIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[external.Record, error] {

	return func(yield func(external.Record, error) bool) {

		for _, uri := range uris {

			for r, err := range it.iterate(ctx, uri) {

				if !yield(r, err) {
					return
				}
			}

		}
	}
}

func (it *OvertureIterator) iterate(ctx context.Context, uri string) iter.Seq2[external.Record, error] {

	// SELECT id, names.primary AS name, ST_AsGeoJSON(geometry) AS geometry FROM read_parquet('/usr/local/overture/parquet/*.parquet') LIMIT 5;

	logger := slog.Default()
	logger = logger.With("uri", uri)

	return func(yield func(external.Record, error) bool) {

		q := fmt.Sprintf(`SELECT id, names.primary AS name, ST_AsGeoJSON(geometry) AS geometry FROM read_parquet('%s')`, uri)
		rows, err := it.db.QueryContext(ctx, q)

		if err != nil {
			logger.Error("Failed to execute query", "error", err)
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {

			var id string
			var name string
			var geom string

			err := rows.Scan(&id, &name, &geom)

			if err != nil {
				logger.Error("Failed to scan row", "error", err)
				yield(nil, err)
				return
			}

			str_f := fmt.Sprintf(`{ "type": "Feature", "properties": {}, "geometry": %s }`, geom)

			f, err := geojson.UnmarshalFeature([]byte(str_f))

			if err != nil {
				logger.Error("Failed to unmarshal geometry in to feature", "error", err)
				yield(nil, err)
				return
			}

			props := map[string]any{
				"id":   id,
				"name": name,
			}

			r, err := NewOvertureRecord(props, f.Geometry)

			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(r, nil) {
				logger.Error("Failed to yield record", "id", id, "error", err)
				return
			}
		}

		err = rows.Err()

		if err != nil {
			logger.Error("Failed to scan rows", "error", err)
			yield(nil, err)
			return
		}

	}
}

func (it *OvertureIterator) Close() error {
	return it.db.Close()
}
