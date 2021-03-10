package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutorPolicy_GetPolicyList(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := &example.Provider{}
	executor := stage.Caskin.GetExecutor(provider)

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
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := &example.Provider{}
	executor := stage.Caskin.GetExecutor(provider)

	provider.Domain = stage.Domain
	provider.User = stage.SubAdminUser

	roles, err := executor.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles, 1)
	objects, err := executor.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	list, err := executor.GetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, list, 4)
}

func TestExecutorPolicy_ModifyPolicyListPerRole(t *testing.T) {
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := &example.Provider{}
	executor := stage.Caskin.GetExecutor(provider)

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
	stage, _ := newStage(t)
	assert.NoError(t, stageAddSubAdmin(stage))
	provider := &example.Provider{}
	executor := stage.Caskin.GetExecutor(provider)

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
		{roles1[0], objects1[0], stage.Domain, caskin.Write},
		{roles0[0], objects1[0], stage.Domain, caskin.Read},
		{roles0[0], objects1[0], stage.Domain, caskin.Write},
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
