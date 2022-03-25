package feature

import (
	"fmt"
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

// Feature
type Feature struct {
	Name        string `json:"name"        toml:"name"`
	Description string `json:"description" toml:"description"`
	Group       string `json:"group"       toml:"group"`
}

func (f *Feature) GetKey() string {
	return f.Name
}

func (f *Feature) ToObject() caskin.Object {
	return nil
}

// Backend
type Backend struct {
	Path        string `json:"path"        toml:"path"`
	Method      string `json:"method"      toml:"method"`
	Description string `json:"description" toml:"description"`
	Group       string `json:"group"       toml:"group"`
}

func (b *Backend) GetKey() string {
	return fmt.Sprint(b.Path, caskin.DefaultSeparator, b.Method)
}

func (b *Backend) ToObject() caskin.Object {
	return nil
}

// Frontend
type Frontend struct {
	Name        string       `json:"name"        toml:"name"`
	Type        FrontendType `json:"type"        toml:"type"`
	Description string       `json:"description" toml:"description"`
	Group       string       `json:"group"       toml:"group"`
}

type FrontendType string

const (
	FrontendTypeMenu        FrontendType = "menu"
	FrontendTypeSubFunction FrontendType = "sub_function"
)

func (f *Frontend) GetKey() string {
	return fmt.Sprint(f.Name, caskin.DefaultSeparator, f.Type)
}

func (f *Frontend) ToObject() caskin.Object {
	return nil
}

type Package struct {
	Key      string   `toml:"key"`
	Feature  []string `toml:"feature"`
	Backend  []string `toml:"backend"`
	Frontend []string `toml:"frontend"`
}

func (p *Package) GetName() string {
	return ""
}
