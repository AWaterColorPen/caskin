package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestServer_UserRole_GetUserRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	list1, err := service.GetUserRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list1, 2)
	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	pairs := []*caskin.UserRolePair{
		{stage.Admin, roles[0]},
		{stage.Member, roles[0]},
	}
	assert.NoError(t, service.ModifyUserRolePerRole(stage.Admin, stage.Domain, roles[0], pairs))

	list2, err := service.GetUserRole(stage.Member, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list2, 3)
}

func TestServer_UserRole_GetUserRoleByRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	list1, err := service.GetUserRoleByRole(stage.Admin, stage.Domain, roles[0])
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	list2, err := service.GetUserRoleByRole(stage.Admin, stage.Domain, roles[1])
	assert.NoError(t, err)
	assert.Len(t, list2, 1)

	pairs := []*caskin.UserRolePair{
		{stage.Admin, roles[0]},
		{stage.Member, roles[0]},
	}
	assert.NoError(t, service.ModifyUserRolePerRole(stage.Admin, stage.Domain, roles[0], pairs))

	list3, err := service.GetUserRoleByRole(stage.Admin, stage.Domain, roles[0])
	assert.NoError(t, err)
	assert.Len(t, list3, 2)
	list4, err := service.GetUserRoleByRole(stage.Admin, stage.Domain, roles[1])
	assert.NoError(t, err)
	assert.Len(t, list4, 1)
}

func TestServer_UserRole_GetUserRoleByUser(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	list1, err := service.GetUserRoleByUser(stage.Admin, stage.Domain, stage.Admin)
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	list2, err := service.GetUserRoleByUser(stage.Admin, stage.Domain, stage.Member)
	assert.NoError(t, err)
	assert.Len(t, list2, 1)

	pairs := []*caskin.UserRolePair{
		{stage.Admin, roles[0]},
		{stage.Member, roles[0]},
	}
	assert.NoError(t, service.ModifyUserRolePerRole(stage.Admin, stage.Domain, roles[0], pairs))

	list3, err := service.GetUserRoleByUser(stage.Admin, stage.Domain, stage.Admin)
	assert.NoError(t, err)
	assert.Len(t, list3, 1)
	list4, err := service.GetUserRoleByUser(stage.Admin, stage.Domain, stage.Member)
	assert.NoError(t, err)
	assert.Len(t, list4, 2)
}

func TestServer_UserRole_ModifyUserRolePerRole(t *testing.T) {
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
	assert.Equal(t, caskin.ErrNoWritePermission, service.ModifyUserRolePerRole(stage.Member, stage.Domain, roles[0], pairs))
	assert.Equal(t, caskin.ErrInputPairArrayNotBelongSameRole, service.ModifyUserRolePerRole(stage.Admin, stage.Domain, roles[0], pairs))
	assert.NoError(t, service.ModifyUserRolePerRole(stage.Admin, stage.Domain, role1, pairs))
}

func TestServer_UserRole_ModifyUserRolePerUser(t *testing.T) {
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
	assert.NoError(t, service.ModifyUserRolePerUser(stage.Member, stage.Domain, stage.Member, pairs))
	list1, err := service.GetUserRoleByUser(stage.Member, stage.Domain, stage.Member)
	assert.NoError(t, err)
	assert.Len(t, list1, 1)

	assert.Equal(t, caskin.ErrInputPairArrayNotBelongSameUser, service.ModifyUserRolePerUser(stage.Admin, stage.Domain, stage.Admin, pairs))
	assert.NoError(t, service.ModifyUserRolePerUser(stage.Admin, stage.Domain, stage.Member, pairs))
	list2, err := service.GetUserRoleByUser(stage.Member, stage.Domain, stage.Member)
	assert.NoError(t, err)
	assert.Len(t, list2, 2)
}
