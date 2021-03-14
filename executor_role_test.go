package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutorRole_GeneralCreate(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	role1 := &example.Role{
		Name: "sub_admin_1",
	}
	assert.Equal(t, caskin.ErrProviderGet, executor.CreateRole(role1))

	provider.Domain = stage.Domain
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateRole(role1))
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateRole(role1))
	role1.ObjectID = 3
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrInValidObjectType, executor.CreateRole(role1))
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrInValidObjectType, executor.CreateRole(role1))
	role1.ObjectID = 2
	assert.NoError(t, executor.CreateRole(role1))

	role2 := &example.Role{
		Name:     "sub_admin_1",
		ObjectID: 2,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateRole(role2))

	role3 := &example.Role{
		Name:     "sub_admin_1",
		ObjectID: 2,
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteRole(role3))
	role3.ID = role1.ID
	assert.NoError(t, executor.DeleteRole(role3))

	role4 := &example.Role{ID: 5, ObjectID: 2}
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteRole(role4))
	assert.NoError(t, executor.CreateRole(role4))
}

func TestExecutorRole_CreateSubNode(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	role1 := &example.Role{
		Name:     "sub_member_1",
		ObjectID: 2,
		ParentID: 3,
	}
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateRole(role1))

	role1.ParentID = 2
	assert.NoError(t, executor.CreateRole(role1))

	role1.ObjectID = 2
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateRole(role1))
	provider.User = stage.AdminUser
	assert.NoError(t, executor.CreateRole(role1))

	role2 := &example.Role{
		Name: "sub_admin_1",
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateRole(role2))

	role3 := &example.Role{
		Name: "sub_admin_1",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteRole(role3))
	role3.ID = role1.ID
	assert.NoError(t, executor.DeleteRole(role3))

	role4 := &example.Role{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteRole(role4))
	assert.NoError(t, executor.CreateRole(role4))
}

func TestExecutorRole_GeneralUpdate(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	provider.User = stage.AdminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	role1 := &example.Role{
		ID:   2,
		Name: "member_momoda",
	}
	assert.NoError(t, executor.UpdateRole(role1))

	role2 := &example.Role{
		Name: "member_momoda",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.UpdateRole(role2))

	role3 := &example.Role{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.UpdateRole(role3))
}

func TestExecutorRole_GeneralRecover(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	role1 := &example.Role{
		Name: "member",
	}
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverRole(role1))
	assert.NoError(t, executor.DeleteRole(role1))

	role2 := &example.Role{
		Name: "member",
	}
	assert.NoError(t, executor.RecoverRole(role2))

	role3 := &example.Role{ID: 5}
	assert.Error(t, executor.RecoverRole(role3))
}

func TestExecutorRole_GeneralDelete(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	role1 := &example.Role{Name: "member"}
	provider.Domain = stage.Domain
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteRole(role1))
	role1.ID = 2
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteRole(role1))

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.NoError(t, executor.DeleteRole(role1))

	roles1, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles1, 1)

	list1, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
}
