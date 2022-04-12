package caskin_test

import (
	"testing"
)

func TestServer_NormalDomain_GetFeature(t *testing.T) {
	//stage, err := .NewPlaygroundWithSqlitePathAndWebFeature(t.TempDir())
	//assert.NoError(t, err)
	//w, err := web_feature.newWebFeature(stage)
	//assert.NoError(t, err)
	//provider := NewCachedProvider(nil, nil)
	//executor := w.GetExecutor(provider)
	//
	//provider.Domain = stage.Options.GetSuperadminDomain()
	//provider.User = stage.Superadmin
	//assert.Error(t, service.SyncLatestVersionToAllDomain())
	//
	//assert.NoError(t, service.BuildVersion())
	//assert.NoError(t, service.SyncLatestVersionToAllDomain())
	//assert.NoError(t, web_feature.reinitializeDomainWithWebFeature(stage, w.GetRoot()))
	//
	//executor1 := stage.Caskin.GetExecutor(provider)
	//provider.Domain = stage.Domain
	//provider.User = stage.AdminUser
	//roles, err := executor1.GetRoles()
	//assert.NoError(t, err)
	//assert.Len(t, roles, 3)
	//
	//policy1, err := executor1.GetPolicyListByRole(roles[2])
	//assert.NoError(t, err)
	//assert.Len(t, policy1, 4)
	//input1 := []*Policy{
	//	policy1[0], policy1[1], policy1[2], policy1[3],
	//	{roles[2], &example.Object{ID: web_feature.featureStartID}, stage.Domain, Read},
	//}
	//
	//assert.NoError(t, service.NormalDomainModifyPolicyListPerRole(roles[2], input1))
	//objects1, err := executor1.GetObjects()
	//assert.NoError(t, err)
	//assert.Len(t, objects1, 5)
	//objects2, err := service.NormalDomainGetFeatureObject()
	//assert.NoError(t, err)
	//assert.Len(t, objects2, 5)
	//
	//provider.User = stage.SubAdminUser
	//objects3, err := executor1.GetObjects()
	//assert.NoError(t, err)
	//assert.Len(t, objects3, 2)
	//objects4, err := service.NormalDomainGetFeatureObject()
	//assert.NoError(t, err)
	//assert.Len(t, objects4, 1)
	//assert.Equal(t, web_feature.featureStartID, objects4[0].GetID())
}

func TestServer_NormalDomain_PolicyList(t *testing.T) {
	//stage, err := web_feature.NewPlaygroundWithSqlitePathAndWebFeature(t.TempDir())
	//assert.NoError(t, err)
	//w, err := web_feature.newWebFeature(stage)
	//assert.NoError(t, err)
	//provider := NewCachedProvider(nil, nil)
	//executor := w.GetExecutor(provider)
	//
	//provider.Domain = stage.Options.GetSuperadminDomain()
	//provider.User = stage.Superadmin
	//assert.Error(t, service.SyncLatestVersionToAllDomain())
	//
	//assert.NoError(t, service.BuildVersion())
	//assert.NoError(t, service.SyncLatestVersionToAllDomain())
	//assert.NoError(t, web_feature.reinitializeDomainWithWebFeature(stage, w.GetRoot()))
	//
	//executor1 := stage.Caskin.GetExecutor(provider)
	//provider.Domain = stage.Domain
	//provider.User = stage.AdminUser
	//roles, err := executor1.GetRoles()
	//assert.NoError(t, err)
	//assert.Len(t, roles, 3)
	//
	//policy1, err := executor1.GetPolicyListByRole(roles[2])
	//assert.NoError(t, err)
	//assert.Len(t, policy1, 4)
	//input1 := []*Policy{
	//	policy1[0], policy1[1], policy1[2], policy1[3],
	//	{roles[2], &example.Object{ID: web_feature.featureStartID}, stage.Domain, Read},
	//}
	//// can't modify feature policy by Caskin.baseService
	//provider.User = stage.AdminUser
	//assert.NoError(t, executor1.ModifyPolicyListPerRole(roles[2], input1))
	//policy2, err := executor1.GetPolicyListByRole(roles[2])
	//assert.NoError(t, err)
	//assert.Len(t, policy2, 4)
	//provider.User = stage.Superadmin
	//policy3, err := service.NormalDomainGetPolicyListByRole(roles[2])
	//assert.NoError(t, err)
	//assert.Len(t, policy3, 0)
	//
	//// TODO issue 1: any way to fix the behaviour
	//// can modify feature policy but can't get by Caskin.baseService when superadmin
	//provider.User = stage.Superadmin
	//assert.NoError(t, executor1.ModifyPolicyListPerRole(roles[2], input1))
	//policy4, err := executor1.GetPolicyListByRole(roles[2])
	//assert.NoError(t, err)
	//assert.Len(t, policy4, 4)
	//policy5, err := service.NormalDomainGetPolicyListByRole(roles[2])
	//assert.NoError(t, err)
	//assert.Len(t, policy5, 1)
	//
	//// can modify feature policy by WebFeature.baseService
	//input2 := []*Policy{
	//	policy1[0], policy1[1], policy1[2], policy1[3],
	//	{roles[2], &example.Object{ID: web_feature.featureStartID + 1}, stage.Domain, Read},
	//	{roles[2], &example.Object{ID: web_feature.featureStartID + 1}, stage.Domain, Write},
	//}
	//provider.User = stage.AdminUser
	//assert.Equal(t, ErrInValidAction, service.NormalDomainModifyPolicyListPerRole(roles[2], input2))
	//input2 = input2[0:5]
	//assert.NoError(t, service.NormalDomainModifyPolicyListPerRole(roles[2], input2))
	//provider.User = stage.SubAdminUser
	//policy6, err := executor1.GetPolicyListByRole(roles[2])
	//assert.NoError(t, err)
	//assert.Len(t, policy6, 4)
	//policy7, err := service.NormalDomainGetPolicyListByRole(roles[2])
	//assert.NoError(t, err)
	//assert.Len(t, policy7, 1)
}
