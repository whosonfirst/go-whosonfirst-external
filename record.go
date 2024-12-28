package external

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"	
)

type RecordWIP interface {
	Id() string
	Name() string
	Placetype() string
	Geometry() orb.Geometry
	Body() string
	AsFeature() (*geojson.Feature, error)
}

type Record struct {
	Id        string       `json:"id"`
	Name      string       `json:"name"`
	Placetype string       `json:"placetype"`
	Geometry  orb.Geometry `json:"geometry"`
}
