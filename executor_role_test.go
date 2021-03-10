package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutorRole(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{
		User:   stage.SuperadminUser,
		Domain: stage.Domain,
	}
	executor := stage.Caskin.GetExecutor(provider)

	usersForRoleList, err := executor.GetAllUsersForRole()
	assert.NoError(t, err)
	assert.Len(t, usersForRoleList, 2)
	assert.Len(t, usersForRoleList[0].Users, 1)
	//bytes, _ := json.Marshal(usersForRoleList)
	//fmt.Println(string(bytes))

	roles, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	usersForRole := &caskin.UsersForRole{
		Role: roles[0],
		Users: caskin.Users{},
	}

	assert.NoError(t, executor.ModifyUsersForRole(usersForRole))
}

func TestExecutorRole_GeneralCreate(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	role1 := &example.Role{
		Name: "sub_admin",
	}
	assert.Equal(t, caskin.ErrProviderGet, executor.CreateRole(role1))

	provider.Domain = stage.Domain
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateRole(role1))
	role1.ObjectID = 2
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateRole(role1))
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateRole(role1))
	// user2 := &example.User{
	// 	PhoneNumber: "12345678904",
	// 	Email:       "member2@qq.com",
	// }
	// assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateUser(user2))
	//
	// user3 := &example.User{
	// 	PhoneNumber: "12345678904",
	// 	Email:       "member2@qq.com",
	// }
	// assert.Equal(t, caskin.ErrEmptyID, executor.DeleteUser(user3))
	// user3.ID = user2.ID
	// assert.NoError(t, executor.DeleteUser(user3))
	//
	// user4 := &example.User{ID: 5}
	// assert.Equal(t, caskin.ErrNotExists, executor.DeleteUser(user4))
	// assert.NoError(t, executor.CreateUser(user4))
}

func TestExecutorRole_GeneralUpdate(t *testing.T) {
	stage, _ := newStage(t)
	executor := stage.Caskin.GetExecutor(&example.Provider{})

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

func TestExecutorRole_GeneralRecover(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{}
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
	assert.Error(t, executor.RecoverUser(user3))
}

func TestExecutorRole_GeneralDelete(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	domain := &example.Domain{Name: "domain_02"}
	assert.NoError(t, executor.CreateDomain(domain))
	assert.NoError(t, executor.ReInitializeDomain(domain))

	provider.Domain = domain
	provider.User = stage.SuperadminUser
	roles, err := executor.GetRoles()
	assert.NoError(t, err)

	for _, v := range []*caskin.RolesForUser{
		{User: stage.MemberUser, Roles: []caskin.Role{roles[0]}},
		{User: stage.AdminUser, Roles: []caskin.Role{roles[1]}},
	} {
		assert.NoError(t, executor.ModifyRolesForUser(v))
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

