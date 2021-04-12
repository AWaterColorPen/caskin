package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestTreeNodeEntryDeleter_RetryWithRelation(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	assert.NoError(t, stage.AddSubAdmin())

	executor := stage.Caskin.GetExecutor(provider)

	object1 := &example.Object{
		Name:     "object_sub_02",
		Type:     caskin.ObjectTypeObject,
		ObjectID: 4,
		ParentID: 4,
	}

	provider.Domain = stage.Domain
	provider.User = stage.SubAdminUser
	assert.NoError(t, executor.CreateObject(object1))

	provider.User = stage.SuperadminUser
	list1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, list1, 6)

	assert.NoError(t, executor.Enforcer.RemoveObjectInDomain(list1[5], stage.Domain))
	list2, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, list2, 6)

	provider.User = stage.AdminUser
	assert.NoError(t, executor.DeleteObject(list2[3]))

	provider.User = stage.SuperadminUser
	list3, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, list3, 4)
}
