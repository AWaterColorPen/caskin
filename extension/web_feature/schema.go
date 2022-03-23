package web_feature

import (
	"fmt"
	"time"

	"github.com/awatercolorpen/caskin"
	"github.com/patrickmn/go-cache"
)

var (
	DefaultSuperRootName = "github.com/awatercolorpen/caskin/web_feature"

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

// Feature
type Feature struct {
	Name        string `toml:"name"        json:"name"`
	Description string `toml:"description" json:"description"`
	Group       string `toml:"group"       json:"group"`
}

func (f *Feature) GetName() string {
	return f.Name
}

// Frontend
type Frontend struct {
	Key         string       `toml:"key"         json:"key"`
	Type        FrontendType `toml:"type"        json:"type"`
	Description string       `toml:"description" json:"description"`
	Group       string       `toml:"group"       json:"group"`
}

type FrontendType string

const (
	FrontendTypeMenu        FrontendType = "menu"
	FrontendTypeSubFunction FrontendType = "sub_function"
)

func (f *Frontend) GetName() string {
	return fmt.Sprint(f.Key, caskin.DefaultSeparator, f.Type)
}

// Backend
type Backend struct {
	Path        string `toml:"path"        json:"path"`
	Method      string `toml:"method"      json:"method"`
	Description string `toml:"description" json:"description"`
	Group       string `toml:"group"       json:"group"`
}

func (b *Backend) GetName() string {
	return fmt.Sprint(b.Path, caskin.DefaultSeparator, b.Method)
}
