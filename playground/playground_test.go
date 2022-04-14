package playground_test

import (
	"testing"

	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestNewPlayground(t *testing.T) {
	playground.DictionaryDsn = "../configs/caskin.toml"
	_, err := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, err)
}
