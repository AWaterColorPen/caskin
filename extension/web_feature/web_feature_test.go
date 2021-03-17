package web_feature_test

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"testing"

	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

const (
	frontendStartID = 8
	backendStartID = 12
	featureStartID = 19
)

func TestWebFeature(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	_, err = newWebFeature(stage)
	assert.NoError(t, err)

	object1 := web_feature.GetFeatureRootObject()
	assert.NotNil(t, object1)
	assert.Equal(t, uint64(frontendStartID - 3), object1.GetID())
	feature, err := caskin.Object2CustomizedData(object1, web_feature.FeatureFactory)
	assert.NoError(t, err)
	assert.Equal(t, web_feature.DefaultFeatureRootName, feature.(*web_feature.Feature).Name)
	assert.Equal(t, web_feature.DefaultFeatureRootDescription, feature.(*web_feature.Feature).Description)
	assert.Equal(t, web_feature.DefaultFeatureRootGroup, feature.(*web_feature.Feature).Group)

	object2 := web_feature.GetFrontendRootObject()
	assert.NotNil(t, object2)
	assert.Equal(t, uint64(frontendStartID - 2), object2.GetID())
	frontend, err := caskin.Object2CustomizedData(object2, web_feature.FrontendFactory)
	assert.NoError(t, err)
	assert.Equal(t, web_feature.DefaultFrontendRootKey, frontend.(*web_feature.Frontend).Key)
	assert.Equal(t, web_feature.DefaultFrontendRootType, frontend.(*web_feature.Frontend).Type)
	assert.Equal(t, web_feature.DefaultFrontendRootDescription, frontend.(*web_feature.Frontend).Description)
	assert.Equal(t, web_feature.DefaultFrontendRootGroup, frontend.(*web_feature.Frontend).Group)

	object3 := web_feature.GetBackendRootObject()
	assert.NotNil(t, object3)
	assert.Equal(t, uint64(frontendStartID - 1), object3.GetID())
	backend, err := caskin.Object2CustomizedData(object3, web_feature.BackendFactory)
	assert.NoError(t, err)
	assert.Equal(t, web_feature.DefaultBackendRootPath, backend.(*web_feature.Backend).Path)
	assert.Equal(t, web_feature.DefaultBackendRootMethod, backend.(*web_feature.Backend).Method)
	assert.Equal(t, web_feature.DefaultBackendRootDescription, backend.(*web_feature.Backend).Description)
	assert.Equal(t, web_feature.DefaultBackendRootGroup, backend.(*web_feature.Backend).Group)
}

func newWebFeature(stage *example.Stage) (*web_feature.WebFeature, error) {
	w, err := web_feature.New(stage.Caskin, nil)
	if err != nil {
		return nil, err
	}
	if err := web_feature.ManualCreateRootObject(stage.Options.MetaDB,
		stage.Options.EntryFactory.NewObject,
		stage.Options.GetSuperadminDomain()); err != nil {
		return nil, err
	}

	provider := caskin.NewCachedProvider(stage.SuperadminUser, stage.Options.GetSuperadminDomain())
	executor := w.GetExecutor(provider)
	frontend := []*web_feature.Frontend{
		{Key: "backend", Type: web_feature.FrontendTypeMenu},
		{Key: "frontend", Type: web_feature.FrontendTypeMenu},
		{Key: "feature", Type: web_feature.FrontendTypeMenu},
		{Key: "feature-sync", Type: web_feature.FrontendTypeSubFunction},
	}
	for _, v := range frontend {
		if err := executor.CreateFrontend(v); err != nil {
			return nil, err
		}
	}
	backend := []*web_feature.Backend{
		{Path: "api/backend", Method: "GET"},
		{Path: "api/backend", Method: "POST"},
		{Path: "api/frontend", Method: "GET"},
		{Path: "api/frontend", Method: "POST"},
		{Path: "api/feature", Method: "GET"},
		{Path: "api/feature", Method: "POST"},
		{Path: "api/sync", Method: "GET"},
	}
	for _, v := range backend {
		if err := executor.CreateBackend(v); err != nil {
			return nil, err
		}
	}
	feature := []*web_feature.Feature{
		{Name: "backend"},
		{Name: "frontend"},
		{Name: "feature"},
		{Name: "feature-sync"},
	}
	for _, v := range feature {
		if err := executor.CreateFeature(v); err != nil {
			return nil, err
		}
	}

	return w, nil
}
