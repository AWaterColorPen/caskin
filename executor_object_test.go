package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, caskin.ErrCanNotOperateRootObjectWithoutSuperadmin, executor.CreateObject(object1))
	object1.ParentID = 3
	assert.Equal(t, caskin.ErrInValidObjectType, executor.CreateObject(object1))
	object1.Type = caskin.ObjectTypeDefault
	assert.NoError(t, executor.CreateObject(object1))

	object2 := &example.Object{
		Name: "object_01",
		Type: caskin.ObjectTypeDefault,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateObject(object2))

	object3 := &example.Object{
		Name: "object_01",
		Type: caskin.ObjectTypeDefault,
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteObject(object3))
	object3.ID = object1.ID
	assert.NoError(t, executor.DeleteObject(object3))

	object4 := &example.Object{ID: 10, ObjectID: 1}
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteObject(object4))
	assert.Equal(t, caskin.ErrCanNotOperateRootObjectWithoutSuperadmin, executor.CreateObject(object4))
}

func TestExecutorObject_CreateSubNode(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	object1 := &example.Object{
		Name:     "object_sub_02",
		Type:     caskin.ObjectTypeObject,
		ObjectID: 4,
		ParentID: 4,
	}
	provider.User = stage.MemberUser
	provider.Domain = stage.Domain
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObject(object1))

	provider.User = stage.SubAdminUser
	assert.NoError(t, executor.CreateObject(object1))

	object1.ParentID = 1
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateObject(object1))
	object1.ParentID = 5
	assert.Equal(t, caskin.ErrInValidObjectType, executor.UpdateObject(object1))
	object1.ParentID = 4
	assert.Equal(t, caskin.ErrObjectTypeObjectIDMustBeItselfID, executor.UpdateObject(object1))

	object1.ObjectID = object1.ID
	assert.NoError(t, executor.UpdateObject(object1))

	object2 := &example.Object{ID: 4}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteObject(object2))
	object3 := &example.Object{ID: 6}
	assert.NoError(t, executor.DeleteObject(object3))

	provider.User = stage.AdminUser
	assert.NoError(t, executor.DeleteObject(object2))
	list1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, list1, 3)

	provider.User = stage.SuperadminUser
	list2, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, list2, 4)
}

func TestExecutorObject_GeneralUpdate(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil)
	assert.NoError(t, stageAddSubAdmin(stage))

	provider.User = stage.AdminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	object := &example.Object{
		ID:       4,
		Name:     "object_sub_01_new_name",
		Type:     caskin.ObjectTypeObject,
		ParentID: 1,
		ObjectID: 4,
	}
	assert.NoError(t, executor.UpdateObject(object))

	object.Type = caskin.ObjectTypeRole
	assert.Equal(t, caskin.ErrInValidObjectType, executor.UpdateObject(object))

	object1 := &example.Object{}
	assert.Equal(t, caskin.ErrEmptyID, executor.UpdateObject(object1))

	object2 := &example.Object{ID: 1, Name: "object_01_new_name", Type: caskin.ObjectTypeObject, ObjectID: 1, ParentID: 0}
	assert.Equal(t, caskin.ErrCanNotOperateRootObjectWithoutSuperadmin, executor.UpdateObject(object2))

	provider.User = stage.MemberUser
	object4 := &example.Object{
		ID:       4,
		Name:     "object_sub_01_new_name2",
		Type:     caskin.ObjectTypeObject,
		ParentID: 1,
		ObjectID: 4,
	}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateObject(object4))

}

func TestExecutorObject_GeneralRecover(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := caskin.NewCachedProvider(nil, nil)
	provider.User = stage.AdminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	object := &example.Object{
		ID:       4,
		ParentID: 1,
	}
	assert.NoError(t, executor.DeleteObject(object))

	object1 := &example.Object{}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverObject(object1))

	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.RecoverObject(object))
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverObject(object))

	object2 := &example.Object{ID: 3}
	assert.NoError(t, executor.DeleteObject(object2))
	provider.User = stage.MemberUser
	object2.ID = 3
	assert.Equal(t, caskin.ErrNoWritePermission, executor.RecoverObject(object2))
}

func TestExecutorObject_GeneralDelete(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil)
	assert.NoError(t, stageAddSubAdmin(stage))

	provider.User = stage.SubAdminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	object := &example.Object{}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteObject(object))

	object1 := &example.Object{ID: 4}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteObject(object1))

	provider.User = stage.AdminUser
	assert.NoError(t, executor.DeleteObject(object1))

	object2 := &example.Object{ID: 5}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteObject(object2))

	object4 := &example.Object{ID: 1}
	assert.Equal(t, caskin.ErrCanNotOperateRootObjectWithoutSuperadmin, executor.DeleteObject(object4))

	// TODO: if object type == object, it will not recover it by admin user now
	// want to support it by a special API
	object3 := &example.Object{ID: 4, ParentID: 1}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.RecoverObject(object3))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.RecoverObject(object3))
}
