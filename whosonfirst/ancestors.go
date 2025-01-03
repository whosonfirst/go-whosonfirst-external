package whosonfirst

import (
	"context"
	"log/slog"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/paulmach/orb/geojson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-external"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/hierarchy"
	hierarchy_filter "github.com/whosonfirst/go-whosonfirst-spatial/hierarchy/filter"
	"slices"
	"strconv"
)

type Ancestors struct {
	ParentId    int64
	Country     string
	Hierarchies []map[string]int64
}

type DeriveAncestorsOptions struct {
	SpatialDatabase  database.SpatialDatabase
	Resolver         *hierarchy.PointInPolygonHierarchyResolver
	ParentCache      *ristretto.Cache[string, *Ancestors]
	ResultsCallback  hierarchy_filter.FilterSPRResultsFunc
	PropertiesReader reader.Reader
}

func DeriveAncestors(ctx context.Context, opts *DeriveAncestorsOptions, r external.Record) (*Ancestors, error) {

	logger := slog.Default()
	logger = logger.With("id", r.Id())

	parent_id := int64(-1)
	ancestors := &Ancestors{}

	f := geojson.NewFeature(r.Geometry())

	f.Properties["wof:id"] = r.Id()
	f.Properties["wof:name"] = r.Name()
	f.Properties["wof:placetype"] = r.Placetype()

	body, err := f.MarshalJSON()

	if err != nil {
		logger.Error("Failed to marshal JSON", "error", err)
		return nil, err
	}

	inputs := &filter.SPRInputs{}
	inputs.IsCurrent = []int64{1}

	possible, err := opts.Resolver.PointInPolygon(ctx, inputs, body)

	if err != nil {
		logger.Error("Failed to resolve PIP", "error", err)
		return nil, err
	}

	parent_spr, err := opts.ResultsCallback(ctx, opts.SpatialDatabase, body, possible)

	if err != nil {
		logger.Error("Failed to process results", "error", err)
		return nil, err
	}

	if parent_spr != nil {

		p_id, err := strconv.ParseInt(parent_spr.Id(), 10, 64)

		if err != nil {
			logger.Error("Failed to parse parse parent ID", "id", parent_spr.Id(), "error", err)
			return nil, err
		}

		parent_id = p_id

		k := parent_spr.Id()

		v, exists := opts.ParentCache.Get(k)

		if exists {
			ancestors = v
		} else {

			hierarchies := make([]map[string]int64, 0)
			country := ""

			if p_id >= 0 {

				parent_body, err := wof_reader.LoadBytes(ctx, opts.PropertiesReader, p_id)

				if err != nil {
					logger.Warn("Failed to derive record from properties reader", "id", p_id, "error", err)
				} else {
					hierarchies = properties.Hierarchies(parent_body)
				}

				country = properties.Country(parent_body)
			}

			if country == "" {

				country_ids := make([]int64, 0)

				for _, h := range hierarchies {

					id, ok := h["country_id"]

					if ok && id > -1 && !slices.Contains(country_ids, id) {
						country_ids = append(country_ids, id)
					}
				}

				switch len(country_ids) {
				case 0:
					country = "XY"
				case 1:

					country_id := country_ids[0]
					country_body, err := wof_reader.LoadBytes(ctx, opts.PropertiesReader, country_id)

					if err != nil {
						logger.Warn("Failed to load record for country", "id", country_id, "error", err)
					} else {
						country = properties.Country(country_body)
					}

					if country == "" {
						country = "XY"
					}

				default:
					country = "XZ"
				}
			}

			ancestors.Hierarchies = hierarchies
			ancestors.Country = country

			opts.ParentCache.Set(k, ancestors, 1)
		}
	}

	ancestors.ParentId = parent_id
	return ancestors, nil
}
