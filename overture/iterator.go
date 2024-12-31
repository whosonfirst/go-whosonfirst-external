package overture

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"iter"
	"log/slog"
	"net/url"
	"strings"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/paulmach/orb/geojson"
	sfom_sql "github.com/sfomuseum/go-database/sql"
	"github.com/whosonfirst/go-whosonfirst-external"
)

const all_properties string = "all"

type OvertureIterator struct {
	external.Iterator
	db         *sql.DB
	properties string
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

	q := u.Query()

	properties := q.Get("properties")

	engine := "duckdb"
	dsn := ""

	db, err := sql.Open(engine, dsn)

	if err != nil {
		return nil, err
	}

	err = sfom_sql.LoadDuckDBExtentions(ctx, db, "SPATIAL")

	if err != nil {
		return nil, err
	}

	it := &OvertureIterator{
		db:         db,
		properties: properties,
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

		props := []string{
			"id",
			"names.primary AS name",
			"ST_AsGeoJSON(geometry) AS geometry",
		}

		if it.properties == all_properties {

			/*

						Wut... ??

			> go run cmd/iterate/main.go -iterator-uri 'overture://parquet/places?properties=all' ~/data/overture/parquet/*.parquet
			2024/12/31 09:32:04 ERROR Failed to scan row uri=/Users/asc/data/overture/parquet/part-00000-9b3cb01a-46a1-4378-9e77-baca19283b5a-c000.zstd.parquet error="sql: Scan error on column index 5, name \"\\\"json\\\"(\\\"names\\\")\": destination not a pointer"
			2024/12/31 09:32:04 INFO Time to iterate records count=0 time=2.324170458s
			2024/12/31 09:32:04 Failed to iterate records, sql: Scan error on column index 5, name "\"json\"(\"names\")": destination not a pointer
			exit status 1

			*/

			other_props := []string{
				"version",
				"JSON(sources)",
				"JSON(names)",
				"JSON(categories)",
				"confidence",
				"JSON(websites)",
				"JSON(socials)",
				"JSON(emails)",
				"JSON(phones)",
				"JSON(brand)",
				"JSON(addresses)",
			}

			props = append(props, other_props...)
		}

		str_props := strings.Join(props, ",")

		q := fmt.Sprintf(`SELECT %s FROM read_parquet('%s')`, str_props, uri)
		rows, err := it.db.QueryContext(ctx, q)

		if err != nil {
			logger.Error("Failed to execute query", "error", err)
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {

			var props map[string]any
			var str_geom string

			if it.properties == all_properties {

				var id string
				var name string
				// str_geom defined above
				var version int
				var str_sources string
				var str_names string
				var str_categories string
				var confidence float32
				var str_websites string
				var str_socials string
				var str_emails string
				var str_phones string
				var str_brand string
				var str_addresses string

				err := rows.Scan(
					&id, &name, &str_geom, &version,
					&str_sources, str_names, &str_categories,
					&confidence,
					&str_websites, &str_socials, &str_emails, &str_phones,
					&str_brand, &str_addresses,
				)

				if err != nil {
					logger.Error("Failed to scan row", "error", err)
					yield(nil, err)
					return
				}

				logger = logger.With("id", id)

				// To do: Add types mapped to the Overture schema

				var sources map[string]any
				var names map[string]any
				var categories map[string]any

				var websites []string
				var socials []string
				var emails []string
				var phones []string

				var brand map[string]any
				var addresses map[string]any

				err = json.Unmarshal([]byte(str_sources), &sources)

				if err != nil {
					logger.Error("Failed to unmarshal sources", "error", err)
					yield(nil, err)
					return
				}

				err = json.Unmarshal([]byte(str_names), &names)

				if err != nil {
					logger.Error("Failed to unmarshal names", "error", err)
					yield(nil, err)
					return
				}

				err = json.Unmarshal([]byte(str_categories), &categories)

				if err != nil {
					logger.Error("Failed to unmarshal categories", "error", err)
					yield(nil, err)
					return
				}

				err = json.Unmarshal([]byte(str_websites), &websites)

				if err != nil {
					logger.Error("Failed to unmarshal websites", "error", err)
					yield(nil, err)
					return
				}

				err = json.Unmarshal([]byte(str_socials), &socials)

				if err != nil {
					logger.Error("Failed to unmarshal socials", "error", err)
					yield(nil, err)
					return
				}

				err = json.Unmarshal([]byte(str_emails), &emails)

				if err != nil {
					logger.Error("Failed to unmarshal emails", "error", err)
					yield(nil, err)
					return
				}

				err = json.Unmarshal([]byte(str_phones), &phones)

				if err != nil {
					logger.Error("Failed to unmarshal phones", "error", err)
					yield(nil, err)
					return
				}

				err = json.Unmarshal([]byte(str_brand), &brand)

				if err != nil {
					logger.Error("Failed to unmarshal brands", "error", err)
					yield(nil, err)
					return
				}

				err = json.Unmarshal([]byte(str_addresses), &addresses)

				if err != nil {
					logger.Error("Failed to unmarshal addresses", "error", err)
					yield(nil, err)
					return
				}

				props = map[string]any{
					"id":         id,
					"name":       name,
					"version":    version,
					"confidence": confidence,
					"sources":    sources,
					"names":      names,
					"categories": categories,
					"websites":   websites,
					"socials":    socials,
					"emails":     emails,
					"phones":     phones,
					"brand":      brand,
					"addresses":  addresses,
				}

			} else {

				var id string
				var name string
				// str_geom defined above

				err := rows.Scan(&id, &name, &str_geom)

				if err != nil {
					logger.Error("Failed to scan row", "error", err)
					yield(nil, err)
					return
				}

				props = map[string]any{
					"id":   id,
					"name": name,
				}

			}

			str_f := fmt.Sprintf(`{ "type": "Feature", "properties": {}, "geometry": %s }`, str_geom)

			f, err := geojson.UnmarshalFeature([]byte(str_f))

			if err != nil {
				logger.Error("Failed to unmarshal geometry in to feature", "error", err)
				yield(nil, err)
				return
			}

			r, err := NewOvertureRecord(props, f.Geometry)

			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(r, nil) {
				logger.Error("Failed to yield record", "id", props["id"], "error", err)
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
