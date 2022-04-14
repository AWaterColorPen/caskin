package caskin

import (
	_ "embed"

	"github.com/casbin/casbin/v2/model"
)

//go:embed configs/casbin_model.conf
var CasbinModelText string

func CasbinModel() (model.Model, error) {
	return model.NewModelFromString(CasbinModelText)
}
