package web_feature

import (
	"fmt"

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

// Frontend
type Frontend struct {
	Key         string       `json:"key"`
	Type        FrontendType `json:"type"`
	Description string       `json:"description"`
	Group       string       `json:"group"`
}

type FrontendType string

const (
	FrontendTypeRoot        FrontendType = "root"
	FrontendTypeMenu        FrontendType = "menu"
	FrontendTypeSubFunction FrontendType = "sub_function"
)

func (f *Frontend) GetName() string {
	return fmt.Sprint(f.Key, DefaultSeparator, f.Type)
}

func (f *Frontend) GetObjectType() caskin.ObjectType {
	return ObjectTypeFrontend
}

// Backend
type Backend struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Group       string `json:"group"`
}

func (b *Backend) GetName() string {
	return fmt.Sprint(b.Path, DefaultSeparator, b.Method)
}

func (b *Backend) GetObjectType() caskin.ObjectType {
	return ObjectTypeBackend
}

func FeatureFactory() caskin.CustomizedData {
	return &Feature{}
}

func FrontendFactory() caskin.CustomizedData {
	return &Frontend{}
}

func BackendFactory() caskin.CustomizedData {
	return &Backend{}
}
