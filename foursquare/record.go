package foursquare

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/whosonfirst/go-whosonfirst-external"
)

const NAMESPACE string = "4sq"

func NewFoursquareRecord(props map[string]any, geom orb.Geometry) (external.Record, error) {

	_, has_id := props["id"]

	if !has_id {
		return nil, fmt.Errorf("Properties missing id")
	}

	_, has_name := props["name"]

	if !has_name {
		return nil, fmt.Errorf("Properties missing name")
	}

	opts := &external.NewExternalRecordOptions{

		Properties: props,
		Geometry:   geom,
		Namespace:  NAMESPACE,
		Placetype:  "venue",
		IdKey:      "id",
		NameKey:    "name",
	}

	return external.NewExternalRecord(opts)
}
