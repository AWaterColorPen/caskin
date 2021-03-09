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
