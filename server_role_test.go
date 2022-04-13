package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServer_GetRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles1, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles1, 2)

	// GetRole use action=read, member can read
	roles2, err := service.GetRole(stage.Member, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles2, 2)

	roles3, err := service.GetRole(stage.Superadmin, caskin.GetSuperadminDomain())
	assert.NoError(t, err)
	assert.Len(t, roles3, 0)
}

func TestServer_CreateRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	role1 := &example.Role{Name: "role_01"}
	assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(stage.Member, stage.Domain, role1))
	assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(stage.Admin, stage.Domain, role1))

	roles, _ := service.GetRole(stage.Member, stage.Domain)
	assert.Len(t, roles, 2)
	role1.ObjectID = roles[0].GetObjectID()
	assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(stage.Member, stage.Domain, role1))
	assert.NoError(t, service.CreateRole(stage.Admin, stage.Domain, role1))

	role2 := &example.Role{Name: "role_01"}
	assert.Equal(t, caskin.ErrAlreadyExists, service.CreateRole(stage.Member, stage.Domain, role2))
}

func TestServer_CreateRole_SubNode(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles, _ := service.GetRole(stage.Member, stage.Domain)
	assert.Len(t, roles, 2)

	role1 := &example.Role{
		Name:     "admin-son-1",
		ObjectID: roles[0].GetObjectID(),
		ParentID: roles[0].GetID(),
	}
	// member can not write
	assert.Equal(t, caskin.ErrNoWritePermission, service.CreateRole(stage.Member, stage.Domain, role1))
	//
	// admin can write
	assert.NoError(t, service.CreateRole(stage.Admin, stage.Domain, role1))
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
