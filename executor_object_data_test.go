package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorObjectData_CreateCheck(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	executor := stage.Caskin.GetExecutor(provider)

	data1 := &example.OneObjectData{
		Name: "object_data_1",
		ObjectID: 3,
	}
	assert.Equal(t, caskin.ErrProviderGet, executor.CreateObjectDataCheckPermission(data1, ObjectTypeTest))

	provider.Domain = stage.Domain
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrInValidObjectType, executor.CreateObjectDataCheckPermission(data1, ObjectTypeTest))
	assert.NoError(t, executor.CreateObjectDataCheckPermission(data1, caskin.ObjectTypeDefault))

	data2 := &example.OneObjectData{
		Name: "object_data_2",
		ObjectID: 2,
	}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.CreateObjectDataCheckPermission(data2, caskin.ObjectTypeRole))
	provider.User = stage.AdminUser
	assert.NoError(t, executor.CreateObjectDataCheckPermission(data2, caskin.ObjectTypeRole))
}

func TestExecutorObjectData_RecoverCheck(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	provider := caskin.NewCachedProvider(nil, nil)
	provider.User = stage.AdminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	data1 := &example.OneObjectData{
		Name: "object_data_1",
		ObjectID: 3,
	}
	assert.Equal(t, caskin.ErrNotExists, executor.RecoverObjectDataCheckPermission(data1))

	data1.ObjectID = 4
	assert.NoError(t, executor.DB.Create(data1))
	assert.Equal(t, caskin.ErrAlreadyExists, executor.RecoverObjectDataCheckPermission(data1))
	assert.NoError(t, executor.DB.DeleteByID(data1, data1.GetID()))

	assert.NoError(t, executor.RecoverObjectDataCheckPermission(data1))

	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.RecoverObjectDataCheckPermission(data1))

	provider.User = stage.SubAdminUser
	assert.NoError(t, executor.RecoverObjectDataCheckPermission(data1))
}

func TestExecutorObjectData_DeleteCheck(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	assert.NoError(t, stage.AddSubAdmin())

	provider.User = stage.AdminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	data1 := &example.OneObjectData{
		ObjectID: 4,
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.DeleteObjectDataCheckPermission(data1))
	data1.ID = 1
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteObjectDataCheckPermission(data1))
	assert.NoError(t, executor.DB.Create(data1))

	data2 := &example.OneObjectData{ID: 1, ObjectID: 1}
	assert.Equal(t, caskin.ErrNotExists, executor.DeleteObjectDataCheckPermission(data2))
	data2.ObjectID = 4
	assert.NoError(t, executor.DeleteObjectDataCheckPermission(data2))
	provider.User = stage.SubAdminUser
	assert.NoError(t, executor.DeleteObjectDataCheckPermission(data2))
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.DeleteObjectDataCheckPermission(data2))
}

func TestExecutorObjectData_UpdateCheck(t *testing.T) {
	stage, _ := example.NewStageWithSqlitePath(t.TempDir())
	provider := caskin.NewCachedProvider(nil, nil)
	assert.NoError(t, stage.AddSubAdmin())

	provider.User = stage.AdminUser
	provider.Domain = stage.Domain
	executor := stage.Caskin.GetExecutor(provider)

	data1 := &example.OneObjectData{
		Name:     "object_data_1",
		ObjectID: 5,
	}
	assert.Equal(t, caskin.ErrEmptyID, executor.UpdateObjectDataCheckPermission(data1, &example.OneObjectData{}, caskin.ObjectTypeRole))
	data1.ID = 1
	assert.Equal(t, caskin.ErrNotExists, executor.UpdateObjectDataCheckPermission(data1, &example.OneObjectData{}, caskin.ObjectTypeRole))
	assert.NoError(t, executor.DB.Create(data1))

	data2 := &example.OneObjectData{
		ID: 1,
		Name: "object_data_3",
	}
	provider.User = stage.SubAdminUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateObjectDataCheckPermission(data2, &example.OneObjectData{}, caskin.ObjectTypeRole))
	data2.ObjectID = 5
	assert.NoError(t, executor.UpdateObjectDataCheckPermission(data2, &example.OneObjectData{}, caskin.ObjectTypeRole))
	provider.User = stage.MemberUser
	assert.Equal(t, caskin.ErrNoWritePermission, executor.UpdateObjectDataCheckPermission(data2, &example.OneObjectData{}, caskin.ObjectTypeRole))
	provider.User = stage.AdminUser
	data2.ObjectID = 4
	assert.Equal(t, caskin.ErrInValidObjectType, executor.UpdateObjectDataCheckPermission(data2, &example.OneObjectData{}, caskin.ObjectTypeRole))
}
