package web_feature

import (
	"fmt"
	"github.com/awatercolorpen/caskin"
)

// Frontend
type Frontend struct {
	Key         string       `json:"key"`
	Type        FrontendType `json:"type"`
	Description string       `json:"description"`
	Group       string       `json:"group"`
}

type FrontendType string

const (
	FrontendTypeMenu        FrontendType = "menu"
	FrontendTypeSubFunction FrontendType = "sub_function"
)

func (f *Frontend) GetName() string {
	return fmt.Sprint(f.Key, DefaultSeparator, f.Type)
}

func (f *Frontend) GetObjectType() caskin.ObjectType {
	return ObjectTypeFrontend
}

func frontendFactory() caskin.CustomizedData {
	return &Frontend{}
}
