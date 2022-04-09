package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestServer_DomainCreate(t *testing.T) {
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

func TestServer_DomainUpdate(t *testing.T) {
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

func TestServer_DomainRecover(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	domain1 := &example.Domain{
		Name: stage.Domain.Name,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, service.DomainRecover(domain1))
	assert.NoError(t, service.DomainDelete(stage.Domain))

	domain2 := &example.Domain{
		Name: stage.Domain.Name,
	}
	assert.NoError(t, service.DomainRecover(domain2))

	domain3 := &example.Domain{ID: 5}
	assert.Error(t, service.DomainRecover(domain3))
}

func TestServer_DomainDelete(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	assert.NoError(t, service.DomainDelete(stage.Domain))

	roles1, err := service.RoleGet(stage.MemberUser, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles1, 2)
	pairs1, err := service.UserRoleGet(stage.MemberUser, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, pairs1, 0)

	assert.NoError(t, service.DomainRecover(stage.Domain))
	roles2, err := service.RoleGet(stage.MemberUser, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles2, 2)
	pairs2, err := service.UserRoleGet(stage.MemberUser, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, pairs2, 0)
}

func TestServer_DomainReset(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	domain := &example.Domain{Name: "domain_02"}
	assert.NoError(t, service.DomainCreate(domain))

	// domain is no reset
	roles1, err := service.RoleGet(stage.SuperadminUser, domain)
	assert.NoError(t, err)
	assert.Len(t, roles1, 0)
	objects1, err := service.ObjectGet(stage.SuperadminUser, domain, caskin.Read)
	assert.NoError(t, err)
	assert.Len(t, objects1, 0)

	// reset domain with role and object
	assert.NoError(t, service.DomainReset(domain))
	roles2, err := service.RoleGet(stage.SuperadminUser, domain)
	assert.NoError(t, err)
	assert.Len(t, roles2, 2)
	objects2, err := service.ObjectGet(stage.SuperadminUser, domain, caskin.Read)
	assert.NoError(t, err)
	assert.Len(t, objects2, 3)

	// delete a object before initialize
	assert.NoError(t, service.ObjectDelete(stage.SuperadminUser, domain, &example.Object{ID: 3}))
	// initialize with new domain creator
	assert.NoError(t, service.DomainReset(stage.Domain))
	roles3, err := service.RoleGet(stage.SuperadminUser, domain)
	assert.NoError(t, err)
	assert.Len(t, roles3, 2)
	objects3, err := service.ObjectGet(stage.SuperadminUser, domain, caskin.Read)
	assert.NoError(t, err)
	assert.Len(t, objects3, 4)
	assert.Equal(t, "test", objects3[2].GetObjectType())
}
