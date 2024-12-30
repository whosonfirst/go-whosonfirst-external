package external

import (
	"fmt"

	"github.com/paulmach/orb/geojson"
)

func AsGeoJSONFeature(r Record) (*geojson.Feature, error) {

	geom := r.Geometry()

	if geom == nil {
		return nil, fmt.Errorf("Record is missing geometry")
	}

	props := make(map[string]any)
	ns := r.Namespace()

	for k, v := range r.Properties() {
		k = fmt.Sprintf("%s:%s", ns, k)
		props[k] = v
	}

	f := geojson.NewFeature(geom)
	f.Properties = props

	return f, nil
}
