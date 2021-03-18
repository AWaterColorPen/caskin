package web_feature_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"github.com/stretchr/testify/assert"
)

func TestExecutorRelation_GetFeatureRelation(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	_, err = executor.GetFeatureRelation()
	assert.Equal(t, caskin.ErrProviderGet, err)

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	_, err = executor.GetFeatureRelation()
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, err)
	provider.Domain = stage.Options.GetSuperadminDomain()
	relation1, err := executor.GetFeatureRelation()
	assert.NoError(t, err)
	assert.Len(t, relation1, 0)

	provider.User = stage.SuperadminUser
	relation2, err := executor.GetFeatureRelation()
	assert.NoError(t, err)
	assert.Len(t, relation2, 4)
}

func TestExecutorRelation_GetFeatureRelationByFeature(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	feature, err := executor.GetFeature()
	assert.NoError(t, err)
	assert.Len(t, feature, 5)

	list1, err := executor.GetFeatureRelationByFeature(feature[0].Object)
	assert.NoError(t, err)
	assert.Len(t, list1, 0)
	list2, err := executor.GetFeatureRelationByFeature(feature[1].Object)
	assert.NoError(t, err)
	assert.Len(t, list2, 3)
	list3, err := executor.GetFeatureRelationByFeature(feature[2].Object)
	assert.NoError(t, err)
	assert.Len(t, list3, 3)
	list4, err := executor.GetFeatureRelationByFeature(feature[3].Object)
	assert.NoError(t, err)
	assert.Len(t, list4, 3)
	list5, err := executor.GetFeatureRelationByFeature(feature[4].Object)
	assert.NoError(t, err)
	assert.Len(t, list5, 3)

}

func TestExecutorRelation_ModifyFeatureRelationPerFeature(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	feature, err := executor.GetFeature()
	assert.NoError(t, err)
	assert.Len(t, feature, 5)

	relation1 := web_feature.FeatureRelation{uint64(7), uint64(8), uint64(12), uint64(14)}
	assert.NoError(t, executor.ModifyFeatureRelationPerFeature(feature[1].Object, relation1))
	list1, err := executor.GetFeatureRelationByFeature(feature[1].Object)
	assert.NoError(t, err)
	assert.Len(t, list1, 3)
	assert.Equal(t, list1[0], uint64(8))
	assert.Equal(t, list1[1], uint64(12))
	assert.Equal(t, list1[2], uint64(14))
}
