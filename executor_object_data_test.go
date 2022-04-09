package caskin_test

import (
	"testing"
)

func TestServer_ObjectData_CreateCheck(t *testing.T) {
	//	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	//	service := stage.Service
	//
	//	data1 := &example.OneObjectData{
	//		Name:     "object_data_1",
	//		ObjectID: 3,
	//	}
	//	assert.Equal(t, caskin.ErrProviderGet, service.ObjectDataCreateCheck(data1, ObjectTypeTest))
	//
	//	assert.Equal(t, caskin.ErrInValidObjectType, service.ObjectDataCreateCheck(data1, ObjectTypeTest))
	//	assert.NoError(t, service.ObjectDataCreateCheck(data1, caskin.ObjectTypeDefault))
	//
	//	data2 := &example.OneObjectData{
	//		Name:     "object_data_2",
	//		ObjectID: 2,
	//	}
	//	assert.Equal(t, caskin.ErrNoWritePermission, service.ObjectDataCreateCheck(data2, caskin.ObjectTypeRole))
	//	assert.NoError(t, service.ObjectDataCreateCheck(data2, caskin.ObjectTypeRole))
}

func TestServer_ObjectData_RecoverCheck(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	//service := stage.Service
	//
	//data1 := &example.OneObjectData{
	//	Name:     "object_data_1",
	//	ObjectID: 3,
	//}
	//assert.Equal(t, caskin.ErrNotExists, service.ObjectDataRecoverCheck(data1))
	//
	//data1.ObjectID = 4
	//assert.NoError(t, service.DB.Create(data1))
	//assert.Equal(t, caskin.ErrAlreadyExists, service.ObjectDataRecoverCheck(data1))
	//assert.NoError(t, service.DB.DeleteByID(data1, data1.GetID()))
	//
	//assert.NoError(t, service.ObjectDataRecoverCheck(data1))
	//
	//provider.User = stage.MemberUser
	//assert.Equal(t, caskin.ErrNoWritePermission, service.ObjectDataRecoverCheck(data1))
	//
	//provider.User = stage.SubAdminUser
	//assert.NoError(t, service.ObjectDataRecoverCheck(data1))
}

func TestServer_ObjectData_DeleteCheck(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//data1 := &example.OneObjectData{
	//	ObjectID: 4,
	//}
	//assert.Equal(t, caskin.ErrEmptyID, service.ObjectDataDeleteCheck(data1))
	//data1.ID = 1
	//assert.Equal(t, caskin.ErrNotExists, service.ObjectDataDeleteCheck(data1))
	//assert.NoError(t, service.DB.Create(data1))
	//
	//data2 := &example.OneObjectData{ID: 1, ObjectID: 1}
	//assert.Equal(t, caskin.ErrNotExists, service.ObjectDataDeleteCheck(data2))
	//data2.ObjectID = 4
	//assert.NoError(t, service.ObjectDataDeleteCheck(data2))
	//provider.User = stage.SubAdminUser
	//assert.NoError(t, service.ObjectDataDeleteCheck(data2))
	//provider.User = stage.MemberUser
	//assert.Equal(t, caskin.ErrNoWritePermission, service.ObjectDataDeleteCheck(data2))
}

func TestServer_ObjectData_UpdateCheck(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//
	//provider.User = stage.AdminUser
	//provider.Domain = stage.Domain
	//service := stage.Service
	//
	//data1 := &example.OneObjectData{
	//	Name:     "object_data_1",
	//	ObjectID: 5,
	//}
	//assert.Equal(t, caskin.ErrEmptyID, service.ObjectDataUpdateCheck(data1, &example.OneObjectData{}, caskin.ObjectTypeRole))
	//data1.ID = 1
	//assert.Equal(t, caskin.ErrNotExists, service.ObjectDataUpdateCheck(data1, &example.OneObjectData{}, caskin.ObjectTypeRole))
	//assert.NoError(t, service.DB.Create(data1))
	//
	//data2 := &example.OneObjectData{
	//	ID:   1,
	//	Name: "object_data_3",
	//}
	//provider.User = stage.SubAdminUser
	//assert.Equal(t, caskin.ErrNoWritePermission, service.ObjectDataUpdateCheck(data2, &example.OneObjectData{}, caskin.ObjectTypeRole))
	//data2.ObjectID = 5
	//assert.NoError(t, service.ObjectDataUpdateCheck(data2, &example.OneObjectData{}, caskin.ObjectTypeRole))
	//provider.User = stage.MemberUser
	//assert.Equal(t, caskin.ErrNoWritePermission, service.ObjectDataUpdateCheck(data2, &example.OneObjectData{}, caskin.ObjectTypeRole))
	//provider.User = stage.AdminUser
	//data2.ObjectID = 4
	//assert.Equal(t, caskin.ErrInValidObjectType, service.ObjectDataUpdateCheck(data2, &example.OneObjectData{}, caskin.ObjectTypeRole))
	//
	//// it should not avoid object_data to change type
	//// it maybe want to change type from default to special or from special to default
	//data2.ObjectID = 3
	//assert.Equal(t, caskin.ErrInValidObjectType, service.ObjectDataUpdateCheck(data2, &example.OneObjectData{}, caskin.ObjectTypeDefault))
	//
	//data2.ObjectID = 2
	//assert.NoError(t, service.ObjectDataUpdateCheck(data2, &example.OneObjectData{}, caskin.ObjectTypeRole))
	//assert.NoError(t, service.DB.Update(data2))
	//
	//// SubAdminUser has no write permission for object_id=2
	//data2.ObjectID = 5
	//provider.User = stage.SubAdminUser
	//assert.Equal(t, caskin.ErrNoWritePermission, service.ObjectDataUpdateCheck(data2, &example.OneObjectData{}, caskin.ObjectTypeRole))
}

func TestServer_ObjectData_Enforce(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.MemberUser
	//data1 := &example.OneObjectData{
	//	Name:     "object_data_1",
	//	ObjectID: 3,
	//}
	//assert.NoError(t, service.ObjectDataCreateCheck(data1, caskin.ObjectTypeDefault))
	//assert.NoError(t, service.EnforceObjectData(data1, caskin.Read))
	//assert.NoError(t, service.EnforceObjectData(data1, caskin.Write))
	//
	//data2 := &example.OneObjectData{
	//	Name:     "object_data_2",
	//	ObjectID: 2,
	//}
	//provider.User = stage.AdminUser
	//assert.NoError(t, service.ObjectDataCreateCheck(data2, caskin.ObjectTypeRole))
	//assert.NoError(t, service.EnforceObjectData(data2, caskin.Read))
	//assert.NoError(t, service.EnforceObjectData(data2, caskin.Write))
	//provider.User = stage.SubAdminUser
	//assert.Equal(t, caskin.ErrNoReadPermission, service.EnforceObjectData(data2, caskin.Read))
	//assert.Equal(t, caskin.ErrNoWritePermission, service.EnforceObjectData(data2, caskin.Write))
}

func TestServer_ObjectData_FilterObjectData(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//data1 := &example.OneObjectData{
	//	Name:     "object_data_1",
	//	ObjectID: 3,
	//}
	//data2 := &example.OneObjectData{
	//	Name:     "object_data_2",
	//	ObjectID: 2,
	//}
	//data3 := &example.OneObjectData{
	//	Name:     "object_data_3",
	//	ObjectID: 5,
	//}
	//list := []any{data1, data2, data3}
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.AdminUser
	//role, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, role, 3)
	//policy, err := service.PolicyGetByRole(role[2])
	//policy = append(policy, &caskin.Policy{
	//	Role: role[2], Object: &example.Object{ID: 3}, Domain: stage.Domain, Action: caskin.Read,
	//})
	//assert.NoError(t, service.PolicyPerRoleModify(role[2], policy))
	//
	//list1, err := service.FilterObjectData(list, caskin.Write)
	//assert.NoError(t, err)
	//assert.Len(t, list1, 3)
	//provider.User = stage.SubAdminUser
	//list2, err := service.FilterObjectData(list, caskin.Read)
	//assert.NoError(t, err)
	//assert.Len(t, list2, 2)
	//list3, err := service.FilterObjectData(list, caskin.Write)
	//assert.NoError(t, err)
	//assert.Len(t, list3, 1)
	//provider.User = stage.MemberUser
	//list4, err := service.FilterObjectData(list, caskin.Write)
	//assert.NoError(t, err)
	//assert.Len(t, list4, 1)
}
