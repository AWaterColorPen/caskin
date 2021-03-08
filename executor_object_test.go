package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutorObject(t *testing.T) {
	stage, _ := getStage(t)
	provider := example.Provider{
		User:   stage.SuperadminUser,
		Domain: stage.Domain,
	}
	executor := stage.Caskin.GetExecutor(provider)

	objects, err := executor.GetObject(caskin.ObjectTypeObject)
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
		Name:    "object_01_sub",
		Type:     objectType,
		DomainID: 1,
		ObjectID: objects[0].GetID(),
		ParentID: object.ID,
	}
	assert.NoError(t, executor.CreateObject(subObject))

	assert.NoError(t, executor.DeleteObject(object))
	objects, err = executor.GetObject(objectType)
	assert.NoError(t, err)
	assert.Len(t, objects, 1)

	assert.NoError(t, executor.RecoverObject(object))
	objects, err = executor.GetObject(objectType)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object.Name = "object_01_new_name"
	assert.Error(t, executor.UpdateObject(object))

	assert.NoError(t, executor.DeleteObject(subObject))
	objects, err = executor.GetObject(objectType)
	assert.NoError(t, err)
	assert.Len(t, objects, 1)

	assert.NoError(t, executor.RecoverObject(subObject))
	objects, err = executor.GetObject(objectType)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	subObject.Name = "object_01__sub_new_name"
	assert.Error(t, executor.UpdateObject(subObject))
}
