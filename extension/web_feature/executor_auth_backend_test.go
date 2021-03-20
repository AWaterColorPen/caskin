package web_feature_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"github.com/stretchr/testify/assert"
)

func TestExecutorAuthBackend_Enforce(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
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
	assert.NoError(t, reinitializeDomainWithWebFeature(stage))

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.NoError(t, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend", Method: "GET"}))
	assert.Equal(t, caskin.ErrNoBackendAPIPermission, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend"}))

	provider.User = stage.SubAdminUser
	assert.Equal(t, caskin.ErrNoBackendAPIPermission, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend", Method: "GET"}))

	executor1 := stage.Caskin.GetExecutor(provider)
	provider.User = stage.AdminUser
	roles, err := executor1.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles, 3)

	input := []*caskin.Policy{
		{roles[2], &example.Object{ID: featureStartID}, stage.Domain, caskin.Read},
	}
	assert.NoError(t, executor.NormalDomainModifyPolicyListPerRole(roles[2], input))
	provider.User = stage.SubAdminUser
	list, err := executor.NormalDomainGetFeatureObject()
	assert.Len(t, list, 1)
	provider.User = stage.SubAdminUser
	assert.NoError(t, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend", Method: "GET"}))
	assert.Equal(t, caskin.ErrNoBackendAPIPermission, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/frontend", Method: "GET"}))

	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrNoBackendAPIPermission, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend", Method: "GET"}))
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrNoBackendAPIPermission, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend", Method: "GET"}))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/test-1", Method: "GET"}))
}
