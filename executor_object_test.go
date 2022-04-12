package caskin_test

import (
	"testing"
)

func TestServer_Object_GetObjects(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.AdminUser
	// objects1, err := service.GetObject()
	// assert.NoError(t, err)
	// assert.Len(t, objects1, 3)
	//
	// objects2, err := service.GetObjects(caskin.ObjectTypeObject)
	// assert.NoError(t, err)
	// assert.Len(t, objects2, 1)
	//
	// objects3, err := service.GetObjects(caskin.ObjectTypeObject, caskin.ObjectTypeRole)
	// assert.NoError(t, err)
	// assert.Len(t, objects3, 2)
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.Member
	// objects4, err := service.GetObject()
	// assert.NoError(t, err)
	// assert.Len(t, objects4, 0)
	//
	// provider.Domain = stage.Options.GetSuperadminDomain()
	// provider.User = stage.Superadmin
	// objects5, err := service.GetObject()
	// assert.NoError(t, err)
	// assert.Len(t, objects5, 0)
}

func TestServer_Object_GetExplicitObjects(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.AdminUser
	// objects1, err := service.GetExplicitObjects(caskin.Read, caskin.ObjectTypeObject, caskin.ObjectTypeRole)
	// assert.NoError(t, err)
	// assert.Len(t, objects1, 2)
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.Member
	// objects2, err := service.GetExplicitObjects(caskin.Read)
	// assert.NoError(t, err)
	// assert.Len(t, objects2, 1)
	//
	// provider.Domain = stage.Options.GetSuperadminDomain()
	// provider.User = stage.Superadmin
	// objects3, err := service.GetExplicitObjects(caskin.Read)
	// assert.NoError(t, err)
	// assert.Len(t, objects3, 0)
}

func TestServer_Object_GeneralCreate(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// service := stage.Service
	//
	// object1 := &example.Object{
	//	Name: "object_01",
	//	Type: ObjectTypeTest,
	// }
	// assert.Equal(t, caskin.ErrProviderGet, service.CreateObject(object1))
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.Member
	// assert.Equal(t, caskin.ErrNoWritePermission, service.CreateObject(object1))
	// provider.User = stage.AdminUser
	// assert.Equal(t, caskin.ErrNoWritePermission, service.CreateObject(object1))
	// object1.ObjectID = 1
	// provider.User = stage.Member
	// assert.Equal(t, caskin.ErrNoWritePermission, service.CreateObject(object1))
	// provider.User = stage.AdminUser
	// assert.Equal(t, caskin.ErrEmptyParentIdOrNotSuperadmin, service.CreateObject(object1))
	// object1.ParentID = 3
	// assert.Equal(t, caskin.ErrInValidObjectType, service.CreateObject(object1))
	// object1.Type = caskin.ObjectTypeDefault
	// assert.NoError(t, service.CreateObject(object1))
	//
	// object2 := &example.Object{
	//	Name: "object_01",
	//	Type: caskin.ObjectTypeDefault,
	// }
	// assert.Equal(t, caskin.ErrAlreadyExists, service.CreateObject(object2))
	//
	// object3 := &example.Object{
	//	Name: "object_01",
	//	Type: caskin.ObjectTypeDefault,
	// }
	// assert.Equal(t, caskin.ErrEmptyID, service.DeleteObject(object3))
	// object3.ID = object1.ID
	// assert.NoError(t, service.DeleteObject(object3))
	//
	// object4 := &example.Object{ID: 10, ObjectID: 1}
	// assert.Equal(t, caskin.ErrNotExists, service.DeleteObject(object4))
	// assert.Equal(t, caskin.ErrEmptyParentIdOrNotSuperadmin, service.CreateObject(object4))
}

func TestServer_Object_CreateSubNode(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// object1 := &example.Object{
	//	Name:     "object_sub_02",
	//	Type:     caskin.ObjectTypeObject,
	//	ObjectID: 4,
	//	ParentID: 4,
	// }
	// provider.User = stage.Member
	// provider.Domain = stage.Domain
	// assert.Equal(t, caskin.ErrNoWritePermission, service.CreateObject(object1))
	//
	// provider.User = stage.SubAdminUser
	// assert.NoError(t, service.CreateObject(object1))
	//
	// object1.ParentID = 1
	// assert.Equal(t, caskin.ErrNoWritePermission, service.UpdateObject(object1))
	// object1.ParentID = 5
	// assert.Equal(t, caskin.ErrInValidObjectType, service.UpdateObject(object1))
	// object1.ParentID = 4
	// assert.Equal(t, caskin.ErrObjectTypeObjectIDMustBeItselfID, service.UpdateObject(object1))
	//
	// object1.ObjectID = object1.ID
	// assert.NoError(t, service.UpdateObject(object1))
	//
	// object2 := &example.Object{ID: 4}
	// assert.Equal(t, caskin.ErrNoWritePermission, service.DeleteObject(object2))
	// object3 := &example.Object{ID: 6}
	// assert.NoError(t, service.DeleteObject(object3))
	//
	// provider.User = stage.AdminUser
	// assert.NoError(t, service.DeleteObject(object2))
	// list1, err := service.GetObject()
	// assert.NoError(t, err)
	// assert.Len(t, list1, 3)
	//
	// provider.User = stage.Superadmin
	// list2, err := service.GetObject()
	// assert.NoError(t, err)
	// assert.Len(t, list2, 4)
}

func TestServer_Object_GeneralUpdate(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// provider.User = stage.AdminUser
	// provider.Domain = stage.Domain
	// service := stage.Service
	//
	// object := &example.Object{
	//	ID:       4,
	//	Name:     "object_sub_01_new_name",
	//	Type:     caskin.ObjectTypeObject,
	//	ParentID: 1,
	//	ObjectID: 4,
	// }
	// assert.NoError(t, service.UpdateObject(object))
	//
	// object.Type = caskin.ObjectTypeRole
	// assert.Equal(t, caskin.ErrCantChangeObjectType, service.UpdateObject(object))
	//
	// object1 := &example.Object{}
	// assert.Equal(t, caskin.ErrEmptyID, service.UpdateObject(object1))
	//
	// object2 := &example.Object{ID: 1, Name: "object_01_new_name", Type: caskin.ObjectTypeObject, ObjectID: 1, ParentID: 0}
	// assert.Equal(t, caskin.ErrEmptyParentIdOrNotSuperadmin, service.UpdateObject(object2))
	//
	// provider.User = stage.Member
	// object4 := &example.Object{
	//	ID:       4,
	//	Name:     "object_sub_01_new_name2",
	//	Type:     caskin.ObjectTypeObject,
	//	ParentID: 1,
	//	ObjectID: 4,
	// }
	// assert.Equal(t, caskin.ErrNoWritePermission, service.UpdateObject(object4))
	//
	// provider.User = stage.Superadmin
	// object5 := &example.Object{
	//	ID:       1,
	//	Type:     caskin.ObjectTypeDefault,
	//	ParentID: 1,
	//	ObjectID: 1,
	// }
	// assert.Equal(t, caskin.ErrParentCanNotBeItself, service.UpdateObject(object5))
	// object5.ParentID = 0
	// assert.Equal(t, caskin.ErrCantChangeObjectType, service.UpdateObject(object5))
	//
	// object6 := &example.Object{
	//	ID:       4,
	//	Type:     caskin.ObjectTypeRole,
	//	ParentID: 2,
	//	ObjectID: 4,
	// }
	// provider.User = stage.AdminUser
	// assert.Equal(t, caskin.ErrCantChangeObjectType, service.UpdateObject(object6))
}

func TestServer_Object_UpdateParent(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// provider.User = stage.AdminUser
	// provider.Domain = stage.Domain
	// service := stage.Service
	//
	// object1 := &example.Object{
	//	Name:     "role_parent_id_to_5",
	//	Type:     caskin.ObjectTypeRole,
	//	ParentID: 5,
	//	ObjectID: 4,
	// }
	// assert.NoError(t, service.CreateObject(object1))
	// object2 := &example.Object{
	//	Name:     "role_parent_id_to_6",
	//	Type:     caskin.ObjectTypeRole,
	//	ParentID: object1.ID,
	//	ObjectID: 4,
	// }
	// assert.NoError(t, service.CreateObject(object2))
	//
	// object3 := &example.Object{
	//	ID:       5,
	//	Name:     "change_role_parent_id_from_2_to_7",
	//	Type:     caskin.ObjectTypeRole,
	//	ParentID: object2.ID,
	//	ObjectID: 4,
	// }
	// assert.Equal(t, caskin.ErrParentToDescendant, service.UpdateObject(object3))
	// object3.ParentID = object1.ID
	// assert.Equal(t, caskin.ErrParentToDescendant, service.UpdateObject(object3))
	//
	// object2.ParentID = 5
	// assert.NoError(t, service.UpdateObject(object2))
}

func TestServer_Object_GeneralRecover(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// provider.User = stage.AdminUser
	// provider.Domain = stage.Domain
	// service := stage.Service
	//
	// object := &example.Object{
	//	ID:       4,
	//	ParentID: 1,
	// }
	// assert.NoError(t, service.DeleteObject(object))
	//
	// object1 := &example.Object{}
	// assert.Equal(t, caskin.ErrAlreadyExists, service.RecoverObject(object1))
	//
	// provider.User = stage.Superadmin
	// assert.NoError(t, service.RecoverObject(object))
	// assert.Equal(t, caskin.ErrAlreadyExists, service.RecoverObject(object))
	//
	// object2 := &example.Object{ID: 3}
	// assert.NoError(t, service.DeleteObject(object2))
	// provider.User = stage.Member
	// object2.ID = 3
	// assert.Equal(t, caskin.ErrNoWritePermission, service.RecoverObject(object2))
}

func TestServer_Object_GeneralDelete(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// provider.User = stage.SubAdminUser
	// provider.Domain = stage.Domain
	// service := stage.Service
	//
	// object := &example.Object{}
	// assert.Equal(t, caskin.ErrEmptyID, service.DeleteObject(object))
	//
	// object1 := &example.Object{ID: 4}
	// assert.Equal(t, caskin.ErrNoWritePermission, service.DeleteObject(object1))
	//
	// provider.User = stage.AdminUser
	// assert.NoError(t, service.DeleteObject(object1))
	//
	// object2 := &example.Object{ID: 5}
	// assert.Equal(t, caskin.ErrNoWritePermission, service.DeleteObject(object2))
	//
	// object4 := &example.Object{ID: 1}
	// assert.Equal(t, caskin.ErrEmptyParentIdOrNotSuperadmin, service.DeleteObject(object4))
	//
	// // TODO issue 2: if object type == object, it will not recover it by admin user now. we want to support it by a special API
	// // want to support it by a special API
	// object3 := &example.Object{ID: 4, ParentID: 1}
	// assert.Equal(t, caskin.ErrNoWritePermission, service.RecoverObject(object3))
	// provider.User = stage.Superadmin
	// assert.NoError(t, service.RecoverObject(object3))
}
