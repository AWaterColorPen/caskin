package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorObject(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := caskin.NewCachedProvider(nil, nil)
	provider.User = stage.SuperadminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	objects, err := executor.GetObjects(caskin.ObjectTypeObject)
	assert.NoError(t, err)
	assert.Len(t, objects, 1)

	domains, err := executor.GetAllDomain()
	assert.NoError(t, err)
	assert.Len(t, domains, 1)

	object := &example.Object{
		Name:     "object_01",
		Type:     ObjectTypeTest,
		ObjectID: objects[0].GetID(),
	}
	assert.NoError(t, executor.CreateObject(object))

	subObject := &example.Object{
		Name:     "object_01_sub",
		Type:     ObjectTypeTest,
		ObjectID: objects[0].GetID(),
		ParentID: object.ID,
	}
	assert.NoError(t, executor.CreateObject(subObject))

	assert.NoError(t, executor.DeleteObject(object))
	objects, err = executor.GetObjects(ObjectTypeTest)
	assert.NoError(t, err)
	assert.Len(t, objects, 0)

	assert.NoError(t, executor.RecoverObject(object))
	objects, err = executor.GetObjects(ObjectTypeTest)
	assert.NoError(t, err)
	assert.Len(t, objects, 1)

	object.Name = "object_01_new_name"
	assert.NoError(t, executor.UpdateObject(object))

	assert.NoError(t, executor.RecoverObject(subObject))
	objects, err = executor.GetObjects(ObjectTypeTest)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	assert.NoError(t, executor.DeleteObject(object))
	objects, err = executor.GetObjects(ObjectTypeTest)
	assert.NoError(t, err)
	assert.Len(t, objects, 0)

	assert.NoError(t, executor.RecoverObject(subObject))
	objects, err = executor.GetObjects(ObjectTypeTest)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	subObject.Name = "object_01__sub_new_name"
	assert.NoError(t, executor.UpdateObject(subObject))
}

func TestExecutorObject_GetObjects(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	objects1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects1, 3)

	objects2, err := executor.GetObjects(caskin.ObjectTypeObject)
	assert.NoError(t, err)
	assert.Len(t, objects2, 1)

	objects3, err := executor.GetObjects(caskin.ObjectTypeObject, caskin.ObjectTypeRole)
	assert.NoError(t, err)
	assert.Len(t, objects3, 2)

	provider.Domain = stage.Domain
	provider.User = stage.MemberUser
	objects4, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects4, 0)
}

func TestExecutorObject_GeneralCreate(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	object1 := &example.Object{
		Name: "object_01",
		Type: ObjectTypeTest,
	}
	assert.Equal(t, caskin.ErrProviderGet, executor.CreateObject(object1))

	provider.Domain = stage.Domain
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObject(object1))
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObject(object1))

	object1.ObjectID = 1
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObject(object1))
	provider.User = stage.AdminUser
	assert.NoError(t, executor.CreateObject(object1))

	object2 := &example.Object{
		Name: "object_01",
		Type: ObjectTypeTest,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateObject(object2))

	object3 := &example.Object{
		Name: "object_01",
		Type: ObjectTypeTest,
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteObject(object3))
	object3.ID = object1.ID
	assert.NoError(t, executor.DeleteObject(object3))

	object4 := &example.Object{ID: 10, ObjectID: 1}
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteObject(object4))
	assert.NoError(t, executor.CreateObject(object4))
}

func TestExecutorObject_CreateSubNode(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	object1 := &example.Object{
		Name:     "sub_object_1",
		ObjectID: 4,
		ParentID: 4,
	}
	provider.Domain = stage.Domain
	provider.User = stage.SubAdminUser
	assert.NoError(t, executor.CreateObject(object1))

	object1.ParentID = 1
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateObject(object1))
	provider.User = stage.AdminUser
	assert.NoError(t, executor.UpdateObject(object1))

	object2 := &example.Object{
		Name:     "object_2",
		ObjectID: 2,
	}
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObject(object2))
	provider.User = stage.AdminUser
	assert.NoError(t, executor.CreateObject(object2))

	object3 := &example.Object{
		Name:     "sub_admin_1",
		ObjectID: 1,
	}
	assert.NoError(t, executor.CreateObject(object3))
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateObject(object3))

	object4 := &example.Object{
		Name: "sub_admin_1",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteObject(object4))
	object4.ID = object3.ID
	assert.NoError(t, executor.DeleteObject(object4))

	object5 := &example.Object{ID: 8, ObjectID: 1}
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteObject(object5))
	assert.NoError(t, executor.CreateObject(object5))
}

func TestExecutorObject_GeneralUpdate(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil)
	provider.User = stage.AdminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	object := &example.Object{
		Name:     "object_01",
		Type:     ObjectTypeTest,
		ObjectID: 1,
	}
	assert.NoError(t, executor.CreateObject(object))

	object.Name = "object_01_new_name"
	assert.NoError(t, executor.UpdateObject(object))

	subObject := &example.Object{
		Name:     "object_01_sub",
		Type:     ObjectTypeTest,
		ObjectID: 1,
		ParentID: object.ID,
	}
	assert.NoError(t, executor.CreateObject(subObject))

	provider.User = stage.MemberUser
	subObject.Name = "object_01_sub_new_name"
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateObject(subObject))

	subObject.ID = 0
	assert.Equal(t, caskin.ErrEmptyID, executor.UpdateObject(subObject))
	subObject.ID = 10
	assert.Equal(t, caskin.ErrNotExists, executor.UpdateObject(subObject))

	object2 := &example.Object{
		Name:     "object_02",
		Type:     ObjectTypeTest,
		ObjectID: 1,
	}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObject(object2))
}

func TestExecutorObject_GeneralRecover(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil)
	provider.User = stage.SuperadminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	objects, _ := executor.GetObjects()
	assert.Len(t, objects, 3)

	assert.NoError(t, executor.DeleteObject(objects[2]))

	objects, _ = executor.GetObjects()
	assert.Len(t, objects, 2)

	object := &example.Object{
		ID: 3,
	}
	assert.Error(t, executor.DeleteObject(object))

	provider.User = stage.MemberUser
	object.ID = 2
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteObject(object))

	provider.User = stage.AdminUser
	assert.NoError(t, executor.DeleteObject(object))
	assert.NoError(t, executor.RecoverObject(object))
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverObject(object))
}

func TestExecutorObject_GeneralDelete(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil)
	provider.User = stage.SuperadminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	objects, _ := executor.GetObjects()
	assert.Len(t, objects, 3)

	assert.NoError(t, executor.DeleteObject(objects[2]))

	objects, _ = executor.GetObjects()
	assert.Len(t, objects, 2)

	object := &example.Object{
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteObject(object))
	object.ID = 4
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteObject(object))
	assert.NoError(t, executor.CreateObject(object))

	object.Name = "object_01_new_name"
	assert.NoError(t, executor.UpdateObject(object))

	subObject := &example.Object{
		Name:     "object_01_sub",
		Type:     ObjectTypeTest,
		DomainID: 1,
		ObjectID: 4,
		ParentID: object.ID,
	}
	assert.NoError(t, executor.CreateObject(subObject))

	provider.User = stage.AdminUser
	subObject.Name = "object_01_sub_new_name"
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateObject(subObject))
}
