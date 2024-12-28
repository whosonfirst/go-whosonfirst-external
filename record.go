package external

import (
	"github.com/paulmach/orb"
)

type Record struct {
	Id        string       `json:"id"`
	Name      string       `json:"name"`
	Placetype string       `json:"placetype"`
	Geometry  orb.Geometry `json:"geometry"`
}
