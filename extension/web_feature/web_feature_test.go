package web_feature_test

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"sync"
	"testing"

	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestWebFeature(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	_, err = newWebFeature(stage)
	assert.NoError(t, err)

	object := web_feature.GetFeatureRootObject()
	assert.NotNil(t, web_feature.GetFeatureRootObject())
	assert.Equal(t, uint64(4), object.GetID())
	customized, err := caskin.Object2CustomizedData(object, web_feature.FeatureFactory)
	assert.NoError(t, err)
	assert.Equal(t, web_feature.DefaultFeatureRootName, customized.(*web_feature.Feature).Name)
	assert.Equal(t, web_feature.DefaultFeatureRootDescriptionName, customized.(*web_feature.Feature).Description)
	assert.Equal(t, web_feature.DefaultFeatureRootGroupName, customized.(*web_feature.Feature).Group)
}

func newWebFeature(stage *example.Stage) (*web_feature.WebFeature, error) {
	w, err := web_feature.New(stage.Caskin, nil)
	if err != nil {
		return nil, err
	}

	provider := caskin.NewCachedProvider(stage.SuperadminUser, stage.Options.GetSuperadminDomain())
	executor := w.GetExecutor(provider)
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

func parallel(input interface{}, function func(interface{}) interface{}, parallelNumber int) chan interface{} {
	wg := sync.WaitGroup{}

	in := make(chan interface{})
	out := make(chan interface{})
	for i := 0; i < parallelNumber; i++ {
		wg.Add(1)
		go func() {
			linq.FromChannel(in).ForEach(func(v interface{}) {
				out <- function(v)
			})
			wg.Done()
		}()
	}

	go linq.From(input).ToChannel(in)
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
