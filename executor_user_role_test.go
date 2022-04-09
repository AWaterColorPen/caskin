package caskin_test

import (
	"testing"
)

func TestServer_UserRole_GetUserRolePair(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.AdminUser
	//list1, err := service.UserRoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, list1, 3)
	//roles, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, roles, 3)
	//
	//pairs := []*caskin.UserRolePair{
	//	{stage.MemberUser, roles[1]},
	//	{stage.SubAdminUser, roles[1]},
	//}
	//assert.NoError(t, service.UserRolePerRoleModify(roles[1], pairs))
	//
	//list3, err := service.UserRoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, list3, 4)
}

func TestServer_UserRole_GetUserRolePairSubAdmin(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.SubAdminUser
	//
	//list, err := service.UserRoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, list, 1)
}

func TestServer_UserRole_GetUserRolePairByRole(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.AdminUser
	//role1 := &example.Role{ID: 2, Name: "xxx"}
	//_, err := service.GetUserRolePairByRole(role1)
	//assert.Equal(t, caskin.ErrNotExists, err)
	//
	//role1.Name = "member"
	//pair1, err := service.GetUserRolePairByRole(role1)
	//assert.NoError(t, err)
	//assert.Len(t, pair1, 1)
	//
	//role2 := &example.Role{Name: "admin"}
	//_, err = service.GetUserRolePairByRole(role2)
	//assert.Equal(t, caskin.ErrEmptyID, err)
	//
	//provider.User = stage.MemberUser
	//_, err = service.GetUserRolePairByRole(role1)
	//assert.Equal(t, caskin.ErrNoReadPermission, err)
	//
	//provider.User = stage.SubAdminUser
	//role3 := &example.Role{ID: 3}
	//pair2, err := service.GetUserRolePairByRole(role3)
	//assert.NoError(t, err)
	//assert.Len(t, pair2, 1)
	//
	//_, err = service.GetUserRolePairByRole(role1)
	//assert.Equal(t, caskin.ErrNoReadPermission, err)
}

func TestServer_UserRole_GetUserRolePairByUser(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.AdminUser
	//user1 := &example.User{ID: 4, PhoneNumber: "xxx"}
	//_, err := service.GetUserRolePairByUser(user1)
	//assert.Equal(t, caskin.ErrNotExists, err)
	//
	//user1.PhoneNumber = stage.SubAdminUser.PhoneNumber
	//pair1, err := service.GetUserRolePairByUser(user1)
	//assert.NoError(t, err)
	//assert.Len(t, pair1, 1)
	//
	//user2 := &example.User{PhoneNumber: stage.SubAdminUser.PhoneNumber}
	//_, err = service.GetUserRolePairByUser(user2)
	//assert.Equal(t, caskin.ErrEmptyID, err)
	//
	//provider.User = stage.MemberUser
	//pair2, err := service.GetUserRolePairByUser(user1)
	//assert.NoError(t, err)
	//assert.Len(t, pair2, 0)
	//
	//provider.User = stage.AdminUser
	//roles, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, roles, 3)
	//pairs := []*caskin.UserRolePair{
	//	{stage.MemberUser, roles[1]},
	//	{stage.SubAdminUser, roles[1]},
	//}
	//assert.NoError(t, service.UserRolePerRoleModify(roles[1], pairs))
	//
	//provider.User = stage.SubAdminUser
	//user3 := &example.User{ID: 4}
	//pair3, err := service.GetUserRolePairByUser(user3)
	//assert.NoError(t, err)
	//assert.Len(t, pair3, 1)
}

func TestServer_UserRole_UserRolePerRoleModify(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.AdminUser
	//objects0, err := service.ObjectGet()
	//assert.NoError(t, err)
	//assert.Len(t, objects0, 5)
	//roles0, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, roles0, 3)
	//
	//provider.User = stage.SubAdminUser
	//roles1, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, roles1, 1)
	//objects1, err := service.ObjectGet()
	//assert.NoError(t, err)
	//assert.Len(t, objects1, 2)
	//
	//pair1 := []*caskin.UserRolePair{
	//	{stage.MemberUser, roles0[0]},
	//	{stage.SubAdminUser, roles0[0]},
	//}
	//assert.Equal(t, caskin.ErrNoWritePermission, service.UserRolePerRoleModify(roles0[0], pair1))
	//
	//pair2 := []*caskin.UserRolePair{
	//	{stage.MemberUser, roles1[0]},
	//	{stage.SubAdminUser, roles1[0]},
	//}
	//assert.NoError(t, service.UserRolePerRoleModify(roles1[0], pair2))
	//
	//provider.User = stage.MemberUser
	//list1, err := service.UserRoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, list1, 2)
	//assert.NoError(t, service.UserRolePerRoleModify(roles1[0], pair2))
	//
	//provider.User = stage.AdminUser
	//pair3 := []*caskin.UserRolePair{
	//	{stage.MemberUser, roles1[0]},
	//	{stage.SubAdminUser, roles1[0]},
	//	{stage.SubAdminUser, roles0[0]},
	//}
	//
	//assert.Equal(t, caskin.ErrInputPairArrayNotBelongSameRole, service.UserRolePerRoleModify(roles1[0], pair3))
}

func TestServer_UserRole_UserRolePerUserModify(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.AdminUser
	//objects0, err := service.ObjectGet()
	//assert.NoError(t, err)
	//assert.Len(t, objects0, 5)
	//roles0, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, roles0, 3)
	//
	//provider.User = stage.SubAdminUser
	//roles1, err := service.RoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, roles1, 1)
	//objects1, err := service.ObjectGet()
	//assert.NoError(t, err)
	//assert.Len(t, objects1, 2)
	//
	//pair1 := []*caskin.UserRolePair{
	//	{stage.MemberUser, roles0[0]},
	//}
	//assert.NoError(t, service.UserRolePerUserModify(stage.MemberUser, pair1))
	//provider.User = stage.AdminUser
	//list1, err := service.UserRoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, list1, 3)
	//
	//provider.User = stage.SubAdminUser
	//pair2 := []*caskin.UserRolePair{
	//	{stage.MemberUser, roles1[0]},
	//}
	//assert.NoError(t, service.UserRolePerUserModify(stage.MemberUser, pair2))
	//list2, err := service.UserRoleGet()
	//assert.NoError(t, err)
	//assert.Len(t, list2, 2)
	//
	//provider.User = stage.AdminUser
	//pair3 := []*caskin.UserRolePair{
	//	{stage.SubAdminUser, roles1[0]},
	//	{stage.MemberUser, roles1[0]},
	//}
	//assert.Equal(t, caskin.ErrInputPairArrayNotBelongSameUser, service.UserRolePerUserModify(stage.SubAdminUser, pair3))
}
