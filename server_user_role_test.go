package caskin_test

import (
	"github.com/awatercolorpen/caskin/example"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestServer_UserRole_GetUserRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	list1, err := service.GetUserRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list1, 2)
	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	pairs := []*caskin.UserRolePair{
		{stage.Admin, roles[0]},
		{stage.Member, roles[0]},
	}
	assert.NoError(t, service.ModifyUserRolePerRole(stage.Admin, stage.Domain, roles[0], pairs))

	list2, err := service.GetUserRole(stage.Member, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list2, 3)
}

func TestServer_UserRole_GetUserRole_TreeNode(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	object1 := &example.Object{Name: "role_sub_01", Type: "role"}
	object1.ParentID = objects[0].GetID()
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	role1 := &example.Role{Name: "role_sub_01", ObjectID: object1.ID}
	assert.NoError(t, service.CreateRole(stage.Admin, stage.Domain, role1))
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.SubAdminUser
	//
	// list, err := service.GetUserRole()
	// assert.NoError(t, err)
	// assert.Len(t, list, 1)
}

func TestServer_UserRole_GetUserRolePairByRole(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.AdminUser
	// role1 := &example.Role{ID: 2, Name: "xxx"}
	// _, err := service.GetUserRolePairByRole(role1)
	// assert.Equal(t, caskin.ErrNotExists, err)
	//
	// role1.Name = "member"
	// pair1, err := service.GetUserRolePairByRole(role1)
	// assert.NoError(t, err)
	// assert.Len(t, pair1, 1)
	//
	// role2 := &example.Role{Name: "admin"}
	// _, err = service.GetUserRolePairByRole(role2)
	// assert.Equal(t, caskin.ErrEmptyID, err)
	//
	// provider.User = stage.Member
	// _, err = service.GetUserRolePairByRole(role1)
	// assert.Equal(t, caskin.ErrNoReadPermission, err)
	//
	// provider.User = stage.SubAdminUser
	// role3 := &example.Role{ID: 3}
	// pair2, err := service.GetUserRolePairByRole(role3)
	// assert.NoError(t, err)
	// assert.Len(t, pair2, 1)
	//
	// _, err = service.GetUserRolePairByRole(role1)
	// assert.Equal(t, caskin.ErrNoReadPermission, err)
}

func TestServer_UserRole_GetUserRolePairByUser(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.AdminUser
	// user1 := &example.User{ID: 4, PhoneNumber: "xxx"}
	// _, err := service.GetUserRolePairByUser(user1)
	// assert.Equal(t, caskin.ErrNotExists, err)
	//
	// user1.PhoneNumber = stage.SubAdminUser.PhoneNumber
	// pair1, err := service.GetUserRolePairByUser(user1)
	// assert.NoError(t, err)
	// assert.Len(t, pair1, 1)
	//
	// user2 := &example.User{PhoneNumber: stage.SubAdminUser.PhoneNumber}
	// _, err = service.GetUserRolePairByUser(user2)
	// assert.Equal(t, caskin.ErrEmptyID, err)
	//
	// provider.User = stage.Member
	// pair2, err := service.GetUserRolePairByUser(user1)
	// assert.NoError(t, err)
	// assert.Len(t, pair2, 0)
	//
	// provider.User = stage.AdminUser
	// roles, err := service.GetRole()
	// assert.NoError(t, err)
	// assert.Len(t, roles, 3)
	// pairs := []*caskin.UserRolePair{
	//	{stage.Member, roles[1]},
	//	{stage.SubAdminUser, roles[1]},
	// }
	// assert.NoError(t, service.ModifyUserRolePerRole(roles[1], pairs))
	//
	// provider.User = stage.SubAdminUser
	// user3 := &example.User{ID: 4}
	// pair3, err := service.GetUserRolePairByUser(user3)
	// assert.NoError(t, err)
	// assert.Len(t, pair3, 1)
}

func TestServer_UserRole_UserRolePerRoleModify(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.AdminUser
	// objects0, err := service.GetObjectID()
	// assert.NoError(t, err)
	// assert.Len(t, objects0, 5)
	// roles0, err := service.GetRole()
	// assert.NoError(t, err)
	// assert.Len(t, roles0, 3)
	//
	// provider.User = stage.SubAdminUser
	// roles1, err := service.GetRole()
	// assert.NoError(t, err)
	// assert.Len(t, roles1, 1)
	// objects1, err := service.GetObjectID()
	// assert.NoError(t, err)
	// assert.Len(t, objects1, 2)
	//
	// pair1 := []*caskin.UserRolePair{
	//	{stage.Member, roles0[0]},
	//	{stage.SubAdminUser, roles0[0]},
	// }
	// assert.Equal(t, caskin.ErrNoWritePermission, service.ModifyUserRolePerRole(roles0[0], pair1))
	//
	// pair2 := []*caskin.UserRolePair{
	//	{stage.Member, roles1[0]},
	//	{stage.SubAdminUser, roles1[0]},
	// }
	// assert.NoError(t, service.ModifyUserRolePerRole(roles1[0], pair2))
	//
	// provider.User = stage.Member
	// list1, err := service.GetUserRole()
	// assert.NoError(t, err)
	// assert.Len(t, list1, 2)
	// assert.NoError(t, service.ModifyUserRolePerRole(roles1[0], pair2))
	//
	// provider.User = stage.AdminUser
	// pair3 := []*caskin.UserRolePair{
	//	{stage.Member, roles1[0]},
	//	{stage.SubAdminUser, roles1[0]},
	//	{stage.SubAdminUser, roles0[0]},
	// }
	//
	// assert.Equal(t, caskin.ErrInputPairArrayNotBelongSameRole, service.ModifyUserRolePerRole(roles1[0], pair3))
}

func TestServer_UserRole_UserRolePerUserModify(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.AdminUser
	// objects0, err := service.GetObjectID()
	// assert.NoError(t, err)
	// assert.Len(t, objects0, 5)
	// roles0, err := service.GetRole()
	// assert.NoError(t, err)
	// assert.Len(t, roles0, 3)
	//
	// provider.User = stage.SubAdminUser
	// roles1, err := service.GetRole()
	// assert.NoError(t, err)
	// assert.Len(t, roles1, 1)
	// objects1, err := service.GetObjectID()
	// assert.NoError(t, err)
	// assert.Len(t, objects1, 2)
	//
	// pair1 := []*caskin.UserRolePair{
	//	{stage.Member, roles0[0]},
	// }
	// assert.NoError(t, service.ModifyUserRolePerUser(stage.Member, pair1))
	// provider.User = stage.AdminUser
	// list1, err := service.GetUserRole()
	// assert.NoError(t, err)
	// assert.Len(t, list1, 3)
	//
	// provider.User = stage.SubAdminUser
	// pair2 := []*caskin.UserRolePair{
	//	{stage.Member, roles1[0]},
	// }
	// assert.NoError(t, service.ModifyUserRolePerUser(stage.Member, pair2))
	// list2, err := service.GetUserRole()
	// assert.NoError(t, err)
	// assert.Len(t, list2, 2)
	//
	// provider.User = stage.AdminUser
	// pair3 := []*caskin.UserRolePair{
	//	{stage.SubAdminUser, roles1[0]},
	//	{stage.Member, roles1[0]},
	// }
	// assert.Equal(t, caskin.ErrInputPairArrayNotBelongSameUser, service.ModifyUserRolePerUser(stage.SubAdminUser, pair3))
}