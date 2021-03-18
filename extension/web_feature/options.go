package web_feature

import (
	"github.com/awatercolorpen/caskin"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	DefaultSuperRootName           = "super-root"
	DefaultFeatureRootName         = "feature-root"
	DefaultFeatureRootDescription  = "root node of feature object"
	DefaultFeatureRootGroup        = ""
	DefaultFrontendRootKey         = "frontend-root"
	DefaultFrontendRootType        = FrontendTypeNil
	DefaultFrontendRootDescription = "root node of frontend object"
	DefaultFrontendRootGroup       = ""
	DefaultBackendRootPath         = "backend-root"
	DefaultBackendRootMethod       = ""
	DefaultBackendRootDescription  = "root node of backend object"
	DefaultBackendRootGroup        = ""

	DefaultSeparator = caskin.DefaultSeparator
)

var (
	// local cache for
	LocalCache = cache.New(2*time.Minute, 5*time.Minute)
)


// Options configuration for web_feature
type Options struct {
	DisableCache bool          `json:"disable_cache"`
	Domain       caskin.Domain `json:"domain"`
}
