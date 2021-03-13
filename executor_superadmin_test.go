package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorSuperadmin_Add(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	user1 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.AddSuperadminUser(user1))
	user1.ID = stage.MemberUser.ID
	assert.Error(t, executor.AddSuperadminUser(user1))
	assert.NoError(t, executor.AddSuperadminUser(stage.MemberUser))

	list1, err := executor.GetAllSuperadminUser()
	assert.NoError(t, err)
	assert.Len(t, list1, 2)
}

func TestExecutorSuperadmin_Delete(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	user1 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteSuperadminUser(user1))
	user1.ID = stage.MemberUser.ID
	assert.Error(t, executor.DeleteSuperadminUser(user1))

	// delete a no superadmin user, it will not return error
	assert.NoError(t, executor.DeleteSuperadminUser(stage.MemberUser))
	list1, err := executor.GetAllSuperadminUser()
	assert.NoError(t, err)
	assert.Len(t, list1, 1)

	assert.NoError(t, executor.DeleteSuperadminUser(stage.SuperadminUser))
	list2, err := executor.GetAllSuperadminUser()
	assert.NoError(t, err)
	assert.Len(t, list2, 0)
}

func TestExecutorSuperadmin_NoSuperadmin(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.NoSuperadmin())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	user1 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.NoError(t, executor.CreateUser(user1))
	assert.Equal(t, caskin.ErrSuperAdminIsNoEnabled, executor.AddSuperadminUser(user1))
	assert.Equal(t, caskin.ErrSuperAdminIsNoEnabled, executor.DeleteSuperadminUser(stage.AdminUser))
	_, err := executor.GetAllSuperadminUser()
	assert.Equal(t, caskin.ErrSuperAdminIsNoEnabled, err)
}
