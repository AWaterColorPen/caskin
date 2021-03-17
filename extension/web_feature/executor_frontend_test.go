package web_feature_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"github.com/stretchr/testify/assert"
)

func TestExecutorFrontend_Create(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	frontend1 := &web_feature.Frontend{Key: "backend", Type: web_feature.FrontendTypeSubFunction}
	assert.Equal(t, caskin.ErrProviderGet, executor.CreateFrontend(frontend1))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.CreateFrontend(frontend1))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateFrontend(frontend1))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.CreateFrontend(frontend1))

	frontend2 := &web_feature.Frontend{Key: "backend", Type: web_feature.FrontendTypeMenu}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateFrontend(frontend2))
}

func TestExecutorFrontend_Recover(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	frontend1 := &web_feature.Frontend{Key: "backend", Type: web_feature.FrontendTypeMenu}
	assert.Equal(t, caskin.ErrProviderGet, executor.RecoverFrontend(frontend1))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.RecoverFrontend(frontend1))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverFrontend(frontend1))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.DeleteFrontend(&example.Object{ID: frontendStartID}))
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.RecoverFrontend(frontend1))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.RecoverFrontend(frontend1))

	frontend2 := &web_feature.Frontend{Key: "backend"}
	assert.Equal(t, caskin.ErrNotExists, executor.RecoverFrontend(frontend2))
}

func TestExecutorFrontend_Delete(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	assert.Equal(t, caskin.ErrProviderGet, executor.DeleteFrontend(&example.Object{ID: frontendStartID}))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.DeleteFrontend(&example.Object{ID: frontendStartID}))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteFrontend(&example.Object{ID: frontendStartID}))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.DeleteFrontend(&example.Object{ID: frontendStartID}))

	assert.Equal(t, caskin.ErrNotExists, executor.DeleteFrontend(&example.Object{ID: featureStartID}))
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteFrontend(&example.Object{ID: frontendStartID - 1}))
}

func TestExecutorFrontend_Update(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	frontend1 := &web_feature.Frontend{Key: "backend", Type: web_feature.FrontendTypeMenu}
	object1 := &example.Object{ID: frontendStartID + 1}
	assert.Equal(t, caskin.ErrProviderGet, executor.UpdateFrontend(frontend1, object1))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.UpdateFrontend(frontend1, object1))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateFrontend(frontend1, object1))
	provider.User = stage.SuperadminUser
	assert.Error(t, executor.UpdateFrontend(frontend1, object1))
	object1.ID = frontendStartID
	assert.NoError(t, executor.UpdateFrontend(frontend1, object1))

	frontend2 := &web_feature.Frontend{Key: "backend", Type: web_feature.FrontendTypeSubFunction}
	object2 := &example.Object{ID: featureStartID}
	assert.Equal(t, caskin.ErrCantChangeObjectType, executor.UpdateFrontend(frontend2, object2))
	object2.ID = frontendStartID - 1
	assert.Equal(t, caskin.ErrCantChangeObjectType, executor.UpdateFrontend(frontend2, object2))
	object2.ID = frontendStartID + 1
	assert.NoError(t, executor.UpdateFrontend(frontend2, object2))

	list1, err := executor.GetFrontend()
	assert.NoError(t, err)
	assert.Len(t, list1, 5)
	assert.Equal(t, "backend_sub_function", list1[2].ObjectCustomizedData.GetName())
}

func TestExecutorFrontend_Get(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	_, err = executor.GetFrontend()
	assert.Equal(t, caskin.ErrProviderGet, err)
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	_, err = executor.GetFrontend()
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, err)
	provider.Domain = stage.Options.GetSuperadminDomain()
	list1, err := executor.GetFrontend()
	assert.Len(t, list1, 0)
	provider.User = stage.SuperadminUser
	list2, err := executor.GetFrontend()
	assert.NoError(t, err)
	assert.Len(t, list2, 5)
	assert.Equal(t, "frontend-root_", list2[0].ObjectCustomizedData.GetName())
}