package caskin_test

import (
	"testing"
)

func TestServer_RoleGet(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	//service := stage.Service
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.AdminUser
	//roles1, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, roles1, 2)
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.MemberUser
	//roles2, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, roles2, 0)
	//
	//provider.Domain = stage.Options.GetSuperadminDomain()
	//provider.User = stage.SuperadminUser
	//roles3, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, roles3, 0)
}

func TestServer_RoleCreate(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//role1 := &example.Role{
	//	Name: "role_01",
	//}
	//assert.Equal(t, caskin.ErrProviderGet, service.RoleCreate(role1))
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.MemberUser
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleCreate(role1))
	//provider.User = stage.AdminUser
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleCreate(role1))
	//role1.ObjectID = 2
	//provider.User = stage.MemberUser
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleCreate(role1))
	//provider.User = stage.AdminUser
	//assert.NoError(t, service.RoleCreate(role1))
	//
	//role2 := &example.Role{
	//	Name: "role_01",
	//}
	//assert.Equal(t, caskin.ErrAlreadyExists, service.RoleCreate(role2))
	//
	//role3 := &example.Role{
	//	Name:     "role_03",
	//	ObjectID: 5,
	//	ParentID: 1,
	//}
	//provider.User = stage.SubAdminUser
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleCreate(role3))
	//role3.ObjectID = 4
	//assert.Equal(t, caskin.ErrInValidObjectType, service.RoleCreate(role3))
}

func TestServer_RoleCreate_SubNode(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//role1 := &example.Role{
	//	Name:     "role_sub_02",
	//	ObjectID: 5,
	//	ParentID: 3,
	//}
	//provider.Domain = stage.Domain
	//provider.User = stage.MemberUser
	//// member can not read or write object5
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleCreate(role1))
	//
	//// subAdmin can read or write object5
	//provider.User = stage.SubAdminUser
	//assert.NoError(t, service.RoleCreate(role1))
	//
	//// make current role a son of member's, subAdminUser does not own the permission
	//role1.ParentID = 2
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleUpdate(role1))
	//
	//role2 := &example.Role{ID: 2}
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleDelete(role2))
	//role3 := &example.Role{ID: 4}
	//assert.NoError(t, service.RoleDelete(role3))
	//
	//provider.User = stage.AdminUser
	//assert.NoError(t, service.RoleDelete(role2))
	//list1, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, list1, 2)
}

func TestServer_RoleUpdate(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//provider.User = stage.AdminUser
	//provider.Domain = stage.Domain
	//service := stage.Service
	//
	//role := &example.Role{
	//	ID: 3, Name: "admin_sub_01_new_name", ParentID: 1, ObjectID: 5}
	//assert.NoError(t, service.RoleUpdate(role))
	//
	//role1 := &example.Role{}
	//assert.Equal(t, caskin.ErrEmptyID, service.RoleUpdate(role1))
	//
	//role2 := &example.Role{ID: 2, Name: "member_new_name", ObjectID: 1, ParentID: 0}
	//assert.Equal(t, caskin.ErrInValidObjectType, service.RoleUpdate(role2))
	//
	//role3 := &example.Role{ID: 2, Name: "member_new_name", ObjectID: 2, ParentID: 2}
	//assert.Equal(t, caskin.ErrParentCanNotBeItself, service.RoleUpdate(role3))
	//
	//provider.User = stage.MemberUser
	//role4 := &example.Role{ID: 3, Name: "admin_sub_01_new_name2", ObjectID: 5, ParentID: 1}
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleUpdate(role4))
	//
	//provider.User = stage.AdminUser
	//assert.NoError(t, service.RoleUpdate(role4))
	//
	//role5 := &example.Role{ID: 3, Name: "admin_sub_01", ObjectID: 4, ParentID: 1}
	//assert.Equal(t, caskin.ErrInValidObjectType, service.RoleUpdate(role5))
	//
	//role6 := &example.Role{ID: 2, Name: "member_new_name", ObjectID: 1, ParentID: 1}
	//assert.Equal(t, caskin.ErrInValidObjectType, service.RoleUpdate(role6))
	//
	//provider.User = stage.SubAdminUser
	//role7 := &example.Role{ID: 3, Name: "member_sub_01", ObjectID: 5, ParentID: 2}
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleUpdate(role7))

}

func TestServer_RoleUpdate_Parent(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//provider.User = stage.AdminUser
	//provider.Domain = stage.Domain
	//service := stage.Service
	//
	//role1 := &example.Role{
	//	Name:     "role_parent_id_to_3",
	//	ObjectID: 5,
	//	ParentID: 3,
	//}
	//assert.NoError(t, service.RoleCreate(role1))
	//role2 := &example.Role{
	//	Name:     "role_parent_id_to_5",
	//	ObjectID: 5,
	//	ParentID: role1.ID,
	//}
	//assert.NoError(t, service.RoleCreate(role2))
	//
	//role3 := &example.Role{
	//	ID:       3,
	//	Name:     "change_role_parent_id_from_2_to_6",
	//	ObjectID: 5,
	//	ParentID: role2.ID,
	//}
	//assert.Equal(t, caskin.ErrParentToDescendant, service.RoleUpdate(role3))
	//role3.ParentID = role1.ID
	//assert.Equal(t, caskin.ErrParentToDescendant, service.RoleUpdate(role3))
	//
	//role2.ParentID = 3
	//assert.NoError(t, service.RoleUpdate(role2))
}

func TestServer_RoleRecover(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//provider.User = stage.AdminUser
	//provider.Domain = stage.Domain
	//service := stage.Service
	//
	//role := &example.Role{
	//	ID:       3,
	//	ParentID: 1,
	//}
	//assert.NoError(t, service.RoleDelete(role))
	//assert.NoError(t, service.RoleRecover(role))
	//assert.Equal(t, caskin.ErrAlreadyExists, service.RoleRecover(role))
	//
	//role1 := &example.Role{}
	//assert.Equal(t, caskin.ErrAlreadyExists, service.RoleRecover(role1))
	//
	//role1 = &example.Role{ID: 2}
	//assert.NoError(t, service.RoleDelete(role1))
	//
	//role2 := &example.Role{ID: 3}
	//assert.NoError(t, service.RoleDelete(role2))
	//
	//provider.User = stage.MemberUser
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleRecover(role1))
	//role2.ID = 3
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleRecover(role2))
	//
	//provider.User = stage.SubAdminUser
	//role3 := &example.Role{ID: 3, ParentID: 1}
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleRecover(role3))
}

func TestServer_RoleDelete(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//provider.User = stage.SubAdminUser
	//provider.Domain = stage.Domain
	//service := stage.Service
	//
	//role := &example.Role{}
	//assert.Equal(t, caskin.ErrEmptyID, service.RoleDelete(role))
	//
	//role.ID = 4
	//assert.Equal(t, caskin.ErrNotExists, service.RoleDelete(role))
	//
	//role1 := &example.Role{ID: 3}
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleDelete(role1))
	//
	//provider.User = stage.AdminUser
	//assert.NoError(t, service.RoleDelete(role1))
	//
	//role4 := &example.Role{ID: 1}
	//assert.NoError(t, service.RoleDelete(role4))
	//
	//role3 := &example.Role{ID: 1}
	//assert.Equal(t, caskin.ErrNoWritePermission, service.RoleRecover(role3))
}

func TestServer_RoleDelete_SubNode(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//provider.User = stage.AdminUser
	//provider.Domain = stage.Domain
	//service := stage.Service
	//
	//role1 := &example.Role{
	//	Name:     "role_sub_02",
	//	ObjectID: 5,
	//	ParentID: 3,
	//}
	//assert.NoError(t, service.RoleCreate(role1))
	//pair1, err := service.UserRoleByUserGet(stage.SubAdminUser)
	//assert.NoError(t, err)
	//assert.Len(t, pair1, 1)
	//
	//pair1 = append(pair1, &caskin.UserRolePair{User: stage.SubAdminUser, Role: role1})
	//assert.NoError(t, service.UserRolePerUserModify(stage.SubAdminUser, pair1))
	//pair2, err := service.UserRoleByUserGet(stage.SubAdminUser)
	//assert.NoError(t, err)
	//assert.Len(t, pair2, 2)
	//
	//role2 := &example.Role{ID: 3}
	//assert.NoError(t, service.RoleDelete(role2))
	//
	//pair3, err := service.UserRoleByUserGet(stage.SubAdminUser)
	//assert.NoError(t, err)
	//assert.Len(t, pair3, 0)
}
