package playground_test

import (
	"github.com/awatercolorpen/caskin"
	"testing"

	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestNewPlayground(t *testing.T) {
	playground.DictionaryDsn = "../configs/caskin.toml"
	_, err := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
}

func TestNewPlaygroundWithWatcher(t *testing.T) {
	playground.DictionaryDsn = "../configs/caskin.toml"
	_, err := playground.NewPlaygroundWithSqlitePathAndWatcher(t.TempDir(), nil)
	assert.NoError(t, err)

	w1 := &caskin.WatcherOption{}
	_, err = playground.NewPlaygroundWithSqlitePathAndWatcher(t.TempDir(), w1)
	assert.NoError(t, err)

	w2 := &caskin.WatcherOption{
		AutoLoad: 1000,
	}
	_, err = playground.NewPlaygroundWithSqlitePathAndWatcher(t.TempDir(), w2)
	assert.NoError(t, err)
}
