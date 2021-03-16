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

func TestGetFeatureRootObject(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	assert.Nil(t, web_feature.GetFeatureRootObject())

	var input []int
	linq.Range(0, 5).ToSlice(&input)
	out := parallel(input, func(i interface{}) interface{} {
		_, _ = newWebFeature(stage)
		return web_feature.GetFeatureRootObject()
	}, 3)

	for v := range out {
		assert.NotNil(t, v)
		assert.Equal(t, v, web_feature.GetFeatureRootObject())
	}
}

func TestGetFeatureRootObject_MultiServer(t *testing.T) {
	stage, err := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
	assert.Nil(t, web_feature.GetFeatureRootObject())

	var input []int
	linq.Range(0, 5).ToSlice(&input)
	out := parallel(input, func(i interface{}) interface{} {
		_, _ = newWebFeature(stage)
		return web_feature.GetFeatureRootObject()
	}, 3)

	for v := range out {
		assert.NotNil(t, v)
		assert.Equal(t, v, web_feature.GetFeatureRootObject())
	}
}

func newWebFeature(stage *example.Stage) (*web_feature.WebFeature, error) {
	w, err := web_feature.New(stage.Caskin, nil)
	if err != nil {
		return nil, err
	}

	// provider := caskin.NewCachedProvider(stage.SuperadminUser, stage.Options.GetSuperadminDomain())
	// executor := w.GetExecutor(provider)
	// executor.

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