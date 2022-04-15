package caskin_test

import (
	"testing"
)

func TestServer_Policy_GetPolicyList(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.AdminUser
	// list1, err := service.GetPolicy()
	// assert.NoError(t, err)
	// assert.Len(t, list1, 12)
	// roles, err := service.GetRole()
	// assert.NoError(t, err)
	// assert.Len(t, roles, 3)
	// objects, err := service.GetObjectID()
	// assert.NoError(t, err)
	// assert.Len(t, objects, 5)
	//
	// domain := stage.Domain
	// policy1 := []*caskin.Policy{
	//	{roles[0], objects[0], domain, caskin.Read},
	//	{roles[0], objects[1], domain, caskin.Read},
	//	{roles[0], objects[1], domain, caskin.Write},
	//	{roles[0], objects[2], domain, caskin.Read},
	//	{roles[0], objects[2], domain, caskin.Write},
	// }
	// assert.NoError(t, service.ModifyPolicyPerRole(roles[0], policy1))
	//
	// list3, err := service.GetPolicy()
	// assert.NoError(t, err)
	// assert.Len(t, list3, 11)
}

func TestServer_Policy_GetPolicyListFromSubAdmin(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.SubAdminUser
	//
	// list, err := service.GetPolicy()
	// assert.NoError(t, err)
	// assert.Len(t, list, 4)
}

func TestServer_Policy_GetPolicyListByRole(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.AdminUser
	// role1 := &example.Role{ID: 2, Name: "xxx"}
	// _, err := service.PolicyGetByRole(role1)
	// assert.Equal(t, caskin.ErrNotExists, err)
	//
	// role1.Name = "member"
	// policy1, err := service.PolicyGetByRole(role1)
	// assert.NoError(t, err)
	// assert.Len(t, policy1, 2)
	//
	// role2 := &example.Role{Name: "admin"}
	// _, err = service.PolicyGetByRole(role2)
	// assert.Equal(t, caskin.ErrEmptyID, err)
	//
	// provider.User = stage.Member
	// _, err = service.PolicyGetByRole(role1)
	// assert.Equal(t, caskin.ErrNoReadPermission, err)
	//
	// provider.User = stage.SubAdminUser
	// role3 := &example.Role{ID: 3}
	// policy2, err := service.PolicyGetByRole(role3)
	// assert.NoError(t, err)
	// assert.Len(t, policy2, 4)
	//
	// _, err = service.PolicyGetByRole(role1)
	// assert.Equal(t, caskin.ErrNoReadPermission, err)
}

func TestServer_Policy_GetPolicyListByObject(t *testing.T) {
	// stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	// assert.NoError(t, stage.AddSubAdmin())
	//
	// service := stage.Service
	//
	// provider.Domain = stage.Domain
	// provider.User = stage.AdminUser
	// object1 := &example.Object{ID: 2, Name: "xxx"}
	// _, err := service.PolicyGetByObject(object1)
	// assert.Equal(t, caskin.ErrNotExists, err)
	//
	// object1.Name = string(caskin.ObjectTypeRole)
	// policy1, err := service.PolicyGetByObject(object1)
	// assert.NoError(t, err)
	// assert.Len(t, policy1, 2)
	//
	// object2 := &example.Object{Name: "object"}
	// _, err = service.PolicyGetByObject(object2)
	// assert.Equal(t, caskin.ErrEmptyID, err)
	//
	// provider.User = stage.Member
	// _, err = service.PolicyGetByObject(object1)
	// assert.Equal(t, caskin.ErrNoReadPermission, err)
	//
	// provider.User = stage.SubAdminUser
	// object3 := &example.Object{ID: 4}
	// policy2, err := service.PolicyGetByObject(object3)
	// assert.NoError(t, err)
	// assert.Len(t, policy2, 2)
	//
	// _, err = service.PolicyGetByObject(object1)
	// assert.Equal(t, caskin.ErrNoReadPermission, err)
}

func TestServer_Policy_ModifyPolicyListPerRole(t *testing.T) {
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
	// policy1 := []*caskin.Policy{
	//	{roles0[0], objects1[0], stage.Domain, caskin.Read},
	// }
	// assert.Equal(t, caskin.ErrNoWritePermission, service.ModifyPolicyPerRole(roles0[0], policy1))
	//
	// policy2 := []*caskin.Policy{
	//	{roles1[0], objects1[0], stage.Domain, caskin.Read},
	//	{roles1[0], objects1[1], stage.Domain, caskin.Read},
	//	{roles1[0], objects0[0], stage.Domain, caskin.Read},
	//	{roles1[0], objects0[0], stage.Domain, caskin.Write},
	// }
	// assert.NoError(t, service.ModifyPolicyPerRole(roles1[0], policy2))
	// list1, err := service.GetPolicy()
	// assert.NoError(t, err)
	// assert.Len(t, list1, 2)
	//
	// provider.User = stage.AdminUser
	// policy3 := []*caskin.Policy{
	//	{roles0[0], objects1[0], stage.Domain, caskin.Read},
	//	{roles0[0], objects1[1], stage.Domain, caskin.Read},
	//	{roles1[0], objects0[0], stage.Domain, caskin.Read},
	//	{roles1[0], objects0[0], stage.Domain, caskin.Write},
	// }
	// assert.Equal(t, caskin.ErrInputPolicyListNotBelongSameRole, service.ModifyPolicyPerRole(roles1[0], policy3))
}

func TestServer_Policy_ModifyPolicyListPerObject(t *testing.T) {
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
	// policy1 := []*caskin.Policy{
	//	{roles1[0], objects0[0], stage.Domain, caskin.Read},
	// }
	// assert.Equal(t, caskin.ErrNoWritePermission, service.ModifyPolicyListPerObject(objects0[0], policy1))
	//
	// policy2 := []*caskin.Policy{
	//	{roles1[0], objects1[0], stage.Domain, caskin.Read},
	//	{roles0[0], objects1[0], stage.Domain, caskin.Read},
	//	{roles0[0], objects1[0], stage.Domain, caskin.Write},
	// }
	// assert.NoError(t, service.ModifyPolicyListPerObject(objects1[0], policy2))
	// list1, err := service.GetPolicy()
	// assert.NoError(t, err)
	// assert.Len(t, list1, 3)
	//
	// provider.User = stage.AdminUser
	// policy3 := []*caskin.Policy{
	//	{roles0[0], objects1[0], stage.Domain, caskin.Read},
	//	{roles0[0], objects1[1], stage.Domain, caskin.Read},
	//	{roles1[0], objects0[0], stage.Domain, caskin.Read},
	//	{roles1[0], objects0[0], stage.Domain, caskin.Write},
	// }
	// assert.Equal(t, caskin.ErrInputPolicyListNotBelongSameObject, service.ModifyPolicyListPerObject(objects1[0], policy3))
}
