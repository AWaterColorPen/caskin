package web_feature_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"github.com/stretchr/testify/assert"
)

func TestExecutorBackend_Create(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	backend1 := &web_feature.Backend{Path: "api/test-1", Method: "GET"}
	assert.Equal(t, caskin.ErrProviderGet, executor.CreateBackend(backend1, &example.Object{}))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.CreateBackend(backend1, &example.Object{}))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateBackend(backend1, &example.Object{}))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.CreateBackend(backend1, &example.Object{}))

	backend2 := &web_feature.Backend{Path: "api/backend", Method: "GET"}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateBackend(backend2, &example.Object{}))
}

func TestExecutorBackend_Recover(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	backend1 := &web_feature.Backend{Path: "api/backend", Method: "GET"}
	assert.Equal(t, caskin.ErrProviderGet, executor.RecoverBackend(backend1))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.RecoverBackend(backend1))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverBackend(backend1))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.DeleteBackend(&example.Object{ID: backendStartID}))
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.RecoverBackend(backend1))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.RecoverBackend(backend1))

	backend2 := &web_feature.Backend{Path: "api/backend"}
	assert.Equal(t, caskin.ErrNotExists, executor.RecoverBackend(backend2))
}

func TestExecutorBackend_Delete(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	assert.Equal(t, caskin.ErrProviderGet, executor.DeleteBackend(&example.Object{ID: backendStartID}))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.DeleteBackend(&example.Object{ID: backendStartID}))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteBackend(&example.Object{ID: backendStartID}))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.DeleteBackend(&example.Object{ID: backendStartID}))

	assert.Equal(t, caskin.ErrNotExists, executor.DeleteBackend(&example.Object{ID: featureStartID}))
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteBackend(&example.Object{ID: backendStartID - 1}))
}

func TestExecutorBackend_Update(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	backend1 := &web_feature.Backend{Path: "api/backend", Method: "GET"}
	object1 := &example.Object{ID: backendStartID + 1}
	assert.Equal(t, caskin.ErrProviderGet, executor.UpdateBackend(backend1, object1))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.UpdateBackend(backend1, object1))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateBackend(backend1, object1))
	provider.User = stage.SuperadminUser
	assert.Error(t, executor.UpdateBackend(backend1, object1))
	object1.ID = backendStartID
	assert.NoError(t, executor.UpdateBackend(backend1, object1))

	backend2 := &web_feature.Backend{Path: "api/test", Method: "GET"}
	object2 := &example.Object{ID: featureStartID}
	assert.Equal(t, caskin.ErrCantChangeObjectType, executor.UpdateBackend(backend2, object2))
	object2.ID = backendStartID - 1
	assert.Equal(t, caskin.ErrCantChangeObjectType, executor.UpdateBackend(backend2, object2))
	object2.ID = backendStartID + 1
	assert.NoError(t, executor.UpdateBackend(backend2, object2))

	list1, err := executor.GetBackend()
	assert.NoError(t, err)
	assert.Len(t, list1, 8)
	assert.Equal(t, "api/test$$GET", list1[2].ObjectCustomizedData.GetName())
}

func TestExecutorBackend_Get(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	_, err = executor.GetBackend()
	assert.Equal(t, caskin.ErrProviderGet, err)
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	_, err = executor.GetBackend()
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, err)
	provider.Domain = stage.Options.GetSuperadminDomain()
	list1, err := executor.GetBackend()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)
	provider.User = stage.SuperadminUser
	list2, err := executor.GetBackend()
	assert.NoError(t, err)
	assert.Len(t, list2, 8)
	assert.Equal(t, "backend-root$$", list2[0].ObjectCustomizedData.GetName())
}
