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
	assert.Equal(t, "3d421dcae2b1033717353d208b6bc49f62a118eb6452e42d5098f2631a446d99", list1[0].SHA256)
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

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	list1, err := executor.GetFeature()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	assert.NoError(t, reinitializeDomainWithWebFeature(stage))

	list2, err := executor.GetFeature()
	assert.NoError(t, err)
	assert.Len(t, list2, 0)
}
