package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestServer_Policy_GetPolicy(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	list1, err := service.GetPolicy(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list1, 9)

	policy1 := []*caskin.Policy{
		{roles[1], objects[0], stage.Domain, caskin.Read},
		{roles[1], objects[1], stage.Domain, caskin.Read},
	}
	assert.NoError(t, service.ModifyPolicyPerRole(stage.Admin, stage.Domain, roles[1], policy1))

	list2, err := service.GetPolicy(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list2, 8)

	list3, err := service.GetPolicy(stage.Member, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list3, 0)
}

func TestServer_Policy_GetPolicyByRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	list1, err := service.GetPolicyByRole(stage.Admin, stage.Domain, roles[0])
	assert.NoError(t, err)
	assert.Len(t, list1, 6)
	list2, err := service.GetPolicyByRole(stage.Admin, stage.Domain, roles[1])
	assert.NoError(t, err)
	assert.Len(t, list2, 3)

	policy1 := []*caskin.Policy{
		{roles[1], objects[0], stage.Domain, caskin.Read},
		{roles[1], objects[1], stage.Domain, caskin.Read},
	}
	assert.NoError(t, service.ModifyPolicyPerRole(stage.Admin, stage.Domain, roles[1], policy1))

	list3, err := service.GetPolicyByRole(stage.Admin, stage.Domain, roles[0])
	assert.NoError(t, err)
	assert.Len(t, list3, 6)

	list4, err := service.GetPolicyByRole(stage.Admin, stage.Domain, roles[1])
	assert.NoError(t, err)
	assert.Len(t, list4, 2)
}

func TestServer_Policy_GetPolicyByObject(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	list1, err := service.GetPolicyByObject(stage.Admin, stage.Domain, objects[0])
	assert.NoError(t, err)
	assert.Len(t, list1, 4)
	list2, err := service.GetPolicyByObject(stage.Admin, stage.Domain, objects[1])
	assert.NoError(t, err)
	assert.Len(t, list2, 5)

	policy1 := []*caskin.Policy{
		{roles[1], objects[0], stage.Domain, caskin.Read},
		{roles[1], objects[1], stage.Domain, caskin.Read},
	}
	assert.NoError(t, service.ModifyPolicyPerRole(stage.Admin, stage.Domain, roles[1], policy1))

	list3, err := service.GetPolicyByObject(stage.Admin, stage.Domain, objects[0])
	assert.NoError(t, err)
	assert.Len(t, list3, 4)

	list4, err := service.GetPolicyByObject(stage.Admin, stage.Domain, objects[1])
	assert.NoError(t, err)
	assert.Len(t, list4, 4)
}

func TestServer_Policy_ModifyPolicyPerRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	object1 := &example.Object{Name: "role_sub_01", Type: "role"}
	object1.ParentID = objects[0].GetID()
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	role1 := &example.Role{Name: "role_sub_01", ObjectID: object1.ID}
	assert.NoError(t, service.CreateRole(stage.Admin, stage.Domain, role1))

	pairs := []*caskin.UserRolePair{
		{stage.Member, role1},
	}
	assert.NoError(t, service.ModifyUserRolePerRole(stage.Admin, stage.Domain, role1, pairs))

	list1, err := service.GetPolicyByRole(stage.Admin, stage.Domain, role1)
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	policy1 := []*caskin.Policy{
		{role1, object1, stage.Domain, caskin.Manage},
	}
	assert.NoError(t, service.ModifyPolicyPerRole(stage.Admin, stage.Domain, role1, policy1))

	list2, err := service.GetPolicyByRole(stage.Member, stage.Domain, role1)
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
	assert.Equal(t, caskin.ErrNoManagePermission, service.DeleteObject(stage.Member, stage.Domain, object1))
	assert.NoError(t, service.DeleteObject(stage.Admin, stage.Domain, object1))

	_, err = service.GetPolicyByRole(stage.Member, stage.Domain, role1)
	assert.Equal(t, caskin.ErrNoReadPermission, err)
}

func TestServer_Policy_ModifyPolicyPerObject(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	object1 := &example.Object{Name: "role_sub_01", Type: "role"}
	object1.ParentID = objects[0].GetID()
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	role1 := &example.Role{Name: "role_sub_01", ObjectID: object1.ID}
	assert.NoError(t, service.CreateRole(stage.Admin, stage.Domain, role1))

	pairs := []*caskin.UserRolePair{
		{stage.Member, role1},
	}
	assert.NoError(t, service.ModifyUserRolePerRole(stage.Admin, stage.Domain, role1, pairs))

	list1, err := service.GetPolicyByRole(stage.Admin, stage.Domain, role1)
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	policy1 := []*caskin.Policy{
		{role1, object1, stage.Domain, caskin.Manage},
	}
	assert.NoError(t, service.ModifyPolicyPerObject(stage.Admin, stage.Domain, object1, policy1))

	list2, err := service.GetPolicyByObject(stage.Member, stage.Domain, object1)
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
	assert.Equal(t, caskin.ErrNoManagePermission, service.DeleteObject(stage.Member, stage.Domain, object1))
	assert.NoError(t, service.DeleteObject(stage.Admin, stage.Domain, object1))

	_, err = service.GetPolicyByObject(stage.Member, stage.Domain, object1)
	assert.Equal(t, caskin.ErrNotExists, err)
}
