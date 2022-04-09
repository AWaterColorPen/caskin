package playground_test

import (
	"github.com/awatercolorpen/caskin/playground"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayground(t *testing.T) {
	stage, err := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	assert.NoError(t, stage.AddSubAdmin())
}
