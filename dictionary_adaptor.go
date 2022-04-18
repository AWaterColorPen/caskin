package caskin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type DictionaryOption struct {
	Dsn  string `json:"dsn"`
	Type string `json:"type"`
}

func NewDictionary(option *DictionaryOption) (IDictionary, error) {
	switch option.Type {
	case "FILE", "":
		return newDictionaryByFile(option)
	default:
		return nil, fmt.Errorf("not supported dictionary type %v", option.Type)
	}
}

type fileDictionary struct {
	Feature       []*Feature       `toml:"feature"`
	Backend       []*Backend       `toml:"backend"`
	Frontend      []*Frontend      `toml:"frontend"`
	Package       []*Package       `toml:"package"`
	CreatorObject []*CreatorObject `toml:"creator_object"`
	CreatorRole   []*CreatorRole   `toml:"creator_role"`
	CreatorPolicy []*CreatorPolicy `toml:"creator_policy"`
}

func (f *fileDictionary) GetFeature() ([]*Feature, error) {
	return f.Feature, nil
}

func (f *fileDictionary) GetBackend() ([]*Backend, error) {
	return f.Backend, nil
}

func (f *fileDictionary) GetFrontend() ([]*Frontend, error) {
	return f.Frontend, nil
}

func (f *fileDictionary) GetFeatureByKey(key string) (*Feature, error) {
	for _, v := range f.Feature {
		if v.Key() == key {
			return v, nil
		}
	}
	return nil, nil
}

func (f *fileDictionary) GetBackendByKey(key string) (*Backend, error) {
	for _, v := range f.Backend {
		if v.Key() == key {
			return v, nil
		}
	}
	return nil, nil
}

func (f *fileDictionary) GetFrontendByKey(key string) (*Frontend, error) {
	for _, v := range f.Frontend {
		if v.Key() == key {
			return v, nil
		}
	}
	return nil, nil
}

func (f *fileDictionary) GetPackage() ([]*Package, error) {
	return f.Package, nil
}

func (f *fileDictionary) GetCreatorObject() ([]*CreatorObject, error) {
	return f.CreatorObject, nil
}

func (f *fileDictionary) GetCreatorRole() ([]*CreatorRole, error) {
	return f.CreatorRole, nil
}

func (f *fileDictionary) GetCreatorPolicy() ([]*CreatorPolicy, error) {
	return f.CreatorPolicy, nil
}

func newDictionaryByFile(option *DictionaryOption) (*fileDictionary, error) {
	b, err := os.ReadFile(option.Dsn)
	if err != nil {
		return nil, err
	}

	dictionary := &fileDictionary{}
	extension := filepath.Ext(option.Dsn)
	switch extension {
	case ".toml":
		if err = toml.Unmarshal(b, dictionary); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("not supported extension %v", extension)
	}
	return dictionary, nil
}
