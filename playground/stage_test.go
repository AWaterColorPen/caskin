package playground_test

import (
	"github.com/awatercolorpen/caskin/playground"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStage(t *testing.T) {
	stage, err := playground.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	assert.NoError(t, stage.AddSubAdmin())
	assert.NoError(t, stage.NoSuperadmin())
}
