package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"testing"

	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestServer_GetFeature(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	feature, err := service.GetFeature()
	assert.NoError(t, err)
	assert.Len(t, feature, 3)
}

func TestServer_GetBackend(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	backend, err := service.GetBackend()
	assert.NoError(t, err)
	assert.Len(t, backend, 6)
}

func TestServer_GetFrontend(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	frontend, err := service.GetFrontend()
	assert.NoError(t, err)
	assert.Len(t, frontend, 3)
}

func TestServer_AuthBackend(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	backend1 := &caskin.Backend{Path: "api/backend", Method: "GET"}
	assert.NoError(t, service.AuthBackend(stage.Admin, stage.Domain, backend1))
	assert.Equal(t, caskin.ErrNoBackendPermission, service.AuthBackend(stage.Member, stage.Domain, backend1))
	backend2 := &caskin.Backend{Path: "no", Method: ""}
	assert.Equal(t, caskin.ErrNoBackendPermission, service.AuthBackend(stage.Admin, stage.Domain, backend2))
	assert.NoError(t, service.AuthBackend(stage.Superadmin, stage.Domain, backend2))
}

func TestServer_AuthFrontend(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	frontend1 := service.AuthFrontend(stage.Admin, stage.Domain)
	assert.Len(t, frontend1, 3)
	frontend2 := service.AuthFrontend(stage.Member, stage.Domain)
	assert.Len(t, frontend2, 0)
	frontend3 := service.AuthFrontend(stage.Superadmin, stage.Domain)
	assert.Len(t, frontend3, 3)
}

func TestServer_GetFeatureObject(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	list1, err := service.GetFeatureObject(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list1, 3)
	list2, err := service.GetFeatureObject(stage.Member, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list2, 0)
	list3, err := service.GetFeatureObject(stage.Superadmin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list3, 3)
}

func TestServer_GetFeaturePolicy(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	list1, err := service.GetFeaturePolicy(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	list2, err := service.GetFeaturePolicy(stage.Member, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list2, 0)
	list3, err := service.GetFeaturePolicy(stage.Superadmin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, list3, 1)
}

func TestServer_GetFeaturePolicyByRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	list1, err := service.GetFeaturePolicyByRole(stage.Admin, stage.Domain, roles[0])
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	list2, err := service.GetFeaturePolicyByRole(stage.Admin, stage.Domain, roles[1])
	assert.NoError(t, err)
	assert.Len(t, list2, 0)
}

func TestServer_ModifyFeaturePolicyPerRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	roles, err := service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.Len(t, roles, 2)

	feature, err := service.GetFeature()
	assert.NoError(t, err)
	assert.Len(t, feature, 3)

	policy1 := []*caskin.Policy{
		{Object: feature[1].ToObject()},
		{Object: feature[2].ToObject()},
	}
	assert.Equal(t, caskin.ErrNoWritePermission, service.ModifyFeaturePolicyPerRole(stage.Member, stage.Domain, roles[1], policy1))
	assert.NoError(t, service.ModifyFeaturePolicyPerRole(stage.Admin, stage.Domain, roles[1], policy1))

	list1, err := service.GetFeaturePolicyByRole(stage.Admin, stage.Domain, roles[1])
	assert.NoError(t, err)
	assert.Len(t, list1, 2)
}

func TestServer_ResetFeature(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service
	assert.NoError(t, service.ResetFeature(stage.Domain))
}
