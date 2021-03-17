package web_feature

import "github.com/awatercolorpen/caskin"

var (
	DefaultSuperRootName           = "super-root"
	DefaultFeatureRootName         = "feature-root"
	DefaultFeatureRootDescription  = "root node of feature object"
	DefaultFeatureRootGroup        = ""
	DefaultFrontendRootKey         = "frontend-root"
	DefaultFrontendRootType        = FrontendTypeRoot
	DefaultFrontendRootDescription = "root node of frontend object"
	DefaultFrontendRootGroup       = ""
	DefaultBackendRootPath         = "backend-root"
	DefaultBackendRootMethod       = ""
	DefaultBackendRootDescription  = "root node of backend object"
	DefaultBackendRootGroup        = ""
)

// Options configuration for web_feature
type Options struct {
	Domain caskin.Domain `json:"domain"`
}
