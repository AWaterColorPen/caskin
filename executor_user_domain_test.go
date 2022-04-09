package caskin_test

import (
	"github.com/awatercolorpen/caskin/playground"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/stretchr/testify/assert"
)

func TestExecutorUserDomain_GetUserInDomain(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	
	assert.NoError(t, stage.AddSubAdmin())
	service := stage.Service

	provider.Domain = stage.Domain
	list1, err := executor.GetUserInDomain(stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list1, 3)
	list2, err := executor.GetUserInDomain(stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list2, 3)
	list3, err := executor.GetUserInDomain(stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list3, 3)
	list4, err := executor.GetUserInDomain(stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list4, 3)
}

func TestExecutorUserDomain_GetDomainByUser(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	
	assert.NoError(t, stage.AddSubAdmin())
	service := stage.Service

	list1, err := executor.GetDomainByUser(stage.AdminUser)
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	list2, err := executor.GetDomainByUser(stage.SubAdminUser)
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
	list3, err := executor.GetDomainByUser(stage.SuperadminUser)
	assert.NoError(t, err)
	assert.Len(t, list3, 2)
	list4, err := executor.GetDomainByUser(stage.MemberUser)
	assert.NoError(t, err)
	assert.Len(t, list4, 1)
}
