package caskin_test

import (
	"github.com/awatercolorpen/caskin/playground"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_UserDomain_GetUserInDomain(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())

	assert.NoError(t, stage.AddSubAdmin())
	service := stage.Service

	list1, err := service.UserByDomainGet(stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list1, 3)
	list2, err := service.UserByDomainGet(stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list2, 3)
	list3, err := service.UserByDomainGet(stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list3, 3)
	list4, err := service.UserByDomainGet(stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list4, 3)
}

func TestServer_UserDomain_GetDomainByUser(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	service := stage.Service

	list1, err := service.DomainByUserGet(stage.AdminUser)
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	list2, err := service.DomainByUserGet(stage.SubAdminUser)
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
	list3, err := service.DomainByUserGet(stage.SuperadminUser)
	assert.NoError(t, err)
	assert.Len(t, list3, 2)
	list4, err := service.DomainByUserGet(stage.MemberUser)
	assert.NoError(t, err)
	assert.Len(t, list4, 1)
}
