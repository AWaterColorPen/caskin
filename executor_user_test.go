package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin/playground"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestServer_UserCreate(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{
		Email: "member2@qq.com",
	}
	assert.NoError(t, service.CreateUser(user1))

	user2 := &example.User{
		Email: "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrAlreadyExists, service.CreateUser(user2))

	user3 := &example.User{
		Email: "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrEmptyID, service.DeleteUser(user3))
	user3.ID = user2.ID
	assert.NoError(t, service.DeleteUser(user3))

	user4 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, service.DeleteUser(user4))
	assert.NoError(t, service.CreateUser(user4))
}

func TestServer_UserUpdate(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{
		ID:    stage.Member.ID,
		Email: "member2@qq.com",
	}
	assert.NoError(t, service.UpdateUser(user1))

	user2 := &example.User{
		// PhoneNumber: stage.Member.PhoneNumber,
	}
	assert.Equal(t, caskin.ErrEmptyID, service.UpdateUser(user2))

	user3 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, service.UpdateUser(user3))
}

func TestServer_UserRecover(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{
		// PhoneNumber: stage.Member.PhoneNumber,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, service.RecoverUser(user1))
	assert.NoError(t, service.DeleteUser(stage.Member))

	user2 := &example.User{
		// PhoneNumber: stage.Member.PhoneNumber,
	}
	assert.NoError(t, service.RecoverUser(user2))

	user3 := &example.User{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, service.RecoverUser(user3))
}

func TestServer_UserDelete(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	domain := &example.Domain{Name: "domain_02"}
	assert.NoError(t, service.CreateDomain(domain))
	assert.NoError(t, service.ResetDomain(domain))

	roles, err := service.GetRole(stage.Superadmin, domain)
	assert.NoError(t, err)

	for k, v := range map[caskin.Role][]*caskin.UserRolePair{
		roles[0]: {{User: stage.Member, Role: roles[0]}},
		roles[1]: {{User: stage.Admin, Role: roles[1]}},
	} {
		assert.NoError(t, service.ModifyUserRolePerRole(stage.Superadmin, domain, k, v))
	}

	assert.NoError(t, service.DeleteUser(stage.Superadmin))
	list1, err := service.GetSuperadmin()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	assert.NoError(t, service.DeleteUser(stage.Member))
	list2, err := service.GetUserRole(stage.Superadmin, domain)
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
}
