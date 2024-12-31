package foursquare

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

	"github.com/paulmach/orb"
	"github.com/whosonfirst/go-whosonfirst-external"
)

const all_properties string = "all"

type FoursquareIterator struct {
	external.Iterator
	db         *sql.DB
	properties string
}

func init() {
	ctx := context.Background()
	err := external.RegisterIterator(ctx, "foursquare", NewFoursquareIterator)
	if err != nil {
		panic(err)
	}
}

func NewFoursquareIterator(ctx context.Context, uri string) (external.Iterator, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	if u.Host != "parquet" {
		return nil, fmt.Errorf("Unsupported data type")
	}

	q := u.Query()
	properties := q.Get("properties")

	engine := "duckdb"
	dsn := ""

	db, err := sql.Open(engine, dsn)

	if err != nil {
		return nil, err
	}

	it := &FoursquareIterator{
		db:         db,
		properties: properties,
	}

	return it, nil
}

func (it *FoursquareIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[external.Record, error] {

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

func (it *FoursquareIterator) iterate(ctx context.Context, uri string) iter.Seq2[external.Record, error] {

	// SELECT fsq_place_id, name, ifnull(latitude, 0.0), ifnull(longitude, 0.0) FROM read_parquet('/usr/local/data/foursquare/parquet/*.parquet')

	logger := slog.Default()
	logger = logger.With("uri", uri)

	return func(yield func(external.Record, error) bool) {

		props := []string{
			"fsq_place_id",
			"name",
			"ifnull(latitude, 0.0)",
			"ifnull(longitude, 0.0)",
		}

		if it.properties == all_properties {

			other_props := []string{
				"ifnull(address, '')",
				"ifnull(locality, '')",
				"ifnull(region, '')",
				"ifnull(postcode, '')",
				"ifnull(admin_region, '')",
				"ifnull(post_town, '')",
				"ifnull(po_box, '')",
				"ifnull(country, '')",
				"ifnull(date_created, '')",
				"ifnull(date_refreshed, '')",
				"ifnull(date_closed, '')",
				"ifnull(tel, '')",
				"ifnull(website, '')",
				"ifnull(facebook_id, 0)",
				"ifnull(instagram, '')",
				"ifnull(twitter, '')",
				"JSON(ifnull(fsq_category_ids, '[]'))",
				"JSON(ifnull(fsq_category_labels, '[]'))",
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
			var geom orb.Geometry

			if it.properties == all_properties {

				var id string
				var name string
				var lat float64
				var lon float64
				var address string
				var locality string
				var region string
				var postcode string
				var admin_region string
				var post_town string
				var po_box string
				var country string
				var date_created string
				var date_refreshed string
				var date_closed string
				var tel string
				var website string
				var facebook_id int64
				var instagram string
				var twitter string
				var str_category_ids string
				var str_category_labels string

				err := rows.Scan(
					&id, &name, &lat, &lon,
					&address, &locality, &region, &postcode, &admin_region, &post_town, &po_box, &country,
					&date_created, &date_refreshed, &date_closed,
					&tel, &website,
					&facebook_id, &instagram, &twitter,
					&str_category_ids, &str_category_labels,
				)

				if err != nil {
					logger.Error("Failed to scan row", "error", err)
					yield(nil, err)
					return
				}

				geom = orb.Point([]float64{lon, lat})

				var category_ids []string
				var category_labels []string

				err = json.Unmarshal([]byte(str_category_ids), &category_ids)

				if err != nil {
					logger.Error("Failed to decode category IDs (%s), %w", str_category_ids, err)
					yield(nil, err)
					return
				}

				err = json.Unmarshal([]byte(str_category_labels), &category_labels)

				if err != nil {
					logger.Error("Failed to decode category labels (%s), %w", str_category_labels, err)
					yield(nil, err)
					return
				}

				props = map[string]any{
					"id":              id,
					"name":            name,
					"address":         address,
					"locality":        locality,
					"region":          region,
					"postcode":        postcode,
					"admin_region":    admin_region,
					"post_town":       post_town,
					"po_box":          po_box,
					"country":         country,
					"date_created":    date_created,
					"date_refreshed":  date_refreshed,
					"date_closed":     date_closed,
					"tel":             tel,
					"website":         website,
					"facebook_id":     facebook_id,
					"instagram":       instagram,
					"twitter":         twitter,
					"category_ids":    category_ids,
					"category_labels": category_labels,
				}

			} else {

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

				geom = orb.Point([]float64{lon, lat})

				props = map[string]any{
					"id":   id,
					"name": name,
				}
			}

			r, err := NewFoursquareRecord(props, geom)

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

func (it *FoursquareIterator) Close() error {
	return it.db.Close()
}
