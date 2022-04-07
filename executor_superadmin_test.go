package caskin_test

import (
	"github.com/awatercolorpen/caskin/playground"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorSuperadmin_Add(t *testing.T) {
	stage, _ := playground.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	user1 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.SuperadminAdd(user1))
	user1.ID = stage.MemberUser.ID
	assert.Error(t, executor.SuperadminAdd(user1))
	assert.NoError(t, executor.SuperadminAdd(stage.MemberUser))

	list1, err := executor.SuperadminGet()
	assert.NoError(t, err)
	assert.Len(t, list1, 2)
}

func TestExecutorSuperadmin_Delete(t *testing.T) {
	stage, _ := playground.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	user1 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.SuperadminDelete(user1))
	user1.ID = stage.MemberUser.ID
	assert.Error(t, executor.SuperadminDelete(user1))

	// delete a no superadmin user, it will not return error
	assert.NoError(t, executor.SuperadminDelete(stage.MemberUser))
	list1, err := executor.SuperadminGet()
	assert.NoError(t, err)
	assert.Len(t, list1, 1)

	assert.NoError(t, executor.SuperadminDelete(stage.SuperadminUser))
	list2, err := executor.SuperadminGet()
	assert.NoError(t, err)
	assert.Len(t, list2, 0)
}
