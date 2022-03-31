package caskin

import (
	"fmt"
)

// Feature it is a package of Backend and Frontend
type Feature struct {
	Name        string `json:"name"        toml:"name"`
	Description string `json:"description" toml:"description"`
	Group       string `json:"group"       toml:"group"`
}

func (f *Feature) GetKey() string {
	return f.Name
}

func (f *Feature) ToObject() Object {
	return nil
}

// Backend it is for backend API
type Backend struct {
	Path        string `json:"path"        toml:"path"`
	Method      string `json:"method"      toml:"method"`
	Description string `json:"description" toml:"description"`
	Group       string `json:"group"       toml:"group"`
}

func (b *Backend) GetKey() string {
	return fmt.Sprint(b.Path, DefaultSeparator, b.Method)
}

func (b *Backend) ToObject() Object {
	return nil
}

// Frontend it is for frontend web component
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
	return fmt.Sprint(f.Name, DefaultSeparator, f.Type)
}

func (f *Frontend) ToObject() Object {
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

type CreatorObject struct {
	Name        string `json:"name"        toml:"name"`
	Type        string `json:"type"        toml:"type"`
	Description string `json:"description" toml:"description"`
}

type CreatorRole struct {
	Name        string `json:"name"        toml:"name"`
	Description string `json:"description" toml:"description"`
}

type CreatorPolicy struct {
	Object string `json:"object" toml:"object"`
	Role   string `json:"role"   toml:"role"`
	Action string `json:"action" toml:"action"`
}
