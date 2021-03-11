package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorDomain_GeneralCreate(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	domain1 := &example.Domain{Name: "domain_02"}
	assert.NoError(t, executor.CreateDomain(domain1))

	domain2 := &example.Domain{Name: "domain_02"}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateDomain(domain2))

	domains1, err := executor.GetAllDomain()
	assert.NoError(t, err)
	assert.Len(t, domains1, 2)

	domain3 := &example.Domain{
		Name: "domain_02",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteDomain(domain3))
	domain3.ID = domain2.ID
	assert.NoError(t, executor.DeleteDomain(domain3))

	domain4 := &example.Domain{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteDomain(domain4))
	assert.NoError(t, executor.CreateDomain(domain4))
}

func TestExecutorDomain_GeneralUpdate(t *testing.T) {
	stage, _ := newStage(t)
	executor := stage.Caskin.GetExecutor(caskin.NewCachedProvider(nil, nil, nil))

	domain1 := &example.Domain{
		ID:   stage.Domain.ID,
		Name: "domain_01_new_name",
	}
	assert.NoError(t, executor.UpdateDomain(domain1))
	domain2 := &example.Domain{
		Name: "domain_01_new_name",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.UpdateDomain(domain2))

	domain3 := &example.Domain{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.UpdateDomain(domain3))
}

func TestExecutorDomain_GeneralRecover(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	domain1 := &example.Domain{
		Name: stage.Domain.Name,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverDomain(domain1))
	assert.NoError(t, executor.DeleteDomain(stage.Domain))

	domain2 := &example.Domain{
		Name: stage.Domain.Name,
	}
	assert.NoError(t, executor.RecoverDomain(domain2))

	domain3 := &example.Domain{ID: 5}
	assert.Error(t, executor.RecoverDomain(domain3))
}

func TestExecutorDomain_GeneralDelete(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	assert.NoError(t, executor.DeleteDomain(stage.Domain))

	provider.Domain = stage.Domain
	provider.User = stage.SuperadminUser

	roles1, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles1, 2)
	pairs1, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, pairs1, 0)

	assert.NoError(t, executor.RecoverDomain(stage.Domain))
	roles2, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles2, 2)
	pairs2, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, pairs2, 0)
}

func TestExecutorDomain_Initialize(t *testing.T) {
	stage, _ := newStage(t)
	provider := caskin.NewCachedProvider(nil, nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	domain := &example.Domain{Name: "domain_02"}
	assert.NoError(t, executor.CreateDomain(domain))

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

	// initialize with new domain creator
	stage.Options.DomainCreator = NewTestCreator
	assert.NoError(t, executor.ReInitializeDomain(stage.Domain))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	roles3, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles3, 2)
	objects3, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects3, 4)
}

type testCreator struct {
	domain  caskin.Domain
	objects caskin.Objects
	roles   caskin.Roles
}

func NewTestCreator(domain caskin.Domain) caskin.Creator {
	return &testCreator{domain: domain}
}

const (
	ObjectTypeTest caskin.ObjectType = "test"
)

func (t *testCreator) BuildCreator() (caskin.Roles, caskin.Objects) {
	role0 := &example.Role{Name: "admin", DomainID: t.domain.GetID()}
	role1 := &example.Role{Name: "member", DomainID: t.domain.GetID()}
	t.roles = []caskin.Role{role0, role1}

	object0 := &example.Object{Name: string(caskin.ObjectTypeObject), Type: caskin.ObjectTypeObject, DomainID: t.domain.GetID()}
	object1 := &example.Object{Name: string(caskin.ObjectTypeRole), Type: caskin.ObjectTypeRole, DomainID: t.domain.GetID()}
	object2 := &example.Object{Name: string(caskin.ObjectTypeDefault), Type: caskin.ObjectTypeDefault, DomainID: t.domain.GetID()}
	object3 := &example.Object{Name: string(ObjectTypeTest), Type: ObjectTypeTest, DomainID: t.domain.GetID()}
	t.objects = []caskin.Object{object0, object1, object2, object3}

	return t.roles, t.objects
}

func (t *testCreator) SetRelation() {
	ooId := t.objects[0].GetID()
	for _, object := range t.objects {
		object.SetObjectId(ooId)
	}

	roId := t.objects[1].GetID()
	for _, role := range t.roles {
		role.SetObjectId(roId)
	}
}

func (t *testCreator) GetRoles() caskin.Roles {
	return t.roles
}

func (t *testCreator) GetObjects() caskin.Objects {
	return t.objects
}

func (t *testCreator) GetPolicy() []*caskin.Policy {
	return []*caskin.Policy{
		{t.roles[0], t.objects[0], t.domain, caskin.Read},
		{t.roles[0], t.objects[0], t.domain, caskin.Write},
		{t.roles[0], t.objects[1], t.domain, caskin.Read},
		{t.roles[0], t.objects[1], t.domain, caskin.Write},
		{t.roles[0], t.objects[2], t.domain, caskin.Read},
		{t.roles[0], t.objects[2], t.domain, caskin.Write},
		{t.roles[0], t.objects[3], t.domain, caskin.Read},
		{t.roles[0], t.objects[3], t.domain, caskin.Write},
		{t.roles[1], t.objects[2], t.domain, caskin.Read},
		{t.roles[1], t.objects[2], t.domain, caskin.Write},
	}
}
