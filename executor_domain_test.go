package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestServer_DomainCreate_General(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	domain1 := &example.Domain{Name: "domain_02"}
	assert.NoError(t, service.DomainCreate(domain1))

	domain2 := &example.Domain{Name: "domain_02"}
	assert.Equal(t, caskin.ErrAlreadyExists, service.DomainCreate(domain2))

	domains1, err := service.DomainGet()
	assert.NoError(t, err)
	assert.Len(t, domains1, 2)

	domain3 := &example.Domain{Name: "domain_02"}
	assert.Equal(t, caskin.ErrEmptyID, service.DomainDelete(domain3))
	domain3.ID = domain2.ID
	assert.NoError(t, service.DomainDelete(domain3))

	domain4 := &example.Domain{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, service.DomainDelete(domain4))
	assert.NoError(t, service.DomainCreate(domain4))
}

func TestExecutorDomain_GeneralUpdate(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	domain1 := &example.Domain{
		ID:   stage.Domain.ID,
		Name: "domain_01_new_name",
	}
	assert.NoError(t, service.DomainUpdate(domain1))
	domain2 := &example.Domain{
		Name: "domain_01_new_name",
	}
	assert.Equal(t, caskin.ErrEmptyID, service.DomainUpdate(domain2))

	domain3 := &example.Domain{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, service.DomainUpdate(domain3))
}

func TestExecutorDomain_GeneralRecover(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	domain1 := &example.Domain{
		Name: stage.Domain.Name,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.DomainRecover(domain1))
	assert.NoError(t, executor.DomainDelete(stage.Domain))

	domain2 := &example.Domain{
		Name: stage.Domain.Name,
	}
	assert.NoError(t, executor.DomainRecover(domain2))

	domain3 := &example.Domain{ID: 5}
	assert.Error(t, executor.DomainRecover(domain3))
}

func TestExecutorDomain_GeneralDelete(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())

	service := stage.Service

	assert.NoError(t, executor.DomainDelete(stage.Domain))

	provider.Domain = stage.Domain
	provider.User = stage.SuperadminUser

	roles1, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles1, 2)
	pairs1, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, pairs1, 0)

	assert.NoError(t, executor.DomainRecover(stage.Domain))
	roles2, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles2, 2)
	pairs2, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, pairs2, 0)
}

func TestExecutorDomain_Reset(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())

	service := stage.Service

	domain := &example.Domain{Name: "domain_02"}
	assert.NoError(t, executor.DomainCreate(domain))

	// domain is no initialization
	provider.Domain = domain
	provider.User = stage.SuperadminUser
	roles1, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles1, 0)
	objects1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects1, 0)

	// initialized domain with role and object
	assert.NoError(t, executor.ReInitializeDomain(domain))
	roles2, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles2, 2)
	objects2, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects2, 3)

	// delete a object before initialize
	provider.Domain = stage.Domain
	assert.NoError(t, executor.DeleteObject(&example.Object{ID: 3}))
	// initialize with new domain creator
	stage.Options.DomainCreator = NewTestCreator
	assert.NoError(t, executor.ReInitializeDomain(stage.Domain))
	provider.User = stage.AdminUser
	roles3, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles3, 2)
	objects3, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects3, 4)
	assert.Equal(t, ObjectTypeTest, objects3[2].GetObjectType())
}
