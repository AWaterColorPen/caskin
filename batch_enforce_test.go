package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestEnforcer_BatchEnforce(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	e := stage.Enforcer
	assert.NotNil(t, e)

	// Get objects visible to admin
	objects, err := stage.Service.GetObject(stage.Admin, stage.Domain, caskin.Read)
	assert.NoError(t, err)
	assert.NotEmpty(t, objects)

	// BatchEnforce should return same results as individual Enforce calls
	results, err := e.BatchEnforce(stage.Admin, objects, stage.Domain, caskin.Read)
	assert.NoError(t, err)
	assert.Len(t, results, len(objects))

	for i, obj := range objects {
		individual, err := e.Enforce(stage.Admin, obj, stage.Domain, caskin.Read)
		assert.NoError(t, err)
		assert.Equal(t, individual, results[i], "mismatch at index %d for object %v", i, obj)
	}
}

func TestEnforcer_BatchEnforce_Empty(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	e := stage.Enforcer

	results, err := e.BatchEnforce(stage.Admin, nil, stage.Domain, caskin.Read)
	assert.NoError(t, err)
	assert.Nil(t, results)

	results, err = e.BatchEnforce(stage.Admin, []caskin.Object{}, stage.Domain, caskin.Read)
	assert.NoError(t, err)
	assert.Nil(t, results)
}

func TestEnforcer_BatchEnforce_MixedPermissions(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	e := stage.Enforcer

	// Member should have Read on some objects but not Manage on all
	objects, err := stage.Service.GetObject(stage.Admin, stage.Domain, caskin.Read)
	assert.NoError(t, err)
	assert.NotEmpty(t, objects)

	// Check member's manage permission — should have fewer allowed
	results, err := e.BatchEnforce(stage.Member, objects, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, results, len(objects))

	// Verify consistency with individual calls
	for i, obj := range objects {
		individual, err := e.Enforce(stage.Member, obj, stage.Domain, caskin.Manage)
		assert.NoError(t, err)
		assert.Equal(t, individual, results[i])
	}
}

func TestFilter_UsesBatchEnforce(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	e := stage.Enforcer

	// Get all objects (admin can see all)
	allObjects, err := stage.Service.GetObject(stage.Admin, stage.Domain, caskin.Read)
	assert.NoError(t, err)

	// Filter for member with Read — should match GetObject result
	filtered := caskin.Filter(e, stage.Member, stage.Domain, caskin.Read, allObjects)

	memberObjects, err := stage.Service.GetObject(stage.Member, stage.Domain, caskin.Read)
	assert.NoError(t, err)

	assert.Len(t, filtered, len(memberObjects))
}

func TestFilter_EmptySource(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	e := stage.Enforcer

	result := caskin.Filter(e, stage.Admin, stage.Domain, caskin.Read, []caskin.Object(nil))
	assert.Nil(t, result)

	result = caskin.Filter(e, stage.Admin, stage.Domain, caskin.Read, []caskin.Object{})
	assert.Nil(t, result)
}
