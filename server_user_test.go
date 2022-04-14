package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestServer_CreateUser(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{Email: "member2@qq.com"}
	assert.NoError(t, service.CreateUser(user1))

	user2 := &example.User{Email: "member2@qq.com"}
	assert.Equal(t, caskin.ErrAlreadyExists, service.CreateUser(user2))

	user3 := &example.User{Email: "member2@qq.com"}
	assert.Equal(t, caskin.ErrEmptyID, service.DeleteUser(user3))
	user3.ID = user2.ID
	assert.NoError(t, service.DeleteUser(user3))

	user4 := &example.User{ID: user2.ID + 1}
	assert.Equal(t, caskin.ErrNotExists, service.DeleteUser(user4))
	assert.NoError(t, service.CreateUser(user4))
}

func TestServer_UpdateUser(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{ID: stage.Member.ID, Email: "member2@qq.com"}
	assert.NoError(t, service.UpdateUser(user1))

	user2 := &example.User{}
	assert.Equal(t, caskin.ErrEmptyID, service.UpdateUser(user2))

	user3 := &example.User{ID: 999}
	assert.Equal(t, caskin.ErrNotExists, service.UpdateUser(user3))
}

func TestServer_RecoverUser(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{}
	assert.Equal(t, caskin.ErrAlreadyExists, service.RecoverUser(user1))
	assert.NoError(t, service.DeleteUser(stage.Member))

	user2 := &example.User{Email: stage.Member.Email}
	assert.NoError(t, service.RecoverUser(user2))

	user3 := &example.User{ID: 999}
	assert.Equal(t, caskin.ErrNotExists, service.RecoverUser(user3))
}

func TestServer_DeleteUser(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	assert.NoError(t, service.DeleteUser(stage.Superadmin))
	list1, err := service.GetSuperadmin()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	assert.NoError(t, service.DeleteUser(stage.Member))
	list2, err := service.GetUserRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
}
