package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutorObject(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{
		User:   stage.SuperadminUser,
		Domain: stage.Domain,
	}
	executor := stage.Caskin.GetExecutor(provider)

	objects, err := executor.GetObjects(caskin.ObjectTypeObject)
	assert.NoError(t, err)
	assert.Len(t, objects, 1)

	domains, err := executor.GetAllDomain()
	assert.NoError(t, err)
	assert.Len(t, domains, 1)

	objectType := caskin.ObjectType("test_data")
	object := &example.Object{
		Name:     "object_01",
		Type:     objectType,
		DomainID: 1,
		ObjectID: objects[0].GetID(),
	}
	assert.NoError(t, executor.CreateObject(object))

	subObject := &example.Object{
		Name:     "object_01_sub",
		Type:     objectType,
		DomainID: 1,
		ObjectID: objects[0].GetID(),
		ParentID: object.ID,
	}
	assert.NoError(t, executor.CreateObject(subObject))

	assert.NoError(t, executor.DeleteObject(object))
	objects, err = executor.GetObjects(objectType)
	assert.NoError(t, err)
	assert.Len(t, objects, 1)

	assert.NoError(t, executor.RecoverObject(object))
	objects, err = executor.GetObjects(objectType)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object.Name = "object_01_new_name"
	assert.NoError(t, executor.UpdateObject(object))

	assert.NoError(t, executor.DeleteObject(subObject))
	objects, err = executor.GetObjects(objectType)
	assert.NoError(t, err)
	assert.Len(t, objects, 1)

	assert.NoError(t, executor.RecoverObject(subObject))
	objects, err = executor.GetObjects(objectType)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	subObject.Name = "object_01__sub_new_name"
	assert.NoError(t, executor.UpdateObject(subObject))
}

func TestExecutorObject_GetObjects(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	objects1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects1, 3)

	objects2, err := executor.GetObjects(caskin.ObjectTypeObject)
	assert.NoError(t, err)
	assert.Len(t, objects2, 1)

	objects3, err := executor.GetObjects(caskin.ObjectTypeObject, caskin.ObjectTypeRole)
	assert.NoError(t, err)
	assert.Len(t, objects3, 2)

	provider.Domain = stage.Domain
	provider.User = stage.MemberUser
	objects4, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects4, 0)
}

func TestExecutorObject_GeneralCreate(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	object1 := &example.Object{
		Name:     "object_01",
		Type:     ObjectTypeTest,
	}
	assert.Equal(t, caskin.ErrProviderGet, executor.CreateObject(object1))

	provider.Domain = stage.Domain
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObject(object1))
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObject(object1))

	object1.ObjectID = 1
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObject(object1))
	provider.User = stage.AdminUser
	assert.NoError(t, executor.CreateObject(object1))

	object2 := &example.Object{
		Name:     "object_01",
		Type:     ObjectTypeTest,
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateObject(object2))

	object3 := &example.Object{
		Name:     "object_01",
		Type:     ObjectTypeTest,
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteObject(object3))
	object3.ID = object1.ID
	assert.NoError(t, executor.DeleteObject(object3))

	object4 := &example.Object{ID: 10, ObjectID: 1}
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteObject(object4))
	assert.NoError(t, executor.CreateObject(object4))
}

func TestExecutorObject_CreateSubNode(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	role1 := &example.Role{
		Name:     "sub_member_1",
		ObjectID: 2,
		ParentID: 3,
	}
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateRole(role1))

	role1.ParentID = 2
	assert.NoError(t, executor.CreateRole(role1))

	role1.ObjectID = 2
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateRole(role1))
	provider.User = stage.AdminUser
	assert.NoError(t, executor.CreateRole(role1))

	role2 := &example.Role{
		Name: "sub_admin_1",
	}
	assert.Equal(t, caskin.ErrAlreadyExists, executor.CreateRole(role2))

	role3 := &example.Role{
		Name: "sub_admin_1",
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteRole(role3))
	role3.ID = role1.ID
	assert.NoError(t, executor.DeleteRole(role3))

	role4 := &example.Role{ID: 5}
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteRole(role4))
	assert.NoError(t, executor.CreateRole(role4))
}

func TestExecutorObject_GeneralUpdate(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{
		User:   stage.AdminUser,
		Domain: stage.Domain,
	}
	executor := stage.Caskin.GetExecutor(provider)

	objectType := caskin.ObjectType("test_data")
	object := &example.Object{
		Name:     "object_01",
		Type:     objectType,
		DomainID: 1,
		ObjectID: 1,
	}
	assert.NoError(t, executor.CreateObject(object))

	object.Name = "object_01_new_name"
	assert.NoError(t, executor.UpdateObject(object))

	subObject := &example.Object{
		Name:     "object_01_sub",
		Type:     objectType,
		DomainID: 1,
		ObjectID: 4,
		ParentID: object.ID,
	}
	assert.NoError(t, executor.CreateObject(subObject))

	// policiesForRole, _ := executor.GetAllPoliciesForRole()
	// bytes, _ := json.Marshal(policiesForRole)
	// fmt.Println(string(bytes))

	provider.User = stage.AdminUser
	subObject.Name = "object_01_sub_new_name"
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateObject(subObject))
	assert.Error(t, executor.CreateObject(object1))

	provider.User = stage.MemberUser
	object2 := &example.Object{
		Name:     "object_02",
		Type:     objectType,
		DomainID: 1,
		ObjectID: objects[0].GetID(),
	}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObject(object2))
}

func TestExecutorObject_GeneralRecover(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{
		User:   stage.SuperadminUser,
		Domain: stage.Domain,
	}
	executor := stage.Caskin.GetExecutor(provider)

	objects, _ := executor.GetObjects()
	assert.Len(t, objects, 3)

	assert.NoError(t, executor.DeleteObject(objects[2]))

	objects, _ = executor.GetObjects()
	assert.Len(t, objects, 2)

	object := &example.Object{
		ID: 3,
	}
	assert.Error(t, executor.DeleteObject(object))

	provider.User = stage.MemberUser
	object.ID = 2
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteObject(object))
	assert.NoError(t, executor.RecoverObject(object))
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverObject(object))
}

func TestExecutor_DeleteObject(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{
		User:   stage.SuperadminUser,
		Domain: stage.Domain,
	}
	executor := stage.Caskin.GetExecutor(provider)

	objects, _ := executor.GetObjects()
	assert.Len(t, objects, 3)

	assert.NoError(t, executor.DeleteObject(objects[2]))

	objects, _ = executor.GetObjects()
	assert.Len(t, objects, 2)

	object := &example.Object{
		ID: 4,
	}
	assert.Error(t, executor.DeleteObject(object))
	assert.NoError(t, executor.CreateObject(object))

	object.Name = "object_01_new_name"
	assert.NoError(t, executor.UpdateObject(object))

	subObject := &example.Object{
		Name:     "object_01_sub",
		Type:     objectType,
		DomainID: 1,
		ObjectID: 4,
		ParentID: object.ID,
	}
	assert.NoError(t, executor.CreateObject(subObject))

	provider.User = stage.AdminUser
	subObject.Name = "object_01_sub_new_name"
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateObject(subObject))
}