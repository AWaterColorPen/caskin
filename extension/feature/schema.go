package feature

import (
	"time"

	"github.com/awatercolorpen/caskin"
	"github.com/patrickmn/go-cache"
)

var (
	DefaultSuperRootName = "github.com/awatercolorpen/caskin/feature"

	// local cache
	LocalCacheDefaultExpiration = 2 * time.Minute
	LocalCacheCleanupInterval   = 5 * time.Minute
	LocalCache                  = cache.New(LocalCacheDefaultExpiration, LocalCacheCleanupInterval)
)

const (
	ObjectTypeFeature  caskin.ObjectType = "feature"
	ObjectTypeFrontend caskin.ObjectType = "frontend"
	ObjectTypeBackend  caskin.ObjectType = "backend"
)
