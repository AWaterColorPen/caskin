package caskin

import (
	_ "embed"
	"github.com/casbin/casbin/v2/model"
)

// TODO use builtin text
// var (
// 	casbinModelText = `[request_definition]
// r = sub, dom, obj, act
//
// [policy_definition]
// p = sub, dom, obj, act
//
// [role_definition]
// g = _, _, _
// g2 = _, _, _
//
// [policy_effect]
// e = some(where (p.eft == allow))
//
// [matchers]
// m = g(r.sub, p.sub, r.dom) && g2(r.obj, p.obj, r.dom) && r.dom == p.dom && r.act == p.act || g(r.sub, "superadmin", "superdomain")`
// )

//go:embed configs/casbin_model.conf
var casbinModelText string

func CasbinModel() (model.Model, error) {
	return model.NewModelFromString(casbinModelText)
}

func CasbinModelText() string {
	return casbinModelText
}
