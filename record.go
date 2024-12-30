package external

import (
	"fmt"

	"github.com/paulmach/orb"
)

type Record interface {
	Id() string
	Name() string
	Placetype() string
	Namespace() string
	Geometry() orb.Geometry
	Properties() map[string]any
}

type ExternalRecord struct {
	Record             `json:",omitempty"`
	ExternalProperties map[string]any `json:"properties"`
	ExternalGeometry   orb.Geometry   `json:"geometry"`
	namespace          string         `json:"namespace"`
	placetype          string         `json:"placetype"`
	id_key             string         `json:"id_key"`
	name_key           string         `json:"name_key"`
}

type NewExternalRecordOptions struct {
	Properties map[string]any
	Geometry   orb.Geometry
	Namespace  string
	Placetype  string
	IdKey      string
	NameKey    string
}

func NewExternalRecord(opts *NewExternalRecordOptions) (Record, error) {

	_, has_id := opts.Properties[opts.IdKey]

	if !has_id {
		return nil, fmt.Errorf("Properties missing %s key", opts.IdKey)
	}

	_, has_name := opts.Properties[opts.NameKey]

	if !has_name {
		return nil, fmt.Errorf("Properties missing %s key", opts.NameKey)
	}

	r := &ExternalRecord{
		ExternalProperties: opts.Properties,
		ExternalGeometry:   opts.Geometry,
		namespace:          opts.Namespace,
		placetype:          opts.Placetype,
		id_key:             opts.IdKey,
		name_key:           opts.NameKey,
	}

	return r, nil
}

func (r *ExternalRecord) Id() string {
	return fmt.Sprintf("%s", r.ExternalProperties[r.id_key])
}

func (r *ExternalRecord) Name() string {
	return fmt.Sprintf("%s", r.ExternalProperties[r.name_key])
}

func (r *ExternalRecord) Placetype() string {
	return r.placetype
}

func (r *ExternalRecord) Namespace() string {
	return r.namespace
}

func (r *ExternalRecord) Geometry() orb.Geometry {
	return r.ExternalGeometry
}

func (r *ExternalRecord) Properties() map[string]any {
	return r.ExternalProperties
}
