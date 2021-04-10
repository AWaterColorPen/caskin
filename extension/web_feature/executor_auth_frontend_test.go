package web_feature_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/stretchr/testify/assert"
)

func TestExecutorAuthFrontend_CaskinStruct(t *testing.T) {
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	assert.Error(t, executor.SyncLatestVersionToAllDomain())

	assert.NoError(t, executor.BuildVersion())
	assert.NoError(t, executor.SyncLatestVersionToAllDomain())
	assert.NoError(t, reinitializeDomainWithWebFeature(stage, w.GetRoot()))

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser

	c1, err := executor.AuthFrontendCaskinStruct("abc")
	assert.NoError(t, err)
	assert.Len(t, c1.M, 324)
	assert.Len(t, c1.G, 0)
	assert.Len(t, c1.G2, 0)
	assert.Len(t, c1.P, 4)
	assert.Equal(t, "abc", c1.P[0][1])
	assert.Equal(t, string(caskin.Read), c1.P[0][4])

	provider.User = stage.MemberUser
	c2, err := executor.AuthFrontendCaskinStruct("abc")
	assert.NoError(t, err)
	assert.Len(t, c2.P, 0)

	provider.User = stage.SuperadminUser
	c3, err := executor.AuthFrontendCaskinStruct("abc")
	assert.NoError(t, err)
	assert.Len(t, c3.P, 5)

	provider.Domain = stage.Options.GetSuperadminDomain()
	c4, err := executor.AuthFrontendCaskinStruct("abc")
	assert.NoError(t, err)
	assert.Len(t, c1.M, 324)
	assert.Len(t, c4.G, 1)
	assert.Len(t, c4.G2, 0)
	assert.Len(t, c4.P, 0)
	assert.Equal(t, "abc", c4.G[0][1])
	assert.Equal(t, caskin.SuperadminRole, c4.G[0][2])
	assert.Equal(t, caskin.SuperadminDomain, c4.G[0][3])

	provider.User = stage.AdminUser
	c5, err := executor.AuthFrontendCaskinStruct("abc")
	assert.NoError(t, err)
	assert.Len(t, c5.G, 0)
}
