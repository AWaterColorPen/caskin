package web_feature

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/awatercolorpen/caskin"
)

type Dump struct {
	Feature         []*caskin.CustomizedDataPair `json:"feature"`
	Frontend        []*caskin.CustomizedDataPair `json:"frontend"`
	Backend         []*caskin.CustomizedDataPair `json:"backend"`
	FeatureTree     Relations                    `json:"feature_tree"`
	FrontendTree    Relations                    `json:"frontend_tree"`
	BackendTree     Relations                    `json:"backend_tree"`
	FeatureRelation Relations                    `json:"feature_relation"`
}

func (d *Dump) ToRelation() *DumpRelation {
	return &DumpRelation{
		FeatureTree:     caskin.SortedInheritanceRelations(d.FeatureTree),
		FrontendTree:    caskin.SortedInheritanceRelations(d.FrontendTree),
		BackendTree:     caskin.SortedInheritanceRelations(d.BackendTree),
		FeatureRelation: caskin.SortedInheritanceRelations(d.FeatureRelation),
	}
}

func (d *Dump) InitFeatureRelationMap(relations Relations) Relations {
	return initFeatureRelationMap(d.FeatureTree, d.FrontendTree, d.BackendTree, relations)
}

func (d *Dump) InitSingleFeatureRelation(relation Relation) Relation {
	return initSingleFeatureRelation(d.FeatureTree, d.FrontendTree, d.BackendTree, relation)
}

func NewDump(feature, frontend, backend []*caskin.CustomizedDataPair, relations Relations) *Dump {
	dump := &Dump{
		Feature:      feature,
		Frontend:     frontend,
		Backend:      backend,
		FeatureTree:  initTreeMapFromPair(feature),
		FrontendTree: initTreeMapFromPair(frontend),
		BackendTree:  initTreeMapFromPair(backend),
	}
	dump.FeatureRelation = dump.InitFeatureRelationMap(relations)
	return dump
}

type DumpRelation struct {
	FeatureTree     Relations `json:"feature_tree"`
	FrontendTree    Relations `json:"frontend_tree"`
	BackendTree     Relations `json:"backend_tree"`
	FeatureRelation Relations `json:"feature_relation"`
}

func (d *DumpRelation) Version() string {
	h := sha256.New()
	b, _ := json.Marshal(d)
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (d *DumpRelation) MergedRelation() Relations {
	m := Relations{}
	for _, relations := range []Relations{d.FeatureTree, d.FrontendTree, d.BackendTree} {
		for k, relation := range relations {
			m[k] = Relation{}
			for _, v := range relation {
				m[k] = append(m[k], v)
			}
		}
	}

	for k, relation := range d.FeatureRelation {
		for _, v := range relation {
			m[k] = append(m[k], v)
		}
	}
	return m
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (d *DumpRelation) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, d)
}

// Value return json value, implement driver.Valuer interface
func (d DumpRelation) Value() (driver.Value, error) {
	bytes, err := json.Marshal(d)
	return string(bytes), err
}
