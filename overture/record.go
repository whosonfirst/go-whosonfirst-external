package overture

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/whosonfirst/go-whosonfirst-external"
)

const NAMESPACE string = "ovtr"

func NewOvertureRecord(props map[string]any, geom orb.Geometry) (external.Record, error) {

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
		Placetype:  "place",
		IdKey:      "id",
		NameKey:    "name",
	}

	return external.NewExternalRecord(opts)
}
