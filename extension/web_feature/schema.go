package web_feature

import (
	"fmt"

	"github.com/awatercolorpen/caskin"
	"gorm.io/datatypes"
)

const (
	ObjectTypeFeature  caskin.ObjectType = "feature"
	ObjectTypeFrontend caskin.ObjectType = "frontend"
	ObjectTypeBackend  caskin.ObjectType = "backend"
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

func (f *Feature) JSONQuery() []*datatypes.JSONQueryExpression {
	var expression []*datatypes.JSONQueryExpression
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(f.Name, "name"))
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(f.Description, "description"))
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(f.Group, "group"))
	return expression
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
	FrontendTypeNil         FrontendType = ""
	FrontendTypeMenu        FrontendType = "menu"
	FrontendTypeSubFunction FrontendType = "sub_function"
)

func (f *Frontend) GetName() string {
	return fmt.Sprint(f.Key, caskin.DefaultSeparator, f.Type)
}

func (f *Frontend) GetObjectType() caskin.ObjectType {
	return ObjectTypeFrontend
}

func (f *Frontend) JSONQuery() []*datatypes.JSONQueryExpression {
	var expression []*datatypes.JSONQueryExpression
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(f.Key, "key"))
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(f.Type, "type"))
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(f.Description, "description"))
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(f.Group, "group"))
	return expression
}

// Backend
type Backend struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Group       string `json:"group"`
}

func (b *Backend) GetName() string {
	return fmt.Sprint(b.Path, caskin.DefaultSeparator, b.Method)
}

func (b *Backend) GetObjectType() caskin.ObjectType {
	return ObjectTypeBackend
}

func (b *Backend) JSONQuery() []*datatypes.JSONQueryExpression {
	var expression []*datatypes.JSONQueryExpression
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(b.Path, "path"))
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(b.Method, "method"))
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(b.Description, "description"))
	expression = append(expression, datatypes.JSONQuery("customized_data").Equals(b.Group, "group"))
	return expression
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
