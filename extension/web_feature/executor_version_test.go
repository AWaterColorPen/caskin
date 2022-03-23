package web_feature_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/extension/web_feature"
	"github.com/stretchr/testify/assert"
)

func TestExecutorVersion_BuildVersion(t *testing.T) {
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	assert.Equal(t, caskin.ErrProviderGet, executor.BuildVersion())
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.BuildVersion())
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrIsNotSuperAdmin, executor.BuildVersion())
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.BuildVersion())

	list1, err := executor.GetVersion()
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	assert.Equal(t, "1dd781e212bc2ff56053e2f09d7b399e0ca9f120b1784eb8501aab4d44c8cbee", list1[0].SHA256)

	// update a backend property can't change build version
	backend1 := &web_feature.Backend{Path: "api/backend", Method: "GET"}
	object1 := &example.Object{ID: backendStartID, ObjectID: 1}
	assert.NoError(t, executor.UpdateBackend(backend1, object1))

	assert.Error(t, executor.BuildVersion())
}

func TestExecutorVersion_SyncVersionToAllDomain(t *testing.T) {
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(nil, nil)
	executor := w.GetExecutor(provider)

	assert.Equal(t, caskin.ErrProviderGet, executor.SyncVersionToAllDomain(nil))
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	assert.Equal(t, caskin.ErrCanOnlyAllowAtValidDomain, executor.SyncVersionToAllDomain(nil))
	provider.Domain = stage.Options.GetSuperadminDomain()
	assert.Equal(t, caskin.ErrIsNotSuperAdmin, executor.SyncVersionToAllDomain(nil))
	provider.User = stage.SuperadminUser
	assert.Error(t, executor.SyncVersionToAllDomain(nil))

	version := &web_feature.WebFeatureVersion{
		SHA256: "sha256",
	}
	assert.Error(t, executor.SyncVersionToAllDomain(version))

	assert.NoError(t, executor.BuildVersion())
	versions, err := executor.GetVersion()
	assert.NoError(t, err)
	assert.Len(t, versions, 1)
	assert.NoError(t, executor.SyncVersionToAllDomain(versions[0]))
	// test twice
	assert.NoError(t, executor.SyncVersionToAllDomain(versions[0]))
}

func TestExecutorVersion_SyncLatestVersionToAllDomain(t *testing.T) {
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
	// test twice
	assert.NoError(t, executor.SyncLatestVersionToAllDomain())

	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	list1, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	assert.NoError(t, reinitializeDomainWithWebFeature(stage, w.GetRoot()))
	list2, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, list2, 5)
	list3, err := executor.NormalDomainGetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, list3, 1)

	provider.User = stage.MemberUser
	list4, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, list4, 0)
	list5, err := executor.NormalDomainGetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, list5, 0)

	// TODO change feature relation, and test re sync will delete some relation
}

func TestExecutorVersion_SyncCompatible(t *testing.T) {
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(stage.SuperadminUser, stage.Options.GetSuperadminDomain())
	executor := w.GetExecutor(provider)
	assert.NoError(t, executor.BuildVersion())
	assert.NoError(t, executor.SyncLatestVersionToAllDomain())
	assert.NoError(t, reinitializeDomainWithWebFeature(stage, w.GetRoot()))

	// delete feature object
	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.DeleteFeature(&example.Object{ID: featureStartID}))
	assert.NoError(t, executor.DeleteBackend(&example.Object{ID: backendStartID + 1}))

	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	assert.Equal(t, caskin.ErrInCompatible, executor.SyncLatestVersionToAllDomain())

	assert.NoError(t, executor.BuildVersion())
	versions, err := executor.GetVersion()
	assert.NoError(t, err)
	assert.Len(t, versions, 2)
	assert.Equal(t, caskin.ErrInCompatible, executor.SyncVersionToAllDomain(versions[0]))
	assert.NoError(t, executor.SyncLatestVersionToAllDomain())
}

func TestExecutorVersion_SyncToUpdateAuth(t *testing.T) {
	stage, err := newStageWithSqlitePathAndWebFeature(t.TempDir())
	assert.NoError(t, err)
	w, err := newWebFeature(stage)
	assert.NoError(t, err)
	provider := caskin.NewCachedProvider(stage.SuperadminUser, stage.Options.GetSuperadminDomain())
	executor := w.GetExecutor(provider)
	assert.NoError(t, executor.BuildVersion())
	assert.NoError(t, executor.SyncLatestVersionToAllDomain())
	assert.NoError(t, reinitializeDomainWithWebFeature(stage, w.GetRoot()))

	// add one feature policy for sub admin
	executor1 := stage.Caskin.GetExecutor(provider)
	provider.Domain = stage.Domain
	provider.User = stage.AdminUser
	roles, err := executor1.GetRoles()
	assert.NoError(t, err)
	assert.Len(t, roles, 3)
	input := []*caskin.Policy{{roles[2], &example.Object{ID: featureStartID}, stage.Domain, caskin.Read}}
	assert.NoError(t, executor.NormalDomainModifyPolicyListPerRole(roles[2], input))

	// delete feature object
	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.DeleteFeature(&example.Object{ID: featureStartID}))
	assert.NoError(t, executor.DeleteBackend(&example.Object{ID: backendStartID + 1}))

	// before sync, it should have permission if the object of backend and frontend is not deleted
	// before sync, it can not manage feature policy any more
	provider.Domain = stage.Domain
	provider.User = stage.SubAdminUser
	assert.NoError(t, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend", Method: "GET"}))
	assert.Equal(t, caskin.ErrNoBackendAPIPermission, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend", Method: "POST"}))
	c1, err := executor.AuthFrontendCaskinStruct("abc")
	assert.NoError(t, err)
	assert.Len(t, c1.P, 1)
	objects1, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, objects1, 0)
	policy1, err := executor.NormalDomainGetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, policy1, 0)

	provider.User = stage.AdminUser
	objects2, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, objects2, 4)
	policy2, err := executor.NormalDomainGetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, policy2, 1)

	// if recover the deleted feature before sync, it can manage feature policy again
	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.RecoverFeature(&web_feature.Feature{Name: "backend"}, &example.Object{}))
	provider.Domain = stage.Domain
	provider.User = stage.SubAdminUser
	objects3, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, objects3, 1)
	policy3, err := executor.NormalDomainGetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, policy3, 1)

	// sync new version
	provider.Domain = stage.Options.GetSuperadminDomain()
	provider.User = stage.SuperadminUser
	assert.NoError(t, executor.DeleteFeature(&example.Object{ID: featureStartID}))
	assert.NoError(t, executor.BuildVersion())
	assert.NoError(t, executor.SyncLatestVersionToAllDomain())

	// after sync, it should have no permission
	provider.Domain = stage.Domain
	provider.User = stage.SubAdminUser
	assert.Equal(t, caskin.ErrNoBackendAPIPermission, executor.AuthBackendAPIEnforce(&web_feature.Backend{Path: "api/backend", Method: "GET"}))
	c2, err := executor.AuthFrontendCaskinStruct("abc")
	assert.NoError(t, err)
	assert.Len(t, c2.P, 0)
	objects4, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, objects4, 0)
	policy4, err := executor.NormalDomainGetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, policy4, 0)

	provider.User = stage.AdminUser
	objects5, err := executor.NormalDomainGetFeatureObject()
	assert.NoError(t, err)
	assert.Len(t, objects5, 4)
	policy5, err := executor.NormalDomainGetPolicyList()
	assert.NoError(t, err)
	assert.Len(t, policy5, 1)
}
