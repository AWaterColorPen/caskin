package caskin_test

import (
	"github.com/awatercolorpen/caskin/playground"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestServer_Superadmin_Add(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrEmptyID, service.SuperadminAdd(user1))
	user1.ID = stage.MemberUser.ID
	assert.Error(t, service.SuperadminAdd(user1))
	assert.NoError(t, service.SuperadminAdd(stage.MemberUser))

	list1, err := service.SuperadminGet()
	assert.NoError(t, err)
	assert.Len(t, list1, 2)
}

func TestServer_Superadmin_Delete(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	user1 := &example.User{
		PhoneNumber: "12345678904",
		Email:       "member2@qq.com",
	}
	assert.Equal(t, caskin.ErrEmptyID, service.SuperadminDelete(user1))
	user1.ID = stage.MemberUser.ID
	assert.Error(t, service.SuperadminDelete(user1))

	// delete a no superadmin user, it will not return error
	assert.NoError(t, service.SuperadminDelete(stage.MemberUser))
	list1, err := service.SuperadminGet()
	assert.NoError(t, err)
	assert.Len(t, list1, 1)

	assert.NoError(t, service.SuperadminDelete(stage.SuperadminUser))
	list2, err := service.SuperadminGet()
	assert.NoError(t, err)
	assert.Len(t, list2, 0)
}
