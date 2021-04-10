package example_test

import (
	"github.com/awatercolorpen/caskin/extension/manager"
	"testing"

	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestNewCaskin(t *testing.T) {
	_, err := example.NewManager(t.TempDir())
	assert.NoError(t, err)

	_, err = example.NewManager(t.TempDir(), func(configuration *manager.Configuration) {
		configuration.Extension = &manager.Extension{
			WebFeature: 0,
		}
	})
	assert.NoError(t, err)
}
