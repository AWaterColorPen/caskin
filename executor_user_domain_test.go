package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutorUserDomain_GetUserInDomain(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	assert.NoError(t, stage.AddSubAdmin())
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	list1, err := executor.GetUserInDomain()
	assert.NoError(t, err)
	assert.Len(t, list1, 3)
	provider.User = stage.SubAdminUser
	list2, err := executor.GetUserInDomain()
	assert.NoError(t, err)
	assert.Len(t, list2, 3)
	provider.User = stage.SuperadminUser
	list3, err := executor.GetUserInDomain()
	assert.NoError(t, err)
	assert.Len(t, list3, 3)
	provider.User = stage.MemberUser
	list4, err := executor.GetUserInDomain()
	assert.NoError(t, err)
	assert.Len(t, list4, 3)
}

func TestExecutorUserDomain_GetDomainByUser(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	assert.NoError(t, stage.AddSubAdmin())
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	list1, err := executor.GetDomainByUser()
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	provider.User = stage.SubAdminUser
	list2, err := executor.GetDomainByUser()
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
	provider.User = stage.SuperadminUser
	list3, err := executor.GetDomainByUser()
	assert.NoError(t, err)
	assert.Len(t, list3, 2)
	provider.User = stage.MemberUser
	list4, err := executor.GetDomainByUser()
	assert.NoError(t, err)
	assert.Len(t, list4, 1)
}
