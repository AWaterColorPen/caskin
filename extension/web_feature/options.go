package web_feature

import "github.com/awatercolorpen/caskin"

var (
	DefaultFeatureRootName            = "feature-root"
	DefaultFeatureRootDescriptionName = "root node of feature object"
	DefaultFeatureRootGroupName       = ""
)

// Options configuration for web_feature
type Options struct {
	Domain caskin.Domain `json:"domain"`
}
