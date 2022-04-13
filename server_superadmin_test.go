package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin/playground"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestServer_AddSuperadmin(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{ID: 999, Email: "member2@qq.com"}
	assert.Equal(t, caskin.ErrNotExists, service.AddSuperadmin(user1))
	user1.ID = stage.Member.ID
	assert.Error(t, service.AddSuperadmin(user1))
	assert.NoError(t, service.AddSuperadmin(stage.Member))

	list1, err := service.GetSuperadmin()
	assert.NoError(t, err)
	assert.Len(t, list1, 2)
}

func TestServer_DeleteSuperadmin(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{Email: "member2@qq.com"}
	assert.Equal(t, caskin.ErrEmptyID, service.DeleteSuperadmin(user1))
	user1.ID = stage.Member.ID
	assert.Error(t, service.DeleteSuperadmin(user1))

	// delete a no superadmin user, it will not return error
	assert.NoError(t, service.DeleteSuperadmin(stage.Member))
	list1, err := service.GetSuperadmin()
	assert.NoError(t, err)
	assert.Len(t, list1, 1)

	assert.NoError(t, service.DeleteSuperadmin(stage.Superadmin))
	list2, err := service.GetSuperadmin()
	assert.NoError(t, err)
	assert.Len(t, list2, 0)
}
