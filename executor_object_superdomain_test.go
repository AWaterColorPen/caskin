package caskin_test

import (
	"github.com/awatercolorpen/caskin/playground"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorObject_Superdomain(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	
	provider.User = stage.SuperadminUser
	provider.Domain = stage.Options.GetSuperadminDomain()
	service := stage.Service

	object1 := &example.Object{
		Name: "object_01",
		Type: ObjectTypeTest,
	}
	assert.NoError(t, executor.CreateObject(object1))

	list1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, list1, 1)

	provider.User = stage.AdminUser
	list2, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, list2, 0)

	provider.User = stage.SuperadminUser
	object1.ObjectID = object1.ID
	assert.Equal(t, caskin.ErrInValidObjectType, executor.UpdateObject(object1))
	object1.ObjectID = 4
	assert.NoError(t, executor.UpdateObject(object1))

	object2 := &example.Object{
		Name:     "object_02",
		Type:     ObjectTypeTest,
		ObjectID: object1.ID,
		ParentID: 1,
	}
	assert.Equal(t, caskin.ErrInValidObjectType, executor.CreateObject(object2))
	object2.ParentID = object1.ID
	assert.Equal(t, caskin.ErrInValidObjectType, executor.CreateObject(object2))
	object2.ObjectID = 1
	assert.NoError(t, executor.CreateObject(object2))

	assert.NoError(t, executor.DeleteObject(object1))
	list3, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, list3, 0)

	assert.NoError(t, executor.RecoverObject(object1))
	list4, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, list4, 1)
}

func TestExecutorObject_Superdomain_NoSuperadmin(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	assert.NoError(t, stage.NoSuperadmin())
	
	provider.User = stage.SuperadminUser
	provider.Domain = stage.Options.GetSuperadminDomain()
	service := stage.Service

	object1 := &example.Object{
		Name: "object_01",
		Type: ObjectTypeTest,
	}
	assert.Error(t, executor.CreateObject(object1))
}
