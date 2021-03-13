package example_test

import (
	"testing"

	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestNewStage(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	assert.NoError(t, stage.AddSubAdmin())
	assert.NoError(t, stage.NoSuperadmin())
}
