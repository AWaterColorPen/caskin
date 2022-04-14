package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServer_GetObject(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects1, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects1, 2)

	objects2, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage, caskin.ObjectTypeRole)
	assert.NoError(t, err)
	assert.Len(t, objects2, 1)

	objects3, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage, caskin.ObjectTypeRole, "test")
	assert.NoError(t, err)
	assert.Len(t, objects3, 2)

	objects4, err := service.GetObject(stage.Member, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects4, 0)

	objects5, err := service.GetObject(stage.Member, stage.Domain, caskin.Read)
	assert.NoError(t, err)
	assert.Len(t, objects5, 2)

	objects6, err := service.GetObject(stage.Member, stage.Domain, caskin.Write)
	assert.NoError(t, err)
	assert.Len(t, objects6, 1)

	objects7, err := service.GetObject(stage.Superadmin, caskin.GetSuperadminDomain(), caskin.Read)
	assert.NoError(t, err)
	assert.Len(t, objects7, 0)
}

func TestServer_CreateObject(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{Name: "test_sub_01", Type: "test"}
	assert.Equal(t, caskin.ErrCantOperateRootObject, service.CreateObject(stage.Member, stage.Domain, object1))
	assert.Equal(t, caskin.ErrCantOperateRootObject, service.CreateObject(stage.Admin, stage.Domain, object1))
	assert.Equal(t, caskin.ErrCantOperateRootObject, service.CreateObject(stage.Admin, stage.Domain, object1))
	object1.ParentID = objects[0].GetID()
	assert.Equal(t, caskin.ErrInValidObjectType, service.CreateObject(stage.Admin, stage.Domain, object1))
	object1.ParentID = objects[1].GetID()
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{Name: "test_sub_01", Type: "test"}
	assert.Equal(t, caskin.ErrAlreadyExists, service.CreateObject(stage.Admin, stage.Domain, object2))

	object3 := &example.Object{Name: "test_sub_01", Type: "test"}
	assert.Equal(t, caskin.ErrEmptyID, service.DeleteObject(stage.Admin, stage.Domain, object3))
	object3.ID = object1.ID
	assert.NoError(t, service.DeleteObject(stage.Admin, stage.Domain, object3))

	object4 := &example.Object{ID: 999}
	assert.Equal(t, caskin.ErrCantOperateRootObject, service.CreateObject(stage.Admin, stage.Domain, object4))
}

func TestServer_ObjectTreeNode(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	object1 := &example.Object{Name: "test_sub_01", Type: "test"}
	object1.ParentID = objects[1].GetID()
	assert.Equal(t, caskin.ErrNoManagePermission, service.CreateObject(stage.Member, stage.Domain, object1))
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{Name: "test_sub_02", Type: "test"}
	object2.ParentID = object1.GetID()
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object2))

	policy := []*caskin.Policy{
		{Role: roles[1], Object: object2, Domain: stage.Domain, Action: caskin.Read},
		{Role: roles[1], Object: object2, Domain: stage.Domain, Action: caskin.Write},
		{Role: roles[1], Object: object2, Domain: stage.Domain, Action: caskin.Manage},
	}
	assert.NoError(t, service.ModifyPolicyPerRole(stage.Admin, stage.Domain, roles[1], policy))

	// delete object tree node
	objects2, err := service.GetObject(stage.Member, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects2, 1)
	assert.Equal(t, caskin.ErrNoManagePermission, service.DeleteObject(stage.Member, stage.Domain, object1))

	assert.NoError(t, service.DeleteObject(stage.Admin, stage.Domain, object1))
	objects3, err := service.GetObject(stage.Member, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects3, 0)

	objects4, err := service.GetObject(stage.Superadmin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects4, 2)

	// recover object tree node
	assert.NoError(t, service.RecoverObject(stage.Admin, stage.Domain, object1))
	objects5, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects5, 3)

	assert.NoError(t, service.RecoverObject(stage.Admin, stage.Domain, object2))
	objects6, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects6, 4)

	objects7, err := service.GetObject(stage.Member, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects7, 0)
}

func TestServer_UpdateObject(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{ID: objects[0].GetID(), Name: "object_01_new_name"}
	assert.Equal(t, caskin.ErrCantOperateRootObject, service.UpdateObject(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{Name: "test_sub_01", Type: "test"}
	object2.ParentID = objects[1].GetID()
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object2))

	object3 := &example.Object{ID: object2.GetID(), Name: "test_sub_01_new_name"}
	assert.Equal(t, caskin.ErrNoManagePermission, service.UpdateObject(stage.Member, stage.Domain, object3))
	assert.Equal(t, caskin.ErrCantOperateRootObject, service.UpdateObject(stage.Admin, stage.Domain, object3))
	object3.ParentID = object2.GetParentID()
	assert.Equal(t, caskin.ErrInValidObjectType, service.UpdateObject(stage.Admin, stage.Domain, object3))
	object3.Type = object2.Type
	assert.NoError(t, service.UpdateObject(stage.Admin, stage.Domain, object3))

	object3.Type = caskin.ObjectTypeRole
	assert.Equal(t, caskin.ErrCantChangeObjectType, service.UpdateObject(stage.Admin, stage.Domain, object3))
}

func TestServer_UpdateObject_Parent(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{Name: "test_sub_01", Type: "test"}
	object1.ParentID = objects[1].GetID()
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{Name: "test_sub_02", Type: "test"}
	object2.ParentID = object1.GetID()
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object2))

	object3 := &example.Object{
		ID:       object1.GetID(),
		Name:     "change_parent_id_from_1_to_2",
		Type:     "test",
		ParentID: object2.ID,
	}
	assert.Equal(t, caskin.ErrParentToDescendant, service.UpdateObject(stage.Admin, stage.Domain, object3))
	object3.ParentID = object3.GetID()
	assert.Equal(t, caskin.ErrParentCanNotBeItself, service.UpdateObject(stage.Admin, stage.Domain, object3))
	object3.ParentID = objects[0].GetID()
	assert.Equal(t, caskin.ErrInValidObjectType, service.UpdateObject(stage.Admin, stage.Domain, object3))

	object4 := &example.Object{
		ID:       object2.GetID(),
		Name:     "change_parent_id_from_2_to_0",
		Type:     "test",
		ParentID: objects[1].GetID(),
	}
	assert.NoError(t, service.UpdateObject(stage.Admin, stage.Domain, object4))
}

func TestServer_RecoverObject(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{ID: objects[0].GetID()}
	assert.Equal(t, caskin.ErrAlreadyExists, service.RecoverObject(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{Name: "test_sub_01", Type: "test"}
	object2.ParentID = objects[1].GetID()
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object2))

	assert.Equal(t, caskin.ErrNoManagePermission, service.DeleteObject(stage.Member, stage.Domain, object2))
	assert.NoError(t, service.DeleteObject(stage.Admin, stage.Domain, object2))

	object3 := &example.Object{ID: object2.GetID()}
	assert.Equal(t, caskin.ErrNoManagePermission, service.RecoverObject(stage.Member, stage.Domain, object3))
	assert.NoError(t, service.RecoverObject(stage.Admin, stage.Domain, object3))

	object4 := &example.Object{ID: 999}
	assert.Equal(t, caskin.ErrNotExists, service.RecoverObject(stage.Admin, stage.Domain, object4))
}

func TestServer_DeleteObject(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{}
	assert.Equal(t, caskin.ErrEmptyID, service.DeleteObject(stage.Admin, stage.Domain, object1))

	object1.ID = objects[1].GetID()
	assert.Equal(t, caskin.ErrCantOperateRootObject, service.DeleteObject(stage.Member, stage.Domain, object1))
	assert.Equal(t, caskin.ErrCantOperateRootObject, service.DeleteObject(stage.Admin, stage.Domain, object1))
	assert.Equal(t, caskin.ErrCantOperateRootObject, service.DeleteObject(stage.Superadmin, stage.Domain, object1))

	object2 := &example.Object{ID: 999}
	assert.Equal(t, caskin.ErrNotExists, service.DeleteObject(stage.Admin, stage.Domain, object2))
}
