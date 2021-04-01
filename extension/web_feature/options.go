package web_feature

import (
	"time"

	"github.com/awatercolorpen/caskin"
	"github.com/patrickmn/go-cache"
)

var (
	DefaultSuperRootName              = "github.com/awatercolorpen/caskin/web_feature"
	DefaultFeatureRootName            = "feature-root"
	DefaultFeatureRootDescription     = "root node of feature object"
	DefaultFeatureRootGroup           = ""
	DefaultFrontendRootKey            = "frontend-root"
	DefaultFrontendRootType           = FrontendTypeNil
	DefaultFrontendRootDescription    = "root node of frontend object"
	DefaultFrontendRootGroup          = ""
	DefaultBackendRootPath            = "backend-root"
	DefaultBackendRootMethod          = ""
	DefaultBackendRootDescription     = "root node of backend object"
	DefaultBackendRootGroup           = ""
	DefaultWebFeatureVersionTableName = "caskin_web_feature_versions"

	// local cache
	LocalCacheDefaultExpiration = 2 * time.Minute
	LocalCacheCleanupInterval   = 5 * time.Minute
	LocalCache                  = cache.New(LocalCacheDefaultExpiration, LocalCacheCleanupInterval)
)

// Options configuration for web_feature
type Options struct {
	Domain       caskin.Domain `json:"domain"`
}
