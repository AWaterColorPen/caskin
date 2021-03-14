package web_feature

import (
	"fmt"

	"github.com/awatercolorpen/caskin"
)

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

func backendFactory() caskin.CustomizedData {
	return &Backend{}
}