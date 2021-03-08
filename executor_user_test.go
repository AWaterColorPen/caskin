package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"testing"

	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorUser_GeneralCreateAndDelete(t *testing.T) {
	stage, _ := getStage(t)
	provider := &example.Provider{
		User: stage.MemberUser,
		Domain: stage.Domain,
	}
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
	stage, _ := getStage(t)
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

func TestExecutorUser_GeneralRecover(t *testing.T) {
	stage, _ := getStage(t)
	provider := &example.Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	user1 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverUser(user1))
	assert.Equal(t, caskin.ErrProviderGet, executor.DeleteUser(stage.MemberUser))

	provider.User = stage.MemberUser
	provider.Domain = stage.Domain
	assert.NoError(t, executor.DeleteUser(stage.MemberUser))

	user2 := &example.User{
		PhoneNumber: stage.MemberUser.PhoneNumber,
	}
	assert.NoError(t, executor.RecoverUser(user2))

	user3 := &example.User{ID: 5}
	assert.Error(t, executor.RecoverUser(user3))
}
