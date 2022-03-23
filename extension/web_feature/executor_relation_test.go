package web_feature_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"github.com/stretchr/testify/assert"
)

func TestExecutorRelation_GetFeatureRelation(t *testing.T) {
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
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
	assert.Len(t, relation2, 5)
	assert.Len(t, relation2[5], 0)
	assert.Len(t, relation2[featureStartID], 3)
	assert.Len(t, relation2[featureStartID+1], 3)
	assert.Len(t, relation2[featureStartID+2], 3)
	assert.Len(t, relation2[featureStartID+3], 3)
}

func TestExecutorRelation_GetFeatureRelationByFeature(t *testing.T) {
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
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
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
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

	relation1 := web_feature.Relation{uint64(1), superObjectID + 3, frontendStartID, backendStartID, backendStartID + 2}
	assert.NoError(t, executor.ModifyFeatureRelationPerFeature(feature[1].Object, relation1))
	list1, err := executor.GetFeatureRelationByFeature(feature[1].Object)
	assert.NoError(t, err)
	assert.Len(t, list1, 4)
	assert.Equal(t, list1[0], superObjectID+3)
	assert.Equal(t, list1[1], frontendStartID)
	assert.Equal(t, list1[2], backendStartID)
	assert.Equal(t, list1[3], backendStartID+2)
}
