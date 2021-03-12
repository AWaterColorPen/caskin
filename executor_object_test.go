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
	assert.NoError(t, executor.DeleteObject(object2))

	provider.User = stage.AdminUser
	list1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, list1, 3)
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
}

func TestExecutorObject_GeneralRecover(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider.User = stage.AdminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	object := &example.Object{
		ID:       4,
		ParentID: 1,
	}
	assert.NoError(t, executor.DeleteObject(object))

	object1 := &example.Object{}
	assert.Equal(t, caskin.ErrEmptyID, executor.RecoverObject(object1))

	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.RecoverObject(object))
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverObject(object))

	object2 := &example.Object{ID: 3}
	assert.NoError(t, executor.DeleteObject(object2))
	provider.User = stage.MemberUser
	object2.ID = 3
	assert.Equal(t, caskin.ErrNoWritePermission, executor.RecoverObject(object2))
}

/**
  删除object的测试：
  1. 成功
  	1.1 测试权限是否正常工作，admin是否能够删除子节点的数据
  	1.2 测试是否能够成功删除
	1.3 测试
  2. 失败
	2.1 传入的值不规范，可能id为0
	2.2 数据不存在，无法删除
	2.3 当前用户没有写的权限
	2.4 和当前的objectType不同而不能删除
*/
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
	assert.Equal(t, caskin.ErrEmptyParentIdOrNotSuperadmin, executor.DeleteObject(object1))

	object2 := &example.Object{ID: 5, ParentID: 1}
	assert.Equal(t, caskin.ErrInValidObjectType, executor.DeleteObject(object2))

	object1.ParentID = 1
	assert.NoError(t, executor.DeleteObject(object1))

	assert.Equal(t, caskin.ErrNotExists, executor.DeleteObject(object1))

	provider.User = stage.AdminUser
	object3 := &example.Object{ID: 5}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteObject(object3))
}
