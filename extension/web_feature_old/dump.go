package web_feature_old

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/awatercolorpen/caskin"
)

type Dump struct {
	FeatureTree     Relations `json:"feature_tree"`
	FrontendTree    Relations `json:"frontend_tree"`
	BackendTree     Relations `json:"backend_tree"`
	FeatureRelation Relations `json:"feature_relation"`
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

func NewDump(feature, frontend, relations Relations) *Dump {
	dump := &Dump{}
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
	return caskin.MergedInheritanceRelations(d.FeatureTree, d.FrontendTree, d.BackendTree, d.FeatureRelation)
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

type DumpFileStruct struct {
	Feature                     []*Feature  `json:"feature"`
	Frontend                    []*Frontend `json:"frontend"`
	Backend                     []*Backend  `json:"backend"`
	VirtualIndexFeatureTree     Relations   `json:"virtual_index_feature_tree"`
	VirtualIndexFrontendTree    Relations   `json:"virtual_index_frontend_tree"`
	VirtualIndexBackendTree     Relations   `json:"virtual_index_backend_tree"`
	VirtualIndexFeatureRelation Relations   `json:"virtual_index_feature_relation"`
}

func (d *DumpFileStruct) IsValid() error {
	if len(d.Feature) != len(d.VirtualIndexFeatureTree) {
		return caskin.ErrInCompatible
	}
	if len(d.Frontend) != len(d.VirtualIndexFrontendTree) {
		return caskin.ErrInCompatible
	}
	if len(d.Backend) != len(d.VirtualIndexBackendTree) {
		return caskin.ErrInCompatible
	}
	if len(d.VirtualIndexFeatureRelation) != len(d.VirtualIndexFeatureTree) {
		return caskin.ErrInCompatible
	}
	return nil
}

func treeToVirtualIndexTree(tree, viTree Relations, index map[uint64]uint64) {
	for k, node := range tree {
		var vi []uint64
		for _, v := range node {
			vi = append(vi, index[v])
		}
		viTree[index[k]] = vi
	}
}
