package web_feature_test

import (
	"path/filepath"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"github.com/stretchr/testify/assert"
)

func TestDumpFileStruct(t *testing.T) {
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.BuildVersion())
	assert.NoError(t, executor.SyncLatestVersionToAllDomain())

	dump, err := executor.Dump()
	assert.NoError(t, err)

	tempFile := filepath.Join(t.TempDir(), "dump_file_struct")
	dfs := &web_feature.DumpFileStruct{}
	assert.NoError(t, dfs.ImportFromDump(dump))
	assert.NoError(t, dfs.IsValid())
	assert.NoError(t, dfs.ExportToFile(tempFile))

	dfs2 := &web_feature.DumpFileStruct{}
	assert.NoError(t, dfs2.ImportFromFile(tempFile))
	assert.NoError(t, dfs2.IsValid())
}
