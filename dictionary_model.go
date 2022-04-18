package caskin

import (
	"encoding/json"
)

var (
	DefaultFeatureRootName = "github.com/awatercolorpen/caskin::feature"
)

// Feature it is a package of Backend and Frontend
type Feature struct {
	Name        string `json:"name"        toml:"name"`
	Description string `json:"description" toml:"description"`
	Group       string `json:"group"       toml:"group"`
}

func (f *Feature) Key() string {
	return f.Name
}

func (f *Feature) ToObject() Object {
	return &NamedObject{Name: f.Key()}
}

// Backend it is for backend API
type Backend struct {
	Path        string `json:"path"        toml:"path"`
	Method      string `json:"method"      toml:"method"`
	Description string `json:"description" toml:"description"`
	Group       string `json:"group"       toml:"group"`
}

func (b *Backend) Key() string {
	s := []string{b.Path, b.Method}
	bb, _ := json.Marshal(s)
	return string(bb)
}

func (b *Backend) ToObject() Object {
	return &NamedObject{Name: b.Key()}
}

// Frontend it is for frontend web component
type Frontend struct {
	Name        string `json:"name"        toml:"name"`
	Type        string `json:"type"        toml:"type"`
	Description string `json:"description" toml:"description"`
	Group       string `json:"group"       toml:"group"`
}

func (f *Frontend) Key() string {
	s := []string{f.Name, f.Type}
	b, _ := json.Marshal(s)
	return string(b)
}

func (f *Frontend) ToObject() Object {
	return &NamedObject{Name: f.Key()}
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

type CreatorObject struct {
	Name        string `json:"name"        toml:"name"`
	Type        string `json:"type"        toml:"type"`
	Description string `json:"description" toml:"description"`
}

func (c *CreatorObject) ToObject() Object {
	o := DefaultFactory().NewObject()
	b, _ := json.Marshal(c)
	_ = json.Unmarshal(b, o)
	return o
}

type CreatorRole struct {
	Name        string `json:"name"        toml:"name"`
	Description string `json:"description" toml:"description"`
}

func (c *CreatorRole) ToRole() Role {
	r := DefaultFactory().NewRole()
	b, _ := json.Marshal(c)
	_ = json.Unmarshal(b, r)
	return r
}

type CreatorPolicy struct {
	Object string   `json:"object" toml:"object"`
	Role   string   `json:"role"   toml:"role"`
	Action []string `json:"action" toml:"action"`
}
