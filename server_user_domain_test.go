package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin/playground"

	"github.com/stretchr/testify/assert"
)

func TestServer_UserDomain_GetUserByDomain(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	list1, err := service.GetUserByDomain(stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list1, 2)
}

func TestServer_UserDomain_GetDomainByUser(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	service := stage.Service

	list1, err := service.GetDomainByUser(stage.Admin)
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	// list2, err := service.GetDomainByUser(stage.SubAdminUser)
	// assert.NoError(t, err)
	// assert.Len(t, list2, 1)
	list3, err := service.GetDomainByUser(stage.Superadmin)
	assert.NoError(t, err)
	assert.Len(t, list3, 2)
	list4, err := service.GetDomainByUser(stage.Member)
	assert.NoError(t, err)
	assert.Len(t, list4, 1)
}
