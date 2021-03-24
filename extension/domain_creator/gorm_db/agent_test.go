package gorm_db_test

import (
	"testing"

	"github.com/awatercolorpen/caskin/extension/domain_creator/gorm_db"
	"github.com/stretchr/testify/assert"
)

func TestAgent_Create(t *testing.T) {
	factory, err := newFactory(t.TempDir())
	assert.NoError(t, err)
	agent := factory.GetAgent()
	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&gorm_db.DomainCreatorObject{Name: "object-1"}))
	assert.Equal(t, gorm_db.ErrRelativeIDOutOfIndex, agent.Create(&gorm_db.DomainCreatorObject{Name: "object-1", RelativeObjectID: 10}))
	assert.Equal(t, gorm_db.ErrNotSupport, agent.Create(&gorm_db.DomainCreatorObject{Name: "object-1", RelativeObjectID: 1, RelativeParentID: 1}))
	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&gorm_db.DomainCreatorObject{Name: "object-1", RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Create(&gorm_db.DomainCreatorObject{Name: "object-1", AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Create(&gorm_db.DomainCreatorObject{Name: "object-2", RelativeObjectID: 1}))
	assert.Error(t, agent.Create(&gorm_db.DomainCreatorObject{Name: "object-2", RelativeObjectID: 1}))

	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&gorm_db.DomainCreatorRole{Name: "role-1"}))
	assert.Equal(t, gorm_db.ErrRelativeIDOutOfIndex, agent.Create(&gorm_db.DomainCreatorRole{Name: "role-1", RelativeObjectID: 10}))
	assert.Equal(t, gorm_db.ErrNotSupport, agent.Create(&gorm_db.DomainCreatorRole{Name: "role-1", RelativeObjectID: 1, RelativeParentID: 1}))
	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&gorm_db.DomainCreatorRole{Name: "role-1", RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Create(&gorm_db.DomainCreatorRole{Name: "role-1", AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Create(&gorm_db.DomainCreatorRole{Name: "role-2", RelativeObjectID: 1}))
	assert.Error(t, agent.Create(&gorm_db.DomainCreatorRole{Name: "role-2", RelativeObjectID: 1}))

	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&gorm_db.DomainCreatorPolicy{RelativeRoleID: 1, RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&gorm_db.DomainCreatorPolicy{RelativeRoleID: 1, RelativeObjectID: 1, AbsoluteRoleID: 1}))
	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Create(&gorm_db.DomainCreatorPolicy{RelativeRoleID: 0, RelativeObjectID: 1}))
	assert.Equal(t, gorm_db.ErrRelativeIDOutOfIndex, agent.Create(&gorm_db.DomainCreatorPolicy{RelativeRoleID: 10, RelativeObjectID: 1}))
	assert.NoError(t, agent.Create(&gorm_db.DomainCreatorPolicy{AbsoluteRoleID: 10, AbsoluteObjectID: 10}))
	assert.NoError(t, agent.Create(&gorm_db.DomainCreatorPolicy{RelativeRoleID: 1, RelativeObjectID: 1}))
}

func TestAgent_Update(t *testing.T) {
	factory, err := newFactory(t.TempDir())
	assert.NoError(t, err)
	agent := factory.GetAgent()
	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&gorm_db.DomainCreatorObject{ID: 1, Name: "object-1"}))
	assert.Equal(t, gorm_db.ErrRelativeIDOutOfIndex, agent.Update(&gorm_db.DomainCreatorObject{ID: 1, Name: "object-1", RelativeObjectID: 10}))
	assert.Equal(t, gorm_db.ErrNotSupport, agent.Update(&gorm_db.DomainCreatorObject{ID: 1, Name: "object-1", RelativeObjectID: 1, RelativeParentID: 1}))
	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&gorm_db.DomainCreatorObject{ID: 1, Name: "object-1", RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Update(&gorm_db.DomainCreatorObject{ID: 1, Name: "object-1", AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Update(&gorm_db.DomainCreatorObject{ID: 1, Name: "object-2", RelativeObjectID: 1}))

	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&gorm_db.DomainCreatorRole{ID: 1, Name: "role-1"}))
	assert.Equal(t, gorm_db.ErrRelativeIDOutOfIndex, agent.Update(&gorm_db.DomainCreatorRole{ID: 1, Name: "role-1", RelativeObjectID: 10}))
	assert.Equal(t, gorm_db.ErrNotSupport, agent.Update(&gorm_db.DomainCreatorRole{ID: 1, Name: "role-1", RelativeObjectID: 1, RelativeParentID: 1}))
	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&gorm_db.DomainCreatorRole{ID: 1, Name: "role-1", RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Update(&gorm_db.DomainCreatorRole{ID: 1, Name: "role-1", AbsoluteObjectID: 1}))
	assert.NoError(t, agent.Update(&gorm_db.DomainCreatorRole{ID: 1, Name: "role-2", RelativeObjectID: 1}))

	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&gorm_db.DomainCreatorPolicy{ID: 1, RelativeRoleID: 1, RelativeObjectID: 1, AbsoluteObjectID: 1}))
	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&gorm_db.DomainCreatorPolicy{ID: 1, RelativeRoleID: 1, RelativeObjectID: 1, AbsoluteRoleID: 1}))
	assert.Equal(t, gorm_db.ErrRelativeIDAndAbsoluteRoleIDInCompatible, agent.Update(&gorm_db.DomainCreatorPolicy{ID: 1, RelativeRoleID: 0, RelativeObjectID: 1}))
	assert.Equal(t, gorm_db.ErrRelativeIDOutOfIndex, agent.Update(&gorm_db.DomainCreatorPolicy{ID: 1, RelativeRoleID: 10, RelativeObjectID: 1}))
	assert.NoError(t, agent.Update(&gorm_db.DomainCreatorPolicy{ID: 1, AbsoluteRoleID: 10, AbsoluteObjectID: 10}))
	assert.NoError(t, agent.Update(&gorm_db.DomainCreatorPolicy{ID: 1, RelativeRoleID: 1, RelativeObjectID: 1}))
}

func TestAgent_Delete(t *testing.T) {
	factory, err := newFactory(t.TempDir())
	assert.NoError(t, err)

	agent := factory.GetAgent()
	assert.Equal(t, gorm_db.ErrRelativeIDOutOfIndex, agent.Delete(&gorm_db.DomainCreatorObject{}, 1))
	assert.Equal(t, gorm_db.ErrRelativeIDOutOfIndex, agent.Delete(&gorm_db.DomainCreatorRole{}, 2))
	assert.Equal(t, gorm_db.ErrRelativeIDOutOfIndex, agent.Delete(&gorm_db.DomainCreatorObject{}, 3))

	assert.NoError(t, agent.Delete(&gorm_db.DomainCreatorPolicy{}, 8))
	assert.NoError(t, agent.Delete(&gorm_db.DomainCreatorPolicy{}, 7))
	assert.NoError(t, agent.Delete(&gorm_db.DomainCreatorRole{}, 2))

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
