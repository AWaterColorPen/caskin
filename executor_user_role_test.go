package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"testing"

	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorUserRole_GetUserRolePair(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := &Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	list1, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list1, 3)
	roles, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles, 3)

	pairs := []*caskin.UserRolePair{
		{stage.MemberUser, roles[1]},
		{stage.SubAdminUser, roles[1]},
	}
	assert.NoError(t, executor.ModifyUserRolePairPerRole(roles[1], pairs))

	list3, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list3, 4)
}

func TestExecutorUserRole_GetUserRolePairSubAdmin(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := &Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.SubAdminUser

	list, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestExecutorUserRole_GetUserRolePairByRole(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := &Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	role1 := &example.Role{ID: 2, Name: "xxx"}
	_, err := executor.GetUserRolePairByRole(role1)
	assert.Equal(t, caskin.ErrNotExists, err)

	role1.Name = "member"
	pair1, err := executor.GetUserRolePairByRole(role1)
	assert.NoError(t, err)
	assert.Len(t, pair1, 1)

	role2 := &example.Role{Name: "admin"}
	_, err = executor.GetUserRolePairByRole(role2)
	assert.Equal(t, caskin.ErrEmptyID, err)

	provider.User = stage.MemberUser
	_, err = executor.GetUserRolePairByRole(role1)
	assert.Equal(t, caskin.ErrNoReadPermission, err)

	provider.User = stage.SubAdminUser
	role3 := &example.Role{ID: 3}
	pair2, err := executor.GetUserRolePairByRole(role3)
	assert.NoError(t, err)
	assert.Len(t, pair2, 1)

	_, err = executor.GetUserRolePairByRole(role1)
	assert.Equal(t, caskin.ErrNoReadPermission, err)
}

func TestExecutorUserRole_GetUserRolePairByUser(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := &Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	user1 := &example.User{ID: 4, PhoneNumber: "xxx"}
	_, err := executor.GetUserRolePairByUser(user1)
	assert.Equal(t, caskin.ErrNotExists, err)

	user1.PhoneNumber = stage.SubAdminUser.PhoneNumber
	pair1, err := executor.GetUserRolePairByUser(user1)
	assert.NoError(t, err)
	assert.Len(t, pair1, 1)

	user2 := &example.User{PhoneNumber: stage.SubAdminUser.PhoneNumber}
	_, err = executor.GetUserRolePairByUser(user2)
	assert.Equal(t, caskin.ErrEmptyID, err)

	provider.User = stage.MemberUser
	pair2, err := executor.GetUserRolePairByUser(user1)
	assert.NoError(t, err)
	assert.Len(t, pair2, 0)

	provider.User = stage.AdminUser
	roles, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles, 3)
	pairs := []*caskin.UserRolePair{
		{stage.MemberUser, roles[1]},
		{stage.SubAdminUser, roles[1]},
	}
	assert.NoError(t, executor.ModifyUserRolePairPerRole(roles[1], pairs))

	provider.User = stage.SubAdminUser
	user3 := &example.User{ID: 4}
	pair3, err := executor.GetUserRolePairByUser(user3)
	assert.NoError(t, err)
	assert.Len(t, pair3, 1)
}

func TestExecutorUserRole_ModifyUserRolePairPerRole(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := &Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	objects0, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects0, 5)
	roles0, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles0, 3)

	provider.User = stage.SubAdminUser
	roles1, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles1, 1)
	objects1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects1, 2)

	pair1 := []*caskin.UserRolePair{
		{stage.MemberUser, roles0[0]},
		{stage.SubAdminUser, roles0[0]},
	}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.ModifyUserRolePairPerRole(roles0[0], pair1))

	pair2 := []*caskin.UserRolePair{
		{stage.MemberUser, roles1[0]},
		{stage.SubAdminUser, roles1[0]},
	}
	assert.NoError(t, executor.ModifyUserRolePairPerRole(roles1[0], pair2))

	provider.User = stage.MemberUser
	list1, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list1, 2)

	provider.User = stage.AdminUser
	pair3 := []*caskin.UserRolePair{
		{stage.MemberUser, roles1[0]},
		{stage.SubAdminUser, roles1[0]},
		{stage.SubAdminUser, roles0[0]},
	}

	assert.Equal(t, caskin.ErrInputPairArrayNotBelongSameRole, executor.ModifyUserRolePairPerRole(roles1[0], pair3))
}

func TestExecutorUserRole_ModifyUserRolePairPerUser(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := &Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	objects0, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects0, 5)
	roles0, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles0, 3)

	provider.User = stage.SubAdminUser
	roles1, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles1, 1)
	objects1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects1, 2)

	pair1 := []*caskin.UserRolePair{
		{stage.MemberUser, roles0[0]},
	}
	assert.NoError(t, executor.ModifyUserRolePairPerUser(stage.MemberUser, pair1))
	provider.User = stage.AdminUser
	list1, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list1, 3)

	provider.User = stage.SubAdminUser
	pair2 := []*caskin.UserRolePair{
		{stage.MemberUser, roles1[0]},
	}
	assert.NoError(t, executor.ModifyUserRolePairPerUser(stage.MemberUser, pair2))
	list2, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list2, 2)

	provider.User = stage.AdminUser
	pair3 := []*caskin.UserRolePair{
		{stage.SubAdminUser, roles1[0]},
		{stage.MemberUser, roles1[0]},
	}
	assert.Equal(t, caskin.ErrInputPairArrayNotBelongSameUser, executor.ModifyUserRolePairPerUser(stage.SubAdminUser, pair3))
}
