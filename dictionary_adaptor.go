package caskin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type DictionaryType string

const (
	FILEDictionary DictionaryType = "FILE"
)

type DictionaryOption struct {
	Type DictionaryType `json:"type"`
	Dsn  string         `json:"dsn"`
}

func NewDictionary(option *DictionaryOption) (IDictionary, error) {
	switch option.Type {
	case FILEDictionary, "":
		return newDictionaryByFile(option)
	default:
		return nil, fmt.Errorf("not supported dictionary type %v", option.Type)
	}
}

type fileDictionary struct {
	Feature  []*Feature  `toml:"feature"`
	Backend  []*Backend  `toml:"backend"`
	Frontend []*Frontend `toml:"frontend"`
	Package  []*Package  `toml:"package"`
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
	return nil, nil
}

func (f *fileDictionary) GetBackendByKey(key string) (*Backend, error) {
	return nil, nil
}

func (f *fileDictionary) GetFrontendByKey(key string) (*Frontend, error) {
	return nil, nil
}

func (f *fileDictionary) isValid() error {
	return nil
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
	return dictionary, dictionary.isValid()
}