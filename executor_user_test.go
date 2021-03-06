package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorUser_GeneralCreate(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	user1 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.NoError(t, executor.CreateUser(user1))

	user2 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateUser(user2))

	user3 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteUser(user3))
	user3.ID = user2.ID
	assert.NoError(t, executor.DeleteUser(user3))

	user4 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteUser(user4))
	assert.NoError(t, executor.CreateUser(user4))
}

func TestExecutorUser_GeneralUpdate(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	executor := stage.Caskin.GetExecutor(caskin.NewCachedProvider(nil, nil))

	user1 := &example.User{
		ID:          stage.MemberUser.ID,
		PhoneNumber: stage.MemberUser.PhoneNumber,
		Email:       "member2@qq.com",
	}
	assert.NoError(t, executor.UpdateUser(user1))

	user2 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.UpdateUser(user2))

	user3 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.UpdateUser(user3))
}

func TestExecutorUser_GeneralRecover(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	user1 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverUser(user1))
	assert.NoError(t, executor.DeleteUser(stage.MemberUser))

	user2 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.NoError(t, executor.RecoverUser(user2))

	user3 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.RecoverUser(user3))
}

func TestExecutorUser_GeneralDelete(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	domain := &example.Domain{Name: "domain_02"}
	assert.NoError(t, executor.CreateDomain(domain))
	assert.NoError(t, executor.ReInitializeDomain(domain))

	provider.Domain = domain
	provider.User = stage.SuperadminUser
	roles, err := executor.GetRoles()
	assert.NoError(t, err)

	for k, v := range map[caskin.Role][]*caskin.UserRolePair{
		roles[0]: {{User: stage.MemberUser, Role: roles[0]}},
		roles[1]: {{User: stage.AdminUser, Role: roles[1]}},
	} {
		assert.NoError(t, executor.ModifyUserRolePairPerRole(k, v))
	}

	assert.NoError(t, executor.DeleteUser(stage.SuperadminUser))
	list1, err := executor.GetAllSuperadminUser()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	assert.NoError(t, executor.DeleteUser(stage.MemberUser))
	list2, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
}
