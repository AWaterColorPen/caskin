package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestServer_ObjectData_CreateCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 3}
	assert.Equal(t, caskin.ErrNoWritePermission, service.CheckCreateObjectData(stage.Admin, stage.Domain, data1, "role"))
	data1.ObjectID = 2
	assert.Equal(t, caskin.ErrInValidObjectType, service.CheckCreateObjectData(stage.Admin, stage.Domain, data1, "role"))
	assert.NoError(t, service.CheckCreateObjectData(stage.Admin, stage.Domain, data1, "test"))

	data2 := &example.OneObjectData{Name: "object_data_2", ObjectID: 1}
	assert.Equal(t, caskin.ErrNoWritePermission, service.CheckCreateObjectData(stage.Member, stage.Domain, data2, caskin.ObjectTypeRole))
	assert.NoError(t, service.CheckCreateObjectData(stage.Admin, stage.Domain, data2, caskin.ObjectTypeRole))
}

func TestServer_ObjectData_RecoverCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 2}
	assert.Equal(t, caskin.ErrNotExists, service.CheckRecoverObjectData(stage.Admin, stage.Domain, data1))

	assert.NoError(t, stage.DB.Create(data1).Error)
	assert.Equal(t, caskin.ErrAlreadyExists, service.CheckRecoverObjectData(stage.Admin, stage.Domain, data1))
	assert.NoError(t, stage.DB.Delete(data1, data1.GetID()).Error)
	assert.NoError(t, service.CheckRecoverObjectData(stage.Admin, stage.Domain, data1))
	assert.NoError(t, service.CheckRecoverObjectData(stage.Member, stage.Domain, data1))
}

func TestServer_ObjectData_DeleteCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 2}
	assert.Equal(t, caskin.ErrEmptyID, service.CheckDeleteObjectData(stage.Admin, stage.Domain, data1))
	data1.ID = 1
	assert.Equal(t, caskin.ErrNotExists, service.CheckDeleteObjectData(stage.Admin, stage.Domain, data1))
	assert.NoError(t, stage.DB.Create(data1).Error)

	data2 := &example.OneObjectData{ID: 1, ObjectID: 1}
	assert.Equal(t, caskin.ErrNotExists, service.CheckDeleteObjectData(stage.Admin, stage.Domain, data2))
	data2.ObjectID = 2
	assert.NoError(t, service.CheckDeleteObjectData(stage.Admin, stage.Domain, data2))
	assert.NoError(t, service.CheckDeleteObjectData(stage.Member, stage.Domain, data2))
}

func TestServer_ObjectData_UpdateCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 1}
	assert.Equal(t, caskin.ErrEmptyID, service.CheckUpdateObjectData(stage.Admin, stage.Domain, data1, "role"))
	data1.ID = 1
	assert.Equal(t, caskin.ErrNotExists, service.CheckUpdateObjectData(stage.Admin, stage.Domain, data1, "role"))
	assert.NoError(t, stage.DB.Create(data1).Error)

	data2 := &example.OneObjectData{ID: 1, ObjectID: 2}
	assert.Equal(t, caskin.ErrInValidObjectType, service.CheckUpdateObjectData(stage.Admin, stage.Domain, data2, "role"))
	data2.ObjectID = 1
	assert.Equal(t, caskin.ErrInValidObjectType, service.CheckUpdateObjectData(stage.Admin, stage.Domain, data2, "test"))
	assert.Equal(t, caskin.ErrNoWritePermission, service.CheckUpdateObjectData(stage.Member, stage.Domain, data2, "role"))
	assert.NoError(t, service.CheckUpdateObjectData(stage.Admin, stage.Domain, data2, "role"))
}

func TestServer_ObjectData_ModifyCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 1}
	assert.Equal(t, caskin.ErrEmptyID, service.CheckModifyObjectData(stage.Admin, stage.Domain, data1))
	data1.ID = 1
	assert.Equal(t, caskin.ErrNotExists, service.CheckModifyObjectData(stage.Admin, stage.Domain, data1))
	assert.NoError(t, stage.DB.Create(data1).Error)

	data2 := &example.OneObjectData{ID: 1, ObjectID: 2}
	assert.Equal(t, caskin.ErrNotExists, service.CheckModifyObjectData(stage.Admin, stage.Domain, data2))
	data2.ObjectID = 1
	assert.Equal(t, caskin.ErrNoWritePermission, service.CheckModifyObjectData(stage.Member, stage.Domain, data2))
	assert.NoError(t, service.CheckModifyObjectData(stage.Admin, stage.Domain, data2))
}

func TestServer_ObjectData_GetCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 1}
	assert.Equal(t, caskin.ErrEmptyID, service.CheckGetObjectData(stage.Admin, stage.Domain, data1))
	data1.ID = 1
	assert.Equal(t, caskin.ErrNotExists, service.CheckGetObjectData(stage.Admin, stage.Domain, data1))
	assert.NoError(t, stage.DB.Create(data1).Error)

	data2 := &example.OneObjectData{ID: 1, ObjectID: 2}
	assert.Equal(t, caskin.ErrNotExists, service.CheckGetObjectData(stage.Admin, stage.Domain, data2))
	data2.ObjectID = 1
	assert.NoError(t, service.CheckGetObjectData(stage.Admin, stage.Domain, data2))
	assert.NoError(t, service.CheckGetObjectData(stage.Member, stage.Domain, data2))
}
