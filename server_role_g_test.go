package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestServer_AddRoleG(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles, _ := service.GetRole(stage.Member, stage.Domain)
	assert.Len(t, roles, 2)

	role1 := &example.Role{Name: "role_01", ObjectID: roles[0].GetObjectID()}
	assert.NoError(t, service.CreateRole(stage.Admin, stage.Domain, role1))

	pairs := []*caskin.UserRolePair{
		{stage.Member, role1},
	}
	assert.NoError(t, service.ModifyUserRolePerRole(stage.Admin, stage.Domain, role1, pairs))

	objects1, err := service.GetObject(stage.Member, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects1, 0)

	assert.Equal(t, caskin.ErrNoWritePermission, service.AddRoleG(stage.Member, stage.Domain, role1, roles[0]))
	assert.NoError(t, service.AddRoleG(stage.Admin, stage.Domain, role1, roles[0]))

	objects2, err := service.GetObject(stage.Member, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects2, 2)
}

func TestServer_RemoveRoleG(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles, _ := service.GetRole(stage.Member, stage.Domain)
	assert.Len(t, roles, 2)

	role1 := &example.Role{Name: "role_01", ObjectID: roles[0].GetObjectID()}
	assert.NoError(t, service.CreateRole(stage.Admin, stage.Domain, role1))

	pairs := []*caskin.UserRolePair{
		{stage.Member, role1},
	}
	assert.NoError(t, service.ModifyUserRolePerRole(stage.Admin, stage.Domain, role1, pairs))
	assert.NoError(t, service.AddRoleG(stage.Admin, stage.Domain, role1, roles[0]))

	objects1, err := service.GetObject(stage.Member, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects1, 2)

	assert.NoError(t, service.RemoveRoleG(stage.Member, stage.Domain, role1, roles[0]))

	objects2, err := service.GetObject(stage.Member, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects2, 0)
}
