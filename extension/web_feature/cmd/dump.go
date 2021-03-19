package cmd

import (
    "fmt"

    "github.com/awatercolorpen/caskin/extension/web_feature"
)

var (
    ErrInvalidExportData = fmt.Errorf("invalid export data")
)

type Dump struct {
	Feature                     []*web_feature.Feature  `json:"feature"`
	Frontend                    []*web_feature.Frontend `json:"frontend"`
	Backend                     []*web_feature.Backend  `json:"backend"`
	VirtualIndexFeatureTree     map[uint64][]uint64     `json:"virtual_index_feature_tree"`
	VirtualIndexFrontendTree    map[uint64][]uint64     `json:"virtual_index_frontend_tree"`
	VirtualIndexBackendTree     map[uint64][]uint64     `json:"virtual_index_backend_tree"`
	VirtualIndexFeatureRelation map[uint64][]uint64     `json:"virtual_index_feature_relation"`
}

func (d *Dump) IsValid() error {
    if len(d.Feature) != len(d.VirtualIndexFeatureTree) {
        return ErrInvalidExportData
    }
    if len(d.Frontend) != len(d.VirtualIndexFrontendTree) {
        return ErrInvalidExportData
    }
    if len(d.Backend) != len(d.VirtualIndexBackendTree) {
        return ErrInvalidExportData
    }
    if len(d.VirtualIndexFeatureRelation) != len(d.VirtualIndexBackendTree) {
        return ErrInvalidExportData
    }
    return nil
}

func ExportFromWebFeature(w *web_feature.WebFeature) (*Dump, error) {
    // dump := &Dump{}
    // executor := w.GetExecutor(nil)
    // pair1, err := executor.GetFeature()
    // if err != nil {
    //     return nil, err
    // }
    // pair2, err := executor.GetFrontend()
    // if err != nil {
    //     return nil, err
    // }
    // pair3, err := executor.GetBackend()
    // if err != nil {
    //     return nil, err
    // }
    // executor.
    return nil, nil
}

func ExportFromFile(w *web_feature.WebFeature) (*Dump, error) {
    return nil, nil
}

func ImportToWebFeature(d *Dump, w *web_feature.WebFeature) error {
    return nil
}