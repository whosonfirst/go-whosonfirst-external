package iterator

import (
	"context"
	"database/sql"
	"fmt"
	"iter"
	"log/slog"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/paulmach/orb"
	"github.com/whosonfirst/go-whosonfirst-external"
)

type FoursquareIterator struct {
	Iterator
	db *sql.DB
}

func init() {
	ctx := context.Background()
	err := RegisterIterator(ctx, "foursquare", NewFoursquareIterator)
	if err != nil {
		panic(err)
	}
}

func NewFoursquareIterator(ctx context.Context, uri string) (Iterator, error) {

	engine := "duckdb"
	dsn := ""

	db, err := sql.Open(engine, dsn)

	if err != nil {
		return nil, err
	}

	it := &FoursquareIterator{
		db: db,
	}

	return it, nil
}

func (it *FoursquareIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[*external.Record, error] {

	return func(yield func(*external.Record, error) bool) {

		for _, uri := range uris {

			for r, err := range it.iterate(ctx, uri) {

				if !yield(r, err) {
					return
				}
			}
		}
	}
}

func (it *FoursquareIterator) iterate(ctx context.Context, uri string) iter.Seq2[*external.Record, error] {

	// SELECT fsq_place_id, name, ifnull(latitude, 0.0), ifnull(longitude, 0.0) FROM read_parquet('/usr/local/data/foursquare/parquet/*.parquet')

	logger := slog.Default()
	logger = logger.With("uri", uri)

	return func(yield func(*external.Record, error) bool) {

		q := fmt.Sprintf(`SELECT fsq_place_id, name, ifnull(latitude, 0.0), ifnull(longitude, 0.0) FROM read_parquet('%s')`, uri)

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
			var lat float64
			var lon float64

			err := rows.Scan(&id, &name, &lat, &lon)

			if err != nil {
				logger.Error("Failed to scan row", "error", err)
				yield(nil, err)
				return
			}

			pt := orb.Point([]float64{lon, lat})

			r := &external.Record{
				Id:        id,
				Name:      name,
				Placetype: "venue",
				Geometry:  &pt,
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

func (it *FoursquareIterator) Close() error {
	return it.db.Close()
}
