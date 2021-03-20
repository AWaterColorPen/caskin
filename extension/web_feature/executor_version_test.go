package web_feature_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"github.com/stretchr/testify/assert"
)

func TestExecutorVersion_BuildVersion(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	assert.Equal(t, caskin.ErrProviderGet, executor.BuildVersion())
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.BuildVersion())
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrIsNotSuperAdmin, executor.BuildVersion())
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.BuildVersion())

	list1, err := executor.GetVersion()
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	assert.Equal(t, "aa2a319fdc891c2d66f0e1a866ea66ef01ee22e74a2f66ff58dae879e98ecdf8", list1[0].SHA256)

	assert.Error(t, executor.BuildVersion())
}

func TestExecutorVersion_SyncVersionToAllDomain(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	assert.Equal(t, caskin.ErrProviderGet, executor.SyncVersionToAllDomain(nil))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.SyncVersionToAllDomain(nil))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrIsNotSuperAdmin, executor.SyncVersionToAllDomain(nil))
	provider.User = stage.SuperadminUser
	assert.Error(t, executor.SyncVersionToAllDomain(nil))

	version := &web_feature.WebFeatureVersion{
		SHA256: "sha256",
	}
	assert.Error(t, executor.SyncVersionToAllDomain(version))

	assert.NoError(t, executor.BuildVersion())
	versions, err := executor.GetVersion()
	assert.NoError(t, err)
	assert.Len(t, versions, 1)
	assert.NoError(t, executor.SyncVersionToAllDomain(versions[0]))
	// test twice
	assert.NoError(t, executor.SyncVersionToAllDomain(versions[0]))
}

func TestExecutorVersion_SyncLatestVersionToAllDomain(t *testing.T) {
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
	// test twice
	assert.NoError(t, executor.SyncLatestVersionToAllDomain())

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	list1, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	assert.NoError(t, reinitializeDomainWithWebFeature(stage))
	list2, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, list2, 5)
	list3, err := executor.NormalDomainGetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, list3, 1)

	provider.User = stage.MemberUser
	list4, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, list4, 0)
	list5, err := executor.NormalDomainGetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, list5, 0)

	// TODO change feature relation, and test re sync will delete some relation
}
