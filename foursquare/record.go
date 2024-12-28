package foursquare

import (
	"fmt"
	
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/whosonfirst/go-whosonfirst-external"
)

type FoursquareRecord struct {
	external.RecordWIP `json:",omitempty"`
	body string
}

func (r *FoursquareRecord) Id() string {
	return "4sq:id="
}

func (r *FoursquareRecord) Name() string {
	return ""
}

func (r *FoursquareRecord) Placetype() string {
	return "venue"
}

func (r *FoursquareRecord) Geometry() orb.Geometry {
	return nil
}

func (r *FoursquareRecord) Body() string {
	return ""
}

func (r *FoursquareRecord) AsFeature() (*geojson.Feature, error) {
	return nil, fmt.Errorf("Not implemented")
}

