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
	"strconv"
)

type Ancestors struct {
	ParentId    int64
	Hierarchies []map[string]int64
}

func (a *Ancestors) MarshalHierarchies() string {
	return ""

	/*


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

			str_hier := make([]string, len(hierarchies))

			for i, h := range hierarchies {

				// colon-separated list
				hier_csv := make([]string, len(candidates))

				for j, k := range candidates {

					id, exists := h[k]
					v := ""

					if exists {
						v = strconv.FormatInt(id, 10)
					}

					hier_csv[j] = v
				}

				str_hier[i] = strings.Join(hier_csv, ":")
			}

			str_hierarchies = strings.Join(str_hier, ",")
		}

	*/

}

type DeriveAncestorsOptions struct {
	SpatialDatabase  database.SpatialDatabase
	Resolver         *hierarchy.PointInPolygonHierarchyResolver
	ParentCache      *ristretto.Cache[string, []map[string]int64]
	ResultsCallback  hierarchy_filter.FilterSPRResultsFunc
	PropertiesReader reader.Reader
}

func DeriveAncestors(ctx context.Context, opts *DeriveAncestorsOptions, r external.Record) (*Ancestors, error) {

	logger := slog.Default()
	logger = logger.With("id", r.Id())

	parent_id := int64(-1)
	hierarchies := make([]map[string]int64, 0)

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
			hierarchies = v
		} else {

			// belongs_to = parent_spr.BelongsTo()

			parent_body, err := wof_reader.LoadBytes(ctx, opts.PropertiesReader, p_id)

			if err != nil {
				logger.Warn("Failed to derive record from properties reader", "id", p_id, "error", err)
			} else {
				hierarchies = properties.Hierarchies(parent_body)
			}

			opts.ParentCache.Set(k, hierarchies, 1)
		}
	}

	a := &Ancestors{
		ParentId:    parent_id,
		Hierarchies: hierarchies,
	}

	return a, nil
}
