package caskin_test

import (
	"github.com/awatercolorpen/caskin/playground"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestServer_UserCreate(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.NoError(t, service.UserCreate(user1))

	user2 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrAlreadyExists, service.UserCreate(user2))

	user3 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrEmptyID, service.UserDelete(user3))
	user3.ID = user2.ID
	assert.NoError(t, service.UserDelete(user3))

	user4 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, service.UserDelete(user4))
	assert.NoError(t, service.UserCreate(user4))
}

func TestServer_UserUpdate(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{
		ID:          stage.MemberUser.ID,
		PhoneNumber: stage.MemberUser.PhoneNumber,
		Email:       "member2@qq.com",
	}
	assert.NoError(t, service.UserUpdate(user1))

	user2 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.Equal(t, caskin.ErrEmptyID, service.UserUpdate(user2))

	user3 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, service.UserUpdate(user3))
}

func TestServer_UserRecover(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, service.UserRecover(user1))
	assert.NoError(t, service.UserDelete(stage.MemberUser))

	user2 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.NoError(t, service.UserRecover(user2))

	user3 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, service.UserRecover(user3))
}

func TestServer_UserDelete(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	domain := &example.Domain{Name: "domain_02"}
	assert.NoError(t, service.DomainCreate(domain))
	assert.NoError(t, service.DomainReset(domain))

	roles, err := service.RoleGet(stage.SuperadminUser, domain)
	assert.NoError(t, err)

	for k, v := range map[caskin.Role][]*caskin.UserRolePair{
		roles[0]: {{User: stage.MemberUser, Role: roles[0]}},
		roles[1]: {{User: stage.AdminUser, Role: roles[1]}},
	} {
		assert.NoError(t, service.UserRolePerRoleModify(stage.SuperadminUser, domain, k, v))
	}

	assert.NoError(t, service.UserDelete(stage.SuperadminUser))
	list1, err := service.SuperadminGet()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	assert.NoError(t, service.UserDelete(stage.MemberUser))
	list2, err := service.UserRoleGet(stage.SuperadminUser, domain)
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
}
