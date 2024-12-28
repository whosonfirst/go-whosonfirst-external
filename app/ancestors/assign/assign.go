package assign

// This assumes a PMTiles spatial database described here:
// https://millsfield.sfomuseum.org/blog/2022/12/19/pmtiles-pip/

/*

./bin/reverse-geocode \
    -workers 5 \
    -emitter-uri csv:///usr/local/data/4sq/4sq.csv.bz2 \
    -spatial-database-uri 'pmtiles://?tiles=file:///usr/local/data/pmtiles/&database=whosonfirst-point-in-polygon-z13-20240406&enable-cache=true&pmtiles-cache-size=4096&zoom=13&layer=whosonfirst'

*/

/*

 COPY (SELECT * FROM read_csv('/usr/local/data/overture-wof.csv', AUTO_DETECT=TRUE)) TO '/usr/local/data/overture-wof.parquet' (FORMAT 'PARQUET', CODEC 'ZSTD');

 COPY (SELECT ST_GeomFromText('external:geometry') AS geometry, 'external:id' AS external_id FROM read_csv('/usr/local/data/overture-wof.csv', header=true, delim=',', quote='"', ignore_errors=true, columns={'geometry': 'VARCHAR','external_id':'VARCHAR'})) TO '/usr/local/data/overture-wof2.parquet' (FORMAT 'PARQUET', CODEC 'ZSTD');

*/

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-csvdict/v2"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-external"
	"github.com/whosonfirst/go-whosonfirst-external/iterator"
	"github.com/whosonfirst/go-whosonfirst-external/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/hierarchy"
	hierarchy_filter "github.com/whosonfirst/go-whosonfirst-spatial/hierarchy/filter"
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

	// START OF json wah-wah...

	// https://pkg.go.dev/github.com/paulmach/orb/geojson#pkg-variables
	// https://github.com/json-iterator/go
	//
	// "Even the most widely used json-iterator will severely degrade in generic (no-schema) or big-volume JSON serialization and deserialization."
	// https://github.com/bytedance/sonic/blob/main/INTRODUCTION.md
	//
	// I have not verified that claim either way but since we're not trafficing in "big-volume" JSON files
	// I am just going to see how this (json-iterator) goes for now.

	var c = jsoniter.Config{
		EscapeHTML:              true,
		SortMapKeys:             false,
		MarshalFloatWith6Digits: true,
	}.Froze()

	geojson.CustomJSONMarshaler = c
	geojson.CustomJSONUnmarshaler = c

	// END OF json wah-wah...

	iter, err := iterator.NewIterator(ctx, opts.IteratorURI)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer iter.Close()

	spatial_db, err := database.NewSpatialDatabase(ctx, opts.SpatialDatabaseURI)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer spatial_db.Disconnect(ctx)

	var properties_reader reader.Reader
	properties_reader = spatial_db

	if opts.PropertiesReaderURI != "{spatial-database-uri}" {

		r, err := reader.NewReader(ctx, opts.PropertiesReaderURI)

		if err != nil {
			return fmt.Errorf("%w", err)
		}

		properties_reader = r
	}

	results_cb := hierarchy_filter.FirstButForgivingSPRResultsFunc

	resolver_opts := &hierarchy.PointInPolygonHierarchyResolverOptions{
		Database: spatial_db,
	}

	resolver, err := hierarchy.NewPointInPolygonHierarchyResolver(ctx, resolver_opts)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	mu := new(sync.RWMutex)
	wg := new(sync.WaitGroup)

	throttle := make(chan bool, workers)

	for i := 0; i < workers; i++ {
		throttle <- true
	}

	parent_cache, err := ristretto.NewCache(&ristretto.Config[string, []map[string]int64]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	ancestors_opts := &whosonfirst.DeriveAncestorsOptions{
		SpatialDatabase:  spatial_db,
		Resolver:         resolver,
		ResultsCallback:  results_cb,
		ParentCache:      parent_cache,
		PropertiesReader: properties_reader,
	}

	var csv_wr *csvdict.Writer

	counter := int64(0)
	last_processed := int64(0)
	processed := int64(0)
	timing := int64(0)

	ticker := time.NewTicker(time.Duration(10) * time.Second)
	defer ticker.Stop()

	start := time.Now()

	go func() {

		for {
			select {
			case <-ticker.C:

				p := atomic.LoadInt64(&processed)
				diff := int64(0)

				if last_processed > 0 {
					diff = p - last_processed
				}

				last_processed = p

				slog.Info("Status", "counter", counter, "processed", p, "diff", diff, "avg t2p", float64(timing)/float64(p), "elaspsed", time.Since(start))
			}
		}
	}()

	candidates := []string{
		"microhood_id",
		"neighbourhood_id",
		"macrohood_id",
		"borough_id",
		"locality_id",
		"localadmin_id",
		"county_id",
		"region_id",
		"country_id",
		"continent_id",
		"empire_id",
	}

	process_record := func(ctx context.Context, r *external.Record) error {

		t1 := time.Now()

		defer func() {

			t2 := time.Since(t1)
			atomic.AddInt64(&timing, t2.Milliseconds())
			atomic.AddInt64(&processed, 1)
		}()

		a, err := whosonfirst.DeriveAncestors(ctx, ancestors_opts, r)

		if err != nil {
			return err
		}

		geom := ""

		if r.Geometry != nil {
			geom = wkt.MarshalString(r.Geometry)
		}

		count_hiers := len(a.Hierarchies)
		csv_rows := make([]map[string]string, count_hiers)

		switch count_hiers {
		case 0:

			out := map[string]string{
				"external:id":           r.Id,
				"external:geometry":     geom,
				"wof:hierarchies_count": "0",
				"wof:hierarchies_idx":   "0",
				"wof:parent_id":         strconv.FormatInt(a.ParentId, 10),
			}

			for _, label := range candidates {
				wof_label := fmt.Sprintf("wof:%s", label)
				wof_value := ""
				out[wof_label] = wof_value
			}

			csv_rows = []map[string]string{
				out,
			}

		default:

			for i, h := range a.Hierarchies {

				out := map[string]string{
					"external:id":           r.Id,
					"external:geometry":     geom,
					"wof:hierarchies_count": strconv.Itoa(count_hiers),
					"wof:hierarchies_idx":   strconv.Itoa(i),
					"wof:parent_id":         strconv.FormatInt(a.ParentId, 10),
				}

				for _, label := range candidates {

					wof_label := fmt.Sprintf("wof:%s", label)
					wof_value := ""

					v, exists := h[label]

					if exists {
						wof_value = strconv.FormatInt(v, 10)
					}

					out[wof_label] = wof_value
				}

				csv_rows[i] = out
			}

		}

		mu.Lock()
		defer mu.Unlock()

		if csv_wr == nil {

			wr, err := csvdict.NewWriter(os.Stdout)

			if err != nil {
				slog.Error("Failed to create CSV writer", "error", err)
				return err
			}

			csv_wr = wr
		}

		for _, out := range csv_rows {
			csv_wr.WriteRow(out)
		}

		csv_wr.Flush()
		return nil
	}

	for r, err := range iter.Iterate(ctx, opts.IteratorSources...) {

		counter += 1

		if err != nil {
			slog.Error("Failed to yield place", "error", err)
			continue
		}

		if opts.StartAfter > 0 && opts.StartAfter > counter {
			slog.Debug("Start after throttle", "after", opts.StartAfter, "count", counter)
			continue
		}

		<-throttle

		wg.Add(1)

		go func(r *external.Record) {

			defer func() {
				throttle <- true
				wg.Done()
			}()

			err = process_record(ctx, r)

			if err != nil {
				slog.Error("Failed to process place", "error", err)
			}

		}(r)
	}

	wg.Wait()
	return nil
}
