package caskin_test

import (
	"testing"
)

func TestServer_RoleGet(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.AdminUser
	// roles1, err := service.GetRole()
	// assert.NoError(t, err)
	// assert.Len(t, roles1, 2)
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.Member
	// roles2, err := service.GetRole()
	// assert.NoError(t, err)
	// assert.Len(t, roles2, 0)
	//
	// provider.Domain = stage.Options.GetSuperadminDomain()
	// provider.User = stage.Superadmin
	// roles3, err := service.GetRole()
	// assert.NoError(t, err)
	// assert.Len(t, roles3, 0)
}

func TestServer_RoleCreate(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// role1 := &example.Role{
	//	Name: "role_01",
	// }
	// assert.Equal(t, caskin.ErrProviderGet, service.CreateRole(role1))
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.Member
	// assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(role1))
	// provider.User = stage.AdminUser
	// assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(role1))
	// role1.ObjectID = 2
	// provider.User = stage.Member
	// assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(role1))
	// provider.User = stage.AdminUser
	// assert.NoError(t, service.CreateRole(role1))
	//
	// role2 := &example.Role{
	//	Name: "role_01",
	// }
	// assert.Equal(t, caskin.ErrAlreadyExists, service.CreateRole(role2))
	//
	// role3 := &example.Role{
	//	Name:     "role_03",
	//	ObjectID: 5,
	//	ParentID: 1,
	// }
	// provider.User = stage.SubAdminUser
	// assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(role3))
	// role3.ObjectID = 4
	// assert.Equal(t, caskin.ErrInValidObjectType, service.CreateRole(role3))
}

func TestServer_RoleCreate_SubNode(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// role1 := &example.Role{
	//	Name:     "role_sub_02",
	//	ObjectID: 5,
	//	ParentID: 3,
	// }
	// provider.Domain = stage.Domain
	// provider.User = stage.Member
	// // member can not read or write object5
	// assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(role1))
	//
	// // subAdmin can read or write object5
	// provider.User = stage.SubAdminUser
	// assert.NoError(t, service.CreateRole(role1))
	//
	// // make current role a son of member's, subAdminUser does not own the permission
	// role1.ParentID = 2
	// assert.Equal(t, caskin.ErrNoWritePermission, service.UpdateRole(role1))
	//
	// role2 := &example.Role{ID: 2}
	// assert.Equal(t, caskin.ErrNoWritePermission, service.DeleteRole(role2))
	// role3 := &example.Role{ID: 4}
	// assert.NoError(t, service.DeleteRole(role3))
	//
	// provider.User = stage.AdminUser
	// assert.NoError(t, service.DeleteRole(role2))
	// list1, err := service.GetRole()
	// assert.NoError(t, err)
	// assert.Len(t, list1, 2)
}

func TestServer_RoleUpdate(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// provider.User = stage.AdminUser
	// provider.Domain = stage.Domain
	// service := stage.Service
	//
	// role := &example.Role{
	//	ID: 3, Name: "admin_sub_01_new_name", ParentID: 1, ObjectID: 5}
	// assert.NoError(t, service.UpdateRole(role))
	//
	// role1 := &example.Role{}
	// assert.Equal(t, caskin.ErrEmptyID, service.UpdateRole(role1))
	//
	// role2 := &example.Role{ID: 2, Name: "member_new_name", ObjectID: 1, ParentID: 0}
	// assert.Equal(t, caskin.ErrInValidObjectType, service.UpdateRole(role2))
	//
	// role3 := &example.Role{ID: 2, Name: "member_new_name", ObjectID: 2, ParentID: 2}
	// assert.Equal(t, caskin.ErrParentCanNotBeItself, service.UpdateRole(role3))
	//
	// provider.User = stage.Member
	// role4 := &example.Role{ID: 3, Name: "admin_sub_01_new_name2", ObjectID: 5, ParentID: 1}
	// assert.Equal(t, caskin.ErrNoWritePermission, service.UpdateRole(role4))
	//
	// provider.User = stage.AdminUser
	// assert.NoError(t, service.UpdateRole(role4))
	//
	// role5 := &example.Role{ID: 3, Name: "admin_sub_01", ObjectID: 4, ParentID: 1}
	// assert.Equal(t, caskin.ErrInValidObjectType, service.UpdateRole(role5))
	//
	// role6 := &example.Role{ID: 2, Name: "member_new_name", ObjectID: 1, ParentID: 1}
	// assert.Equal(t, caskin.ErrInValidObjectType, service.UpdateRole(role6))
	//
	// provider.User = stage.SubAdminUser
	// role7 := &example.Role{ID: 3, Name: "member_sub_01", ObjectID: 5, ParentID: 2}
	// assert.Equal(t, caskin.ErrNoWritePermission, service.UpdateRole(role7))

}

func TestServer_RoleUpdate_Parent(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// provider.User = stage.AdminUser
	// provider.Domain = stage.Domain
	// service := stage.Service
	//
	// role1 := &example.Role{
	//	Name:     "role_parent_id_to_3",
	//	ObjectID: 5,
	//	ParentID: 3,
	// }
	// assert.NoError(t, service.CreateRole(role1))
	// role2 := &example.Role{
	//	Name:     "role_parent_id_to_5",
	//	ObjectID: 5,
	//	ParentID: role1.ID,
	// }
	// assert.NoError(t, service.CreateRole(role2))
	//
	// role3 := &example.Role{
	//	ID:       3,
	//	Name:     "change_role_parent_id_from_2_to_6",
	//	ObjectID: 5,
	//	ParentID: role2.ID,
	// }
	// assert.Equal(t, caskin.ErrParentToDescendant, service.UpdateRole(role3))
	// role3.ParentID = role1.ID
	// assert.Equal(t, caskin.ErrParentToDescendant, service.UpdateRole(role3))
	//
	// role2.ParentID = 3
	// assert.NoError(t, service.UpdateRole(role2))
}

func TestServer_RoleRecover(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// provider.User = stage.AdminUser
	// provider.Domain = stage.Domain
	// service := stage.Service
	//
	// role := &example.Role{
	//	ID:       3,
	//	ParentID: 1,
	// }
	// assert.NoError(t, service.DeleteRole(role))
	// assert.NoError(t, service.RecoverRole(role))
	// assert.Equal(t, caskin.ErrAlreadyExists, service.RecoverRole(role))
	//
	// role1 := &example.Role{}
	// assert.Equal(t, caskin.ErrAlreadyExists, service.RecoverRole(role1))
	//
	// role1 = &example.Role{ID: 2}
	// assert.NoError(t, service.DeleteRole(role1))
	//
	// role2 := &example.Role{ID: 3}
	// assert.NoError(t, service.DeleteRole(role2))
	//
	// provider.User = stage.Member
	// assert.Equal(t, caskin.ErrNoWritePermission, service.RecoverRole(role1))
	// role2.ID = 3
	// assert.Equal(t, caskin.ErrNoWritePermission, service.RecoverRole(role2))
	//
	// provider.User = stage.SubAdminUser
	// role3 := &example.Role{ID: 3, ParentID: 1}
	// assert.Equal(t, caskin.ErrNoWritePermission, service.RecoverRole(role3))
}

func TestServer_RoleDelete(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// provider.User = stage.SubAdminUser
	// provider.Domain = stage.Domain
	// service := stage.Service
	//
	// role := &example.Role{}
	// assert.Equal(t, caskin.ErrEmptyID, service.DeleteRole(role))
	//
	// role.ID = 4
	// assert.Equal(t, caskin.ErrNotExists, service.DeleteRole(role))
	//
	// role1 := &example.Role{ID: 3}
	// assert.Equal(t, caskin.ErrNoWritePermission, service.DeleteRole(role1))
	//
	// provider.User = stage.AdminUser
	// assert.NoError(t, service.DeleteRole(role1))
	//
	// role4 := &example.Role{ID: 1}
	// assert.NoError(t, service.DeleteRole(role4))
	//
	// role3 := &example.Role{ID: 1}
	// assert.Equal(t, caskin.ErrNoWritePermission, service.RecoverRole(role3))
}

func TestServer_RoleDelete_SubNode(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// provider.User = stage.AdminUser
	// provider.Domain = stage.Domain
	// service := stage.Service
	//
	// role1 := &example.Role{
	//	Name:     "role_sub_02",
	//	ObjectID: 5,
	//	ParentID: 3,
	// }
	// assert.NoError(t, service.CreateRole(role1))
	// pair1, err := service.GetUserRoleByUser(stage.SubAdminUser)
	// assert.NoError(t, err)
	// assert.Len(t, pair1, 1)
	//
	// pair1 = append(pair1, &caskin.UserRolePair{User: stage.SubAdminUser, Role: role1})
	// assert.NoError(t, service.ModifyUserRolePerUser(stage.SubAdminUser, pair1))
	// pair2, err := service.GetUserRoleByUser(stage.SubAdminUser)
	// assert.NoError(t, err)
	// assert.Len(t, pair2, 2)
	//
	// role2 := &example.Role{ID: 3}
	// assert.NoError(t, service.DeleteRole(role2))
	//
	// pair3, err := service.GetUserRoleByUser(stage.SubAdminUser)
	// assert.NoError(t, err)
	// assert.Len(t, pair3, 0)
}
