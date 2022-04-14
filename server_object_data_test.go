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
	assert.Equal(t, caskin.ErrNoWritePermission, service.ObjectDataCreateCheck(stage.Admin, stage.Domain, data1, "role"))
	data1.ObjectID = 2
	assert.Equal(t, caskin.ErrInValidObjectType, service.ObjectDataCreateCheck(stage.Admin, stage.Domain, data1, "role"))
	assert.NoError(t, service.ObjectDataCreateCheck(stage.Admin, stage.Domain, data1, "test"))

	data2 := &example.OneObjectData{Name: "object_data_2", ObjectID: 1}
	assert.Equal(t, caskin.ErrNoWritePermission, service.ObjectDataCreateCheck(stage.Member, stage.Domain, data2, caskin.ObjectTypeRole))
	assert.NoError(t, service.ObjectDataCreateCheck(stage.Admin, stage.Domain, data2, caskin.ObjectTypeRole))
}

func TestServer_ObjectData_RecoverCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 2}
	assert.Equal(t, caskin.ErrNotExists, service.ObjectDataRecoverCheck(stage.Admin, stage.Domain, data1))

	assert.NoError(t, stage.DB.Create(data1).Error)
	assert.Equal(t, caskin.ErrAlreadyExists, service.ObjectDataRecoverCheck(stage.Admin, stage.Domain, data1))
	assert.NoError(t, stage.DB.Delete(data1, data1.GetID()).Error)
	assert.NoError(t, service.ObjectDataRecoverCheck(stage.Admin, stage.Domain, data1))
	assert.NoError(t, service.ObjectDataRecoverCheck(stage.Member, stage.Domain, data1))
}

func TestServer_ObjectData_DeleteCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 2}
	assert.Equal(t, caskin.ErrEmptyID, service.ObjectDataDeleteCheck(stage.Admin, stage.Domain, data1))
	data1.ID = 1
	assert.Equal(t, caskin.ErrNotExists, service.ObjectDataDeleteCheck(stage.Admin, stage.Domain, data1))
	assert.NoError(t, stage.DB.Create(data1).Error)

	data2 := &example.OneObjectData{ID: 1, ObjectID: 1}
	assert.Equal(t, caskin.ErrNotExists, service.ObjectDataDeleteCheck(stage.Admin, stage.Domain, data2))
	data2.ObjectID = 2
	assert.NoError(t, service.ObjectDataDeleteCheck(stage.Admin, stage.Domain, data2))
	assert.NoError(t, service.ObjectDataDeleteCheck(stage.Member, stage.Domain, data2))
}

func TestServer_ObjectData_UpdateCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 1}
	assert.Equal(t, caskin.ErrEmptyID, service.ObjectDataUpdateCheck(stage.Admin, stage.Domain, data1, "role"))
	data1.ID = 1
	assert.Equal(t, caskin.ErrNotExists, service.ObjectDataUpdateCheck(stage.Admin, stage.Domain, data1, "role"))
	assert.NoError(t, stage.DB.Create(data1).Error)

	data2 := &example.OneObjectData{ID: 1, ObjectID: 2}
	assert.Equal(t, caskin.ErrInValidObjectType, service.ObjectDataUpdateCheck(stage.Admin, stage.Domain, data2, "role"))
	data2.ObjectID = 1
	assert.Equal(t, caskin.ErrInValidObjectType, service.ObjectDataUpdateCheck(stage.Admin, stage.Domain, data2, "test"))
	assert.Equal(t, caskin.ErrNoWritePermission, service.ObjectDataUpdateCheck(stage.Member, stage.Domain, data2, "role"))
	assert.NoError(t, service.ObjectDataUpdateCheck(stage.Admin, stage.Domain, data2, "role"))
}

func TestServer_ObjectData_ModifyCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 1}
	assert.Equal(t, caskin.ErrEmptyID, service.ObjectDataModifyCheck(stage.Admin, stage.Domain, data1))
	data1.ID = 1
	assert.Equal(t, caskin.ErrNotExists, service.ObjectDataModifyCheck(stage.Admin, stage.Domain, data1))
	assert.NoError(t, stage.DB.Create(data1).Error)

	data2 := &example.OneObjectData{ID: 1, ObjectID: 2}
	assert.Equal(t, caskin.ErrNotExists, service.ObjectDataModifyCheck(stage.Admin, stage.Domain, data2))
	data2.ObjectID = 1
	assert.Equal(t, caskin.ErrNoWritePermission, service.ObjectDataModifyCheck(stage.Member, stage.Domain, data2))
	assert.NoError(t, service.ObjectDataModifyCheck(stage.Admin, stage.Domain, data2))
}

func TestServer_ObjectData_GetCheck(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: 1}
	assert.Equal(t, caskin.ErrEmptyID, service.ObjectDataGetCheck(stage.Admin, stage.Domain, data1))
	data1.ID = 1
	assert.Equal(t, caskin.ErrNotExists, service.ObjectDataGetCheck(stage.Admin, stage.Domain, data1))
	assert.NoError(t, stage.DB.Create(data1).Error)

	data2 := &example.OneObjectData{ID: 1, ObjectID: 2}
	assert.Equal(t, caskin.ErrNotExists, service.ObjectDataGetCheck(stage.Admin, stage.Domain, data2))
	data2.ObjectID = 1
	assert.NoError(t, service.ObjectDataGetCheck(stage.Admin, stage.Domain, data2))
	assert.NoError(t, service.ObjectDataGetCheck(stage.Member, stage.Domain, data2))
}
