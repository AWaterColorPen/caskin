package web_feature_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorNormalDomain_GetFeature(t *testing.T) {
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	assert.Error(t, executor.SyncLatestVersionToAllDomain())

	assert.NoError(t, executor.BuildVersion())
	assert.NoError(t, executor.SyncLatestVersionToAllDomain())
	assert.NoError(t, reinitializeDomainWithWebFeature(stage, w.GetRoot()))

	executor1 := stage.Caskin.GetExecutor(provider)
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	roles, err := executor1.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles, 3)

	policy1, err := executor1.GetPolicyListByRole(roles[2])
	assert.NoError(t, err)
	assert.Len(t, policy1, 4)
	input1 := []*caskin.Policy{
		policy1[0], policy1[1], policy1[2], policy1[3],
		{roles[2], &example.Object{ID: featureStartID}, stage.Domain, caskin.Read},
	}

	assert.NoError(t, executor.NormalDomainModifyPolicyListPerRole(roles[2], input1))
	objects1, err := executor1.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects1, 5)
	objects2, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, objects2, 5)

	provider.User = stage.SubAdminUser
	objects3, err := executor1.GetObjects()
	assert.NoError(t, err)
	assert.Len(t, objects3, 2)
	objects4, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, objects4, 1)
	assert.Equal(t, featureStartID, objects4[0].GetID())
}

func TestExecutorNormalDomain_PolicyList(t *testing.T) {
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	assert.Error(t, executor.SyncLatestVersionToAllDomain())

	assert.NoError(t, executor.BuildVersion())
	assert.NoError(t, executor.SyncLatestVersionToAllDomain())
	assert.NoError(t, reinitializeDomainWithWebFeature(stage, w.GetRoot()))

	executor1 := stage.Caskin.GetExecutor(provider)
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	roles, err := executor1.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles, 3)

	policy1, err := executor1.GetPolicyListByRole(roles[2])
	assert.NoError(t, err)
	assert.Len(t, policy1, 4)
	input1 := []*caskin.Policy{
		policy1[0], policy1[1], policy1[2], policy1[3],
		{roles[2], &example.Object{ID: featureStartID}, stage.Domain, caskin.Read},
	}
	// can't modify feature policy by Caskin.baseService
	provider.User = stage.AdminUser
	assert.NoError(t, executor1.ModifyPolicyListPerRole(roles[2], input1))
	policy2, err := executor1.GetPolicyListByRole(roles[2])
	assert.NoError(t, err)
	assert.Len(t, policy2, 4)
	provider.User = stage.SuperadminUser
	policy3, err := executor.NormalDomainGetPolicyListByRole(roles[2])
	assert.NoError(t, err)
	assert.Len(t, policy3, 0)

	// TODO issue 1: any way to fix the behaviour
	// can modify feature policy but can't get by Caskin.baseService when superadmin
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor1.ModifyPolicyListPerRole(roles[2], input1))
	policy4, err := executor1.GetPolicyListByRole(roles[2])
	assert.NoError(t, err)
	assert.Len(t, policy4, 4)
	policy5, err := executor.NormalDomainGetPolicyListByRole(roles[2])
	assert.NoError(t, err)
	assert.Len(t, policy5, 1)

	// can modify feature policy by WebFeature.baseService
	input2 := []*caskin.Policy{
		policy1[0], policy1[1], policy1[2], policy1[3],
		{roles[2], &example.Object{ID: featureStartID + 1}, stage.Domain, caskin.Read},
		{roles[2], &example.Object{ID: featureStartID + 1}, stage.Domain, caskin.Write},
	}
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrInValidAction, executor.NormalDomainModifyPolicyListPerRole(roles[2], input2))
	input2 = input2[0:5]
	assert.NoError(t, executor.NormalDomainModifyPolicyListPerRole(roles[2], input2))
	provider.User = stage.SubAdminUser
	policy6, err := executor1.GetPolicyListByRole(roles[2])
	assert.NoError(t, err)
	assert.Len(t, policy6, 4)
	policy7, err := executor.NormalDomainGetPolicyListByRole(roles[2])
	assert.NoError(t, err)
	assert.Len(t, policy7, 1)
}
