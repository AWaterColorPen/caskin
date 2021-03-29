package domain_creator_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestCreator(t *testing.T) {
	sqlitePath := t.TempDir()
	factory, err := newFactory(sqlitePath)
	assert.NoError(t, err)

	stage, _ := example.NewStageWithSqlitePath(sqlitePath)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)
	stage.Options.DomainCreator = factory.NewCreator

	domain := &example.Domain{Name: "domain_02"}
	assert.NoError(t, executor.CreateDomain(domain))
	provider.User = stage.SuperadminUser
	provider.Domain = domain
	assert.NoError(t, executor.ReInitializeDomain(domain))
	roles1, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles1, 2)
	objects1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects1, 3)
}
