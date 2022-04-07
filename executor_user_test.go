package caskin_test

import (
	"github.com/awatercolorpen/caskin/playground"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorUser_GeneralCreate(t *testing.T) {
	stage, _ := playground.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	user1 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.NoError(t, executor.UserCreate(user1))

	user2 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.UserCreate(user2))

	user3 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.UserDelete(user3))
	user3.ID = user2.ID
	assert.NoError(t, executor.UserDelete(user3))

	user4 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.UserDelete(user4))
	assert.NoError(t, executor.UserCreate(user4))
}

func TestExecutorUser_GeneralUpdate(t *testing.T) {
	stage, _ := playground.NewStageWithSqlitePath(t.TempDir())
	executor := stage.Caskin.GetExecutor(caskin.NewCachedProvider(nil, nil))

	user1 := &example.User{
		ID:          stage.MemberUser.ID,
		PhoneNumber: stage.MemberUser.PhoneNumber,
		Email:       "member2@qq.com",
	}
	assert.NoError(t, executor.UserUpdate(user1))

	user2 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.UserUpdate(user2))

	user3 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.UserUpdate(user3))
}

func TestExecutorUser_GeneralRecover(t *testing.T) {
	stage, _ := playground.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	user1 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.UserRecover(user1))
	assert.NoError(t, executor.UserDelete(stage.MemberUser))

	user2 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.NoError(t, executor.UserRecover(user2))

	user3 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.UserRecover(user3))
}

func TestExecutorUser_GeneralDelete(t *testing.T) {
	stage, _ := playground.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	domain := &example.Domain{Name: "domain_02"}
	assert.NoError(t, executor.DomainCreate(domain))
	assert.NoError(t, executor.DomainInitialize(domain))

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

	assert.NoError(t, executor.UserDelete(stage.SuperadminUser))
	list1, err := executor.SuperadminGet()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	assert.NoError(t, executor.UserDelete(stage.MemberUser))
	list2, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
}
