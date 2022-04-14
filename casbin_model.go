package caskin

import (
	_ "embed"
	"github.com/casbin/casbin/v2/model"
)

//go:embed configs/casbin_model.conf
var casbinModelText string

func CasbinModel() (model.Model, error) {
	return model.NewModelFromString(casbinModelText)
}

func CasbinModelText() string {
	return casbinModelText
}
