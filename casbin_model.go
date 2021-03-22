package caskin

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/casbin/casbin/v2/model"
)

var (
	casbinModel1 *casbinModel
	casbinModel2 *casbinModel
	err1         error
	err2         error
	once1        sync.Once
	once2        sync.Once
)

type casbinModel struct {
	m    model.Model
	text string
}

func CasbinModel(options *Options) (model.Model, error) {
	switch options.IsDisableSuperAdmin() {
	case true:
		return CasbinModelNoSuperadmin()
	default:
		return CasbinModelSuperadmin()
	}
}

func CasbinModelText(options *Options) (string, error) {
	switch options.IsDisableSuperAdmin() {
	case true:
		return CasbinModelTextNoSuperadmin()
	default:
		return CasbinModelTextSuperadmin()
	}
}

func CasbinModelSuperadmin() (model.Model, error) {
	once1.Do(func() {
		casbinModel1, err1 = getCasbinModelFromUrl(ModelConfPathSuperadmin)
	})
	if err1 != nil {
		return nil, err1
	}
	return casbinModel1.m, nil
}

func CasbinModelTextSuperadmin() (string, error) {
	once1.Do(func() {
		casbinModel1, err1 = getCasbinModelFromUrl(ModelConfPathSuperadmin)
	})
	if err1 != nil {
		return "", err1
	}
	return casbinModel1.text, nil
}

func CasbinModelNoSuperadmin() (model.Model, error) {
	once2.Do(func() {
		casbinModel2, err2 = getCasbinModelFromUrl(ModelConfPathNoSuperadmin)
	})
	if err2 != nil {
		return nil, err2
	}
	return casbinModel2.m, nil
}

func CasbinModelTextNoSuperadmin() (string, error) {
	once2.Do(func() {
		casbinModel2, err2 = getCasbinModelFromUrl(ModelConfPathNoSuperadmin)
	})
	if err2 != nil {
		return "", err2
	}
	return casbinModel2.text, nil
}

func getCasbinModelTextFromUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getCasbinModelFromUrl(url string) (*casbinModel, error) {
	text, err := getCasbinModelTextFromUrl(url)
	if err != nil {
		return nil, err
	}
	m, err := model.NewModelFromString(text)
	if err != nil {
		return nil, err
	}
	return &casbinModel{m, text}, nil
}
