package caskin

import (
	"bytes"
	"net/http"

	"github.com/casbin/casbin/v2/model"
)

var (
	casbinModelText1 = casbinModelText2 + ` || g(r.sub, "superadmin", "superdomain")`

	casbinModelText2 = `[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _
g2 = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && g2(r.obj, p.obj, r.dom) && r.dom == p.dom && r.act == p.act`
)

type casbinModel struct {
	m    model.Model
	text string
}

func CasbinModel(isDisableSuperadmin bool) (model.Model, error) {
	switch isDisableSuperadmin {
	case true:
		return CasbinModelNoSuperadmin()
	default:
		return CasbinModelSuperadmin()
	}
}

func CasbinModelText(isDisableSuperadmin bool) (string, error) {
	switch isDisableSuperadmin {
	case true:
		return CasbinModelTextNoSuperadmin()
	default:
		return CasbinModelTextSuperadmin()
	}
}

func CasbinModelSuperadmin() (model.Model, error) {
	return model.NewModelFromString(casbinModelText1)
}

func CasbinModelTextSuperadmin() (string, error) {
	return casbinModelText1, nil
}

func CasbinModelNoSuperadmin() (model.Model, error) {
	return model.NewModelFromString(casbinModelText2)
}

func CasbinModelTextNoSuperadmin() (string, error) {
	return casbinModelText2, nil
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
