package feature

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Dictionary interface {
	GetFeature() ([]*Feature, error)
	GetBackend() ([]*Backend, error)
	GetFrontend() ([]*Frontend, error)
}

type DictionaryType string

const (
	FILEDictionary DictionaryType = "FILE"
)

type DictionaryOption struct {
	Type DictionaryType `json:"type"`
	Dsn  string         `json:"dsn"`
}

func NewDictionary(option *DictionaryOption) (Dictionary, error) {
	switch option.Type {
	case FILEDictionary, "":
		return newDictionaryByFile(option)
	default:
		return nil, fmt.Errorf("not supported dictionary type %v", option.Type)
	}
}

type fileDictionary struct {
	Feature  []*Feature  `toml:"Feature"`
	Backend  []*Backend  `toml:"backend"`
	Frontend []*Frontend `toml:"frontend"`
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

func (f *fileDictionary) isValid() error {
	return nil
}

func newDictionaryByFile(option *DictionaryOption) (*fileDictionary, error) {
	b, err := ioutil.ReadFile(option.Dsn)
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
