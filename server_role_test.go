package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestServer_GetRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles1, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles1, 2)

	// GetRole use action=read, member can read
	roles2, err := service.GetRole(stage.Member, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles2, 2)

	roles3, err := service.GetRole(stage.Superadmin, caskin.GetSuperadminDomain())
	assert.NoError(t, err)
	assert.Len(t, roles3, 0)
}

func TestServer_CreateRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	role1 := &example.Role{Name: "role_01"}
	assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(stage.Member, stage.Domain, role1))
	assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(stage.Admin, stage.Domain, role1))

	roles, _ := service.GetRole(stage.Member, stage.Domain)
	assert.Len(t, roles, 2)
	role1.ObjectID = roles[0].GetObjectID()
	assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(stage.Member, stage.Domain, role1))
	assert.NoError(t, service.CreateRole(stage.Admin, stage.Domain, role1))

	role2 := &example.Role{Name: "role_01"}
	assert.Equal(t, caskin.ErrAlreadyExists, service.CreateRole(stage.Member, stage.Domain, role2))
}

func TestServer_UpdateRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	role1 := &example.Role{ID: roles[0].GetID(), Name: "admin_01_new_name"}
	assert.Equal(t, caskin.ErrNoWritePermission, service.UpdateRole(stage.Admin, stage.Domain, role1))
	role1.ObjectID = roles[0].GetObjectID()
	assert.Equal(t, caskin.ErrNoWritePermission, service.UpdateRole(stage.Member, stage.Domain, role1))
	assert.NoError(t, service.UpdateRole(stage.Admin, stage.Domain, role1))

	role2 := &example.Role{Name: "member_01_new_name", ObjectID: objects[1].GetID()}
	assert.Equal(t, caskin.ErrEmptyID, service.UpdateRole(stage.Member, stage.Domain, role2))
	role2.ID = roles[1].GetID()
	assert.Equal(t, caskin.ErrInValidObjectType, service.UpdateRole(stage.Admin, stage.Domain, role2))
}

func TestServer_RecoverRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	role1 := &example.Role{ID: roles[1].GetID()}
	assert.NoError(t, service.DeleteRole(stage.Admin, stage.Domain, role1))

	role2 := &example.Role{ID: roles[1].GetID()}
	assert.Equal(t, caskin.ErrNoWritePermission, service.RecoverRole(stage.Member, stage.Domain, role2))
	assert.NoError(t, service.RecoverRole(stage.Admin, stage.Domain, role2))
	assert.Equal(t, caskin.ErrAlreadyExists, service.RecoverRole(stage.Admin, stage.Domain, role2))
}

func TestServer_DeleteRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	role1 := &example.Role{ID: roles[0].GetID(), ObjectID: roles[0].GetObjectID()}
	assert.NoError(t, service.DeleteRole(stage.Admin, stage.Domain, role1))
	assert.Equal(t, caskin.ErrNoWritePermission, service.RecoverRole(stage.Admin, stage.Domain, role1))
	assert.NoError(t, service.RecoverRole(stage.Superadmin, stage.Domain, role1))
}
