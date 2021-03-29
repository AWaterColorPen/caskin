package domain_creator_test

import (
	"github.com/awatercolorpen/caskin/extension/domain_creator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgent_Create(t *testing.T) {
	factory, err := newFactory(t.TempDir())
	assert.NoError(t, err)
	agent := factory.GetAgent()
	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&domain_creator.DomainCreatorObject{Name: "object-1"}))
	assert.Equal(t, domain_creator.ErrRelativeIDOutOfIndex, agent.Create(&domain_creator.DomainCreatorObject{Name: "object-1", RelativeObjectID: 10}))
	assert.Equal(t, domain_creator.ErrNotSupport, agent.Create(&domain_creator.DomainCreatorObject{Name: "object-1", RelativeObjectID: 1, RelativeParentID: 1}))
	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&domain_creator.DomainCreatorObject{Name: "object-1", RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Create(&domain_creator.DomainCreatorObject{Name: "object-1", AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Create(&domain_creator.DomainCreatorObject{Name: "object-2", RelativeObjectID: 1}))
	assert.Error(t, agent.Create(&domain_creator.DomainCreatorObject{Name: "object-2", RelativeObjectID: 1}))

	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&domain_creator.DomainCreatorRole{Name: "role-1"}))
	assert.Equal(t, domain_creator.ErrRelativeIDOutOfIndex, agent.Create(&domain_creator.DomainCreatorRole{Name: "role-1", RelativeObjectID: 10}))
	assert.Equal(t, domain_creator.ErrNotSupport, agent.Create(&domain_creator.DomainCreatorRole{Name: "role-1", RelativeObjectID: 1, RelativeParentID: 1}))
	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&domain_creator.DomainCreatorRole{Name: "role-1", RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Create(&domain_creator.DomainCreatorRole{Name: "role-1", AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Create(&domain_creator.DomainCreatorRole{Name: "role-2", RelativeObjectID: 1}))
	assert.Error(t, agent.Create(&domain_creator.DomainCreatorRole{Name: "role-2", RelativeObjectID: 1}))

	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&domain_creator.DomainCreatorPolicy{RelativeRoleID: 1, RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&domain_creator.DomainCreatorPolicy{RelativeRoleID: 1, RelativeObjectID: 1, AbsoluteRoleID: 1}))
	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&domain_creator.DomainCreatorPolicy{RelativeRoleID: 0, RelativeObjectID: 1}))
	assert.Equal(t, domain_creator.ErrRelativeIDOutOfIndex, agent.Create(&domain_creator.DomainCreatorPolicy{RelativeRoleID: 10, RelativeObjectID: 1}))
	assert.NoError(t, agent.Create(&domain_creator.DomainCreatorPolicy{AbsoluteRoleID: 10, AbsoluteObjectID: 10}))
	assert.NoError(t, agent.Create(&domain_creator.DomainCreatorPolicy{RelativeRoleID: 1, RelativeObjectID: 1}))
}

func TestAgent_Update(t *testing.T) {
	factory, err := newFactory(t.TempDir())
	assert.NoError(t, err)
	agent := factory.GetAgent()
	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&domain_creator.DomainCreatorObject{ID: 1, Name: "object-1"}))
	assert.Equal(t, domain_creator.ErrRelativeIDOutOfIndex, agent.Update(&domain_creator.DomainCreatorObject{ID: 1, Name: "object-1", RelativeObjectID: 10}))
	assert.Equal(t, domain_creator.ErrNotSupport, agent.Update(&domain_creator.DomainCreatorObject{ID: 1, Name: "object-1", RelativeObjectID: 1, RelativeParentID: 1}))
	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&domain_creator.DomainCreatorObject{ID: 1, Name: "object-1", RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Update(&domain_creator.DomainCreatorObject{ID: 1, Name: "object-1", AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Update(&domain_creator.DomainCreatorObject{ID: 1, Name: "object-2", RelativeObjectID: 1}))

	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&domain_creator.DomainCreatorRole{ID: 1, Name: "role-1"}))
	assert.Equal(t, domain_creator.ErrRelativeIDOutOfIndex, agent.Update(&domain_creator.DomainCreatorRole{ID: 1, Name: "role-1", RelativeObjectID: 10}))
	assert.Equal(t, domain_creator.ErrNotSupport, agent.Update(&domain_creator.DomainCreatorRole{ID: 1, Name: "role-1", RelativeObjectID: 1, RelativeParentID: 1}))
	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&domain_creator.DomainCreatorRole{ID: 1, Name: "role-1", RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Update(&domain_creator.DomainCreatorRole{ID: 1, Name: "role-1", AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Update(&domain_creator.DomainCreatorRole{ID: 1, Name: "role-2", RelativeObjectID: 1}))

	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&domain_creator.DomainCreatorPolicy{ID: 1, RelativeRoleID: 1, RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&domain_creator.DomainCreatorPolicy{ID: 1, RelativeRoleID: 1, RelativeObjectID: 1, AbsoluteRoleID: 1}))
	assert.Equal(t, domain_creator.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&domain_creator.DomainCreatorPolicy{ID: 1, RelativeRoleID: 0, RelativeObjectID: 1}))
	assert.Equal(t, domain_creator.ErrRelativeIDOutOfIndex, agent.Update(&domain_creator.DomainCreatorPolicy{ID: 1, RelativeRoleID: 10, RelativeObjectID: 1}))
	assert.NoError(t, agent.Update(&domain_creator.DomainCreatorPolicy{ID: 1, AbsoluteRoleID: 10, AbsoluteObjectID: 10}))
	assert.NoError(t, agent.Update(&domain_creator.DomainCreatorPolicy{ID: 1, RelativeRoleID: 1, RelativeObjectID: 1}))
}

func TestAgent_Delete(t *testing.T) {
	factory, err := newFactory(t.TempDir())
	assert.NoError(t, err)

	agent := factory.GetAgent()
	assert.Equal(t, domain_creator.ErrRelativeIDOutOfIndex, agent.Delete(&domain_creator.DomainCreatorObject{}, 1))
	assert.Equal(t, domain_creator.ErrRelativeIDOutOfIndex, agent.Delete(&domain_creator.DomainCreatorRole{}, 2))
	assert.Equal(t, domain_creator.ErrRelativeIDOutOfIndex, agent.Delete(&domain_creator.DomainCreatorObject{}, 3))

	assert.NoError(t, agent.Delete(&domain_creator.DomainCreatorPolicy{}, 8))
	assert.NoError(t, agent.Delete(&domain_creator.DomainCreatorPolicy{}, 7))
	assert.NoError(t, agent.Delete(&domain_creator.DomainCreatorRole{}, 2))

	list1, err := agent.GetDomainCreatorObject()
	assert.NoError(t, err)
	assert.Len(t, list1, 3)
	list2, err := agent.GetDomainCreatorRole()
	assert.NoError(t, err)
	assert.Len(t, list2, 1)
	list3, err := agent.GetDomainCreatorPolicy()
	assert.NoError(t, err)
	assert.Len(t, list3, 6)
}
