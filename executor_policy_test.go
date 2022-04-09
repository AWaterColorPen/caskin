package caskin_test

import (
	"github.com/awatercolorpen/caskin/playground"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorPolicy_GetPolicyList(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	
	service := stage.Service

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	list1, err := executor.GetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, list1, 12)
	roles, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles, 3)
	objects, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects, 5)

	domain := stage.Domain
	policy1 := []*caskin.Policy{
		{roles[0], objects[0], domain, caskin.Read},
		{roles[0], objects[1], domain, caskin.Read},
		{roles[0], objects[1], domain, caskin.Write},
		{roles[0], objects[2], domain, caskin.Read},
		{roles[0], objects[2], domain, caskin.Write},
	}
	assert.NoError(t, executor.ModifyPolicyListPerRole(roles[0], policy1))

	list3, err := executor.GetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, list3, 11)
}

func TestExecutorPolicy_GetPolicyListFromSubAdmin(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	
	service := stage.Service

	provider.Domain = stage.Domain
	provider.User = stage.SubAdminUser

	list, err := executor.GetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, list, 4)
}

func TestExecutorPolicy_GetPolicyListByRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	
	service := stage.Service

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	role1 := &example.Role{ID: 2, Name: "xxx"}
	_, err := executor.GetPolicyListByRole(role1)
	assert.Equal(t, caskin.ErrNotExists, err)

	role1.Name = "member"
	policy1, err := executor.GetPolicyListByRole(role1)
	assert.NoError(t, err)
	assert.Len(t, policy1, 2)

	role2 := &example.Role{Name: "admin"}
	_, err = executor.GetPolicyListByRole(role2)
	assert.Equal(t, caskin.ErrEmptyID, err)

	provider.User = stage.MemberUser
	_, err = executor.GetPolicyListByRole(role1)
	assert.Equal(t, caskin.ErrNoReadPermission, err)

	provider.User = stage.SubAdminUser
	role3 := &example.Role{ID: 3}
	policy2, err := executor.GetPolicyListByRole(role3)
	assert.NoError(t, err)
	assert.Len(t, policy2, 4)

	_, err = executor.GetPolicyListByRole(role1)
	assert.Equal(t, caskin.ErrNoReadPermission, err)
}

func TestExecutorPolicy_GetPolicyListByObject(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	
	service := stage.Service

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	object1 := &example.Object{ID: 2, Name: "xxx"}
	_, err := executor.GetPolicyListByObject(object1)
	assert.Equal(t, caskin.ErrNotExists, err)

	object1.Name = string(caskin.ObjectTypeRole)
	policy1, err := executor.GetPolicyListByObject(object1)
	assert.NoError(t, err)
	assert.Len(t, policy1, 2)

	object2 := &example.Object{Name: "object"}
	_, err = executor.GetPolicyListByObject(object2)
	assert.Equal(t, caskin.ErrEmptyID, err)

	provider.User = stage.MemberUser
	_, err = executor.GetPolicyListByObject(object1)
	assert.Equal(t, caskin.ErrNoReadPermission, err)

	provider.User = stage.SubAdminUser
	object3 := &example.Object{ID: 4}
	policy2, err := executor.GetPolicyListByObject(object3)
	assert.NoError(t, err)
	assert.Len(t, policy2, 2)

	_, err = executor.GetPolicyListByObject(object1)
	assert.Equal(t, caskin.ErrNoReadPermission, err)
}

func TestExecutorPolicy_ModifyPolicyListPerRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	
	service := stage.Service

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	objects0, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects0, 5)
	roles0, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles0, 3)

	provider.User = stage.SubAdminUser
	roles1, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles1, 1)
	objects1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects1, 2)

	policy1 := []*caskin.Policy{
		{roles0[0], objects1[0], stage.Domain, caskin.Read},
	}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.ModifyPolicyListPerRole(roles0[0], policy1))

	policy2 := []*caskin.Policy{
		{roles1[0], objects1[0], stage.Domain, caskin.Read},
		{roles1[0], objects1[1], stage.Domain, caskin.Read},
		{roles1[0], objects0[0], stage.Domain, caskin.Read},
		{roles1[0], objects0[0], stage.Domain, caskin.Write},
	}
	assert.NoError(t, executor.ModifyPolicyListPerRole(roles1[0], policy2))
	list1, err := executor.GetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, list1, 2)

	provider.User = stage.AdminUser
	policy3 := []*caskin.Policy{
		{roles0[0], objects1[0], stage.Domain, caskin.Read},
		{roles0[0], objects1[1], stage.Domain, caskin.Read},
		{roles1[0], objects0[0], stage.Domain, caskin.Read},
		{roles1[0], objects0[0], stage.Domain, caskin.Write},
	}
	assert.Equal(t, caskin.ErrInputPolicyListNotBelongSameRole, executor.ModifyPolicyListPerRole(roles1[0], policy3))
}

func TestExecutorPolicy_ModifyPolicyListPerObject(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	assert.NoError(t, stage.AddSubAdmin())
	
	service := stage.Service

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	objects0, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects0, 5)
	roles0, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles0, 3)

	provider.User = stage.SubAdminUser
	roles1, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles1, 1)
	objects1, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects1, 2)

	policy1 := []*caskin.Policy{
		{roles1[0], objects0[0], stage.Domain, caskin.Read},
	}
	assert.Equal(t, caskin.ErrNoWritePermission, executor.ModifyPolicyListPerObject(objects0[0], policy1))

	policy2 := []*caskin.Policy{
		{roles1[0], objects1[0], stage.Domain, caskin.Read},
		{roles0[0], objects1[0], stage.Domain, caskin.Read},
		{roles0[0], objects1[0], stage.Domain, caskin.Write},
	}
	assert.NoError(t, executor.ModifyPolicyListPerObject(objects1[0], policy2))
	list1, err := executor.GetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, list1, 3)

	provider.User = stage.AdminUser
	policy3 := []*caskin.Policy{
		{roles0[0], objects1[0], stage.Domain, caskin.Read},
		{roles0[0], objects1[1], stage.Domain, caskin.Read},
		{roles1[0], objects0[0], stage.Domain, caskin.Read},
		{roles1[0], objects0[0], stage.Domain, caskin.Write},
	}
	assert.Equal(t, caskin.ErrInputPolicyListNotBelongSameObject, executor.ModifyPolicyListPerObject(objects1[0], policy3))
}
