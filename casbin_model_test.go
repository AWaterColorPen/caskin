package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/stretchr/testify/assert"
)

func TestCasbinModel(t *testing.T) {
	_, err := caskin.CasbinModel()
	assert.NoError(t, err)
}

func TestCasbinModelText(t *testing.T) {
	_, err := caskin.CasbinModelText()
	assert.NoError(t, err)
}
