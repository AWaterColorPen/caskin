package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

const (
	ObjectTypeFeature  caskin.ObjectType = "feature"
	ObjectTypeFrontend caskin.ObjectType = "frontend"
	ObjectTypeBackend  caskin.ObjectType = "backend"
)

var (
	DefaultSeparator = "_"
)

// Feature
type Feature struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Group       string `json:"group"`
}

func (f *Feature) GetName() string {
	return f.Name
}

func (f *Feature) GetObjectType() caskin.ObjectType {
	return ObjectTypeFeature
}

func featureFactory() caskin.ObjectCustomizedData {
	return &Feature{}
}
