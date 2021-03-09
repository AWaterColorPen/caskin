package caskin_test

import (
	"encoding/json"
	"fmt"
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

func TestExecutor_GetObjects(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{
		User:   stage.SuperadminUser,
		Domain: stage.Domain,
	}
	executor := stage.Caskin.GetExecutor(provider)

	objects, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects, 3)

	objects, err = executor.GetObjects(caskin.ObjectTypeObject)
	assert.NoError(t, err)
	assert.Len(t, objects, 1)

	objects, err = executor.GetObjects(caskin.ObjectTypeObject, caskin.ObjectTypeRole)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
}

func TestExecutor_CreateObject(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{
		User:   stage.SuperadminUser,
		Domain: stage.Domain,
	}
	executor := stage.Caskin.GetExecutor(provider)

	objects, _ := executor.GetObjects(caskin.ObjectTypeObject)
	assert.Len(t, objects, 1)

	objectType := caskin.ObjectType("test_data")
	object := &example.Object{
		Name:     "object_01",
		Type:     objectType,
		DomainID: 1,
		ObjectID: objects[0].GetID(),
	}
	assert.NoError(t, executor.CreateObject(object))

	object2 := &example.Object{
		Name:     "object_01",
		Type:     objectType,
		DomainID: 1,
		ObjectID: objects[0].GetID(),
	}
	assert.Error(t, executor.CreateObject(object2))
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
		ID:        4,
	}
	assert.Error(t, executor.DeleteObject(object))
}

func TestExecutor_RecoverObject(t *testing.T) {
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
		ID:        3,
	}
	assert.NoError(t, executor.RecoverObject(object))
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverObject(object))
}

func TestExecutor_UpdateObject(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{
		User:   stage.SuperadminUser,
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

	policiesForRole, _ := executor.GetAllPoliciesForRole()
	bytes, _ := json.Marshal(policiesForRole)
	fmt.Println(string(bytes))

	provider.User = stage.AdminUser
	subObject.Name = "object_01_sub_new_name"
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateObject(subObject))
}
