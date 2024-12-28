package overture

import (
	"fmt"
	
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/whosonfirst/go-whosonfirst-external"
)

type OvertureRecord struct {
	external.RecordWIP `json:",omitempty"`
	body string
}

func (r *OvertureRecord) Id() string {
	return "4sq:id="
}

func (r *OvertureRecord) Name() string {
	return ""
}

func (r *OvertureRecord) Placetype() string {
	return "venue"
}

func (r *OvertureRecord) Geometry() orb.Geometry {
	return nil
}

func (r *OvertureRecord) Body() string {
	return ""
}

func (r *OvertureRecord) AsFeature() (*geojson.Feature, error) {
	return nil, fmt.Errorf("Not implemented")
}

