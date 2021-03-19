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
	assert.Error(t, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend"}))

	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoBackendAPIPermission, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend", Method: "GET"}))

	// TODO modify member role 's policy
}
