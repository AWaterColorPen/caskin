package web_feature_old

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"

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

func (d *DumpFileStruct) ImportFromDump(dump *Dump) error {
	d.VirtualIndexFeatureTree = Relations{}
	d.VirtualIndexFrontendTree = Relations{}
	d.VirtualIndexBackendTree = Relations{}
	d.VirtualIndexFeatureRelation = Relations{}

	indexFe, indexFr, indexBa := map[uint64]uint64{}, map[uint64]uint64{}, map[uint64]uint64{}

	for i, v := range dump.Feature {
		d.Feature = append(d.Feature, v.ObjectCustomizedData.(*Feature))
		indexFe[v.Object.GetID()] = uint64(i)
	}
	for i, v := range dump.Frontend {
		d.Frontend = append(d.Frontend, v.ObjectCustomizedData.(*Frontend))
		indexFr[v.Object.GetID()] = uint64(i)
	}
	for i, v := range dump.Backend {
		d.Backend = append(d.Backend, v.ObjectCustomizedData.(*Backend))
		indexBa[v.Object.GetID()] = uint64(i)
	}

	treeToVirtualIndexTree(dump.FeatureTree, d.VirtualIndexFeatureTree, indexFe)
	treeToVirtualIndexTree(dump.FrontendTree, d.VirtualIndexFrontendTree, indexFr)
	treeToVirtualIndexTree(dump.BackendTree, d.VirtualIndexBackendTree, indexBa)

	for k, node := range dump.FeatureRelation {
		var tree []uint64
		for _, v := range node {
			if u, ok := indexFr[v]; ok {
				tree = append(tree, u)
			}
			if u, ok := indexBa[v]; ok {
				tree = append(tree, u)
			}
		}
		d.VirtualIndexFeatureRelation[indexFe[k]] = tree
	}
	return nil
}

func (d *DumpFileStruct) ImportFromFile(name string) error {
	// go 1.16
	// b, err := os.ReadFile(name)
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, d)
}

func (d *DumpFileStruct) ExportToWebFeature(w *WebFeature) error {
	var executor *Executor
	indexFe, indexFr, indexBa := map[uint64]uint64{}, map[uint64]uint64{}, map[uint64]uint64{}
	for i, v := range d.Feature {
		o := executor.objectFactory()
		if err := executor.CreateFeature(v, o); err != nil {
			return err
		}
		indexFe[uint64(i)] = o.GetID()
	}
	for i, v := range d.Frontend {
		o := executor.objectFactory()
		if err := executor.CreateFrontend(v, o); err != nil {
			return err
		}
		indexFr[uint64(i)] = o.GetID()
	}
	for i, v := range d.Backend {
		o := executor.objectFactory()
		if err := executor.CreateBackend(v, o); err != nil {
			return err
		}
		indexBa[uint64(i)] = o.GetID()
	}

	return nil
}

func (d *DumpFileStruct) ExportToFile(name string) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, b, 0644)
	// go 1.16
	// return os.WriteFile(name, b, 0644)
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
