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
	assert.Equal(t, caskin.ErrEmptyID, executor.SuperadminUserAdd(user1))
	user1.ID = stage.MemberUser.ID
	assert.Error(t, executor.SuperadminUserAdd(user1))
	assert.NoError(t, executor.SuperadminUserAdd(stage.MemberUser))

	list1, err := executor.SuperadminUserGet()
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
	assert.Equal(t, caskin.ErrEmptyID, executor.SuperadminUserDelete(user1))
	user1.ID = stage.MemberUser.ID
	assert.Error(t, executor.SuperadminUserDelete(user1))

	// delete a no superadmin user, it will not return error
	assert.NoError(t, executor.SuperadminUserDelete(stage.MemberUser))
	list1, err := executor.SuperadminUserGet()
	assert.NoError(t, err)
	assert.Len(t, list1, 1)

	assert.NoError(t, executor.SuperadminUserDelete(stage.SuperadminUser))
	list2, err := executor.SuperadminUserGet()
	assert.NoError(t, err)
	assert.Len(t, list2, 0)
}
