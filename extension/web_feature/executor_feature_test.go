package web_feature_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"github.com/stretchr/testify/assert"
)

func TestExecutorFeature_Create(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	feature1 := &web_feature.Feature{Name: "new-feature"}
	assert.Equal(t, caskin.ErrProviderGet, executor.CreateFeature(feature1, &example.Object{}))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.CreateFeature(feature1, &example.Object{}))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateFeature(feature1, &example.Object{}))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.CreateFeature(feature1, &example.Object{}))

	feature2 := &web_feature.Feature{Name: "feature"}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateFeature(feature2, &example.Object{}))
	feature3 := &web_feature.Feature{Name: "new-feature-2"}
	assert.Equal(t, caskin.ErrInValidObject, executor.CreateFeature(feature3, &example.Object{ID: 10}))
	assert.Equal(t, caskin.ErrInValidObjectType, executor.CreateFeature(feature3, &example.Object{ParentID: frontendStartID}))
	assert.NoError(t, executor.CreateFeature(feature3, &example.Object{ParentID: featureStartID}))
}

func TestExecutorFeature_Recover(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	feature1 := &web_feature.Feature{Name: "backend"}
	assert.Equal(t, caskin.ErrProviderGet, executor.RecoverFeature(feature1, &example.Object{}))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.RecoverFeature(feature1, &example.Object{}))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverFeature(feature1, &example.Object{}))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.DeleteFeature(&example.Object{ID: featureStartID}))
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.RecoverFeature(feature1, &example.Object{}))
	provider.User = stage.SuperadminUser
	assert.Equal(t, caskin.ErrInValidObject, executor.RecoverFeature(feature1, &example.Object{ID: 10}))
	assert.NoError(t, executor.RecoverFeature(feature1, &example.Object{}))

	feature2 := &web_feature.Feature{Name: "new-feature"}
	assert.Equal(t, caskin.ErrNotExists, executor.RecoverFeature(feature2, &example.Object{}))

}

func TestExecutorFeature_Delete(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	assert.Equal(t, caskin.ErrProviderGet, executor.DeleteFeature(&example.Object{ID: featureStartID}))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.DeleteFeature(&example.Object{ID: featureStartID}))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteFeature(&example.Object{ID: featureStartID}))
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.DeleteFeature(&example.Object{ID: featureStartID}))
	// it can't delete frontend and backend when it delete feature
	list1, err := executor.GetFrontend()
	assert.NoError(t, err)
	assert.Len(t, list1, 5)
	list2, err := executor.GetBackend()
	assert.NoError(t, err)
	assert.Len(t, list2, 8)

	assert.Equal(t, caskin.ErrNotExists, executor.DeleteFeature(&example.Object{ID: frontendStartID}))
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteFeature(&example.Object{ID: frontendStartID - 1}))
}

func TestExecutorFeature_Update(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	feature1 := &web_feature.Feature{Name: "backend", Group: "1"}
	object1 := &example.Object{ID: featureStartID + 1}
	assert.Equal(t, caskin.ErrProviderGet, executor.UpdateFeature(feature1, object1))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.UpdateFeature(feature1, object1))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateFeature(feature1, object1))
	provider.User = stage.SuperadminUser
	assert.Error(t, executor.UpdateFeature(feature1, object1))
	object1.ID = featureStartID
	object1.ObjectID = 1
	assert.NoError(t, executor.UpdateFeature(feature1, object1))
	assert.Equal(t, superObjectID, object1.ObjectID)

	feature2 := &web_feature.Feature{Name: "new-feature"}
	object2 := &example.Object{ID: backendStartID}
	assert.Equal(t, caskin.ErrCantChangeObjectType, executor.UpdateFeature(feature2, object2))
	object2.ID = featureStartID - 1
	assert.Equal(t, caskin.ErrCantChangeObjectType, executor.UpdateFeature(feature2, object2))
	object2.ID = featureStartID + 1
	assert.NoError(t, executor.UpdateFeature(feature2, object2))

	list1, err := executor.GetFeature()
	assert.NoError(t, err)
	assert.Len(t, list1, 5)
	assert.Equal(t, "new-feature", list1[2].ObjectCustomizedData.GetName())
}

func TestExecutorFeature_Get(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	_, err = executor.GetFeature()
	assert.Equal(t, caskin.ErrProviderGet, err)
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	_, err = executor.GetFeature()
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, err)

	provider.User = stage.SuperadminUser
	_, err = executor.GetFeature()
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, err)

	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.AdminUser
	list1, err := executor.GetFeature()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	provider.User = stage.SuperadminUser
	list2, err := executor.GetFeature()
	assert.NoError(t, err)
	assert.Len(t, list2, 5)
	assert.Equal(t, "feature-root", list2[0].ObjectCustomizedData.GetName())
	assert.Equal(t, superObjectID+1, list2[1].Object.GetParentID())
}
