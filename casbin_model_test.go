package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/stretchr/testify/assert"
)

func TestCasbinModel(t *testing.T) {
	_, err1 := caskin.CasbinModel(false)
	assert.NoError(t, err1)
	_, err2 := caskin.CasbinModel(true)
	assert.NoError(t, err2)
}

func TestCasbinModelSuperadmin(t *testing.T) {
	_, err := caskin.CasbinModelSuperadmin()
	assert.NoError(t, err)
}

func TestCasbinModelNoSuperadmin(t *testing.T) {
	_, err := caskin.CasbinModelNoSuperadmin()
	assert.NoError(t, err)
}
