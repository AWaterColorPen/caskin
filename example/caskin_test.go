package example_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestNewCaskin(t *testing.T) {
	_, err := example.NewCaskin(&caskin.Options{}, t.TempDir())
	assert.NoError(t, err)
}
