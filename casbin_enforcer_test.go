package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
)

func TestEnforcer_EnforceRole(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	e := stage.Enforcer
	assert.NotNil(t, e)

	roles, err := stage.Service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(roles), 2)

	// Create a child role that inherits from admin role
	object1 := &example.Object{Name: "enforce_role_obj", Type: "role"}
	object1.ParentID = roles[0].GetObjectID()
	assert.NoError(t, stage.Service.CreateObject(stage.Admin, stage.Domain, object1))

	role1 := &example.Role{Name: "enforce_role_child", ObjectID: object1.ID}
	assert.NoError(t, stage.Service.CreateRole(stage.Admin, stage.Domain, role1))

	// Before adding role inheritance, role1 does NOT inherit from roles[0]
	ok, err := e.EnforceRole(role1, roles[0], stage.Domain)
	assert.NoError(t, err)
	assert.False(t, ok)

	// Add role inheritance: role1 inherits from roles[0]
	assert.NoError(t, stage.Service.AddRoleG(stage.Admin, stage.Domain, role1, roles[0]))

	// After adding inheritance, role1 should inherit roles[0]
	ok2, err2 := e.EnforceRole(role1, roles[0], stage.Domain)
	assert.NoError(t, err2)
	assert.True(t, ok2)

	// role1 should NOT inherit roles[1] (unrelated)
	ok3, err3 := e.EnforceRole(role1, roles[1], stage.Domain)
	assert.NoError(t, err3)
	assert.False(t, ok3)
}

func TestEnforcer_GetPoliciesForObjectInDomain(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	e := stage.Enforcer
	assert.NotNil(t, e)

	objects, err := stage.Service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(objects), 1)

	// Get policies for the first object in domain
	policies := e.GetPoliciesForObjectInDomain(objects[0], stage.Domain)
	// The admin role should have policies on the root role object
	assert.NotEmpty(t, policies)

	// Create a new object and verify it has no policies initially
	roles, err := stage.Service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	newObj := &example.Object{Name: "new_obj_for_policy_test", Type: "role"}
	newObj.ParentID = objects[0].GetID()
	assert.NoError(t, stage.Service.CreateObject(stage.Admin, stage.Domain, newObj))

	noPolicies := e.GetPoliciesForObjectInDomain(newObj, stage.Domain)
	assert.Empty(t, noPolicies)

	// Add a policy and verify
	policy := []*caskin.Policy{{roles[0], newObj, stage.Domain, caskin.Read}}
	assert.NoError(t, stage.Service.ModifyPolicyPerRole(stage.Admin, stage.Domain, roles[0], append(
		func() []*caskin.Policy {
			existing, _ := stage.Service.GetPolicyByRole(stage.Admin, stage.Domain, roles[0])
			return existing
		}(),
		policy[0],
	)))

	withPolicy := e.GetPoliciesForObjectInDomain(newObj, stage.Domain)
	assert.NotEmpty(t, withPolicy)
}

func TestEnforcer_GetRolesInDomain(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	e := stage.Enforcer
	assert.NotNil(t, e)

	// GetRolesInDomain returns roles that participate in role-inheritance (g policy).
	// The base playground only has user→role policies, so there are no role→role
	// relationships initially — the result is empty.
	rolesInitial := e.GetRolesInDomain(stage.Domain)
	assert.Empty(t, rolesInitial)

	// Create a role-inheritance relationship so GetRolesInDomain returns something.
	domainRoles, err := stage.Service.GetRole(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(domainRoles), 2)

	childObj := &example.Object{Name: "roles_in_domain_child_obj", Type: "role"}
	childObj.ParentID = domainRoles[0].GetObjectID()
	assert.NoError(t, stage.Service.CreateObject(stage.Admin, stage.Domain, childObj))

	childRole := &example.Role{Name: "roles_in_domain_child", ObjectID: childObj.ID}
	assert.NoError(t, stage.Service.CreateRole(stage.Admin, stage.Domain, childRole))
	assert.NoError(t, stage.Service.AddRoleG(stage.Admin, stage.Domain, childRole, domainRoles[0]))

	// Now GetRolesInDomain should return the roles involved in that inheritance.
	rolesAfter := e.GetRolesInDomain(stage.Domain)
	assert.NotEmpty(t, rolesAfter)
}

func TestEnforcer_GetObjectsInDomain(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	e := stage.Enforcer
	assert.NotNil(t, e)

	// The playground sets up objects in the domain
	objects := e.GetObjectsInDomain(stage.Domain)
	// Should return objects (at least root objects from ResetDomain)
	assert.NotEmpty(t, objects)
}

func TestEnforcer_GetPoliciesInDomain(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	e := stage.Enforcer
	assert.NotNil(t, e)

	// The playground sets up policies in the domain (admin and member roles have policies)
	policies := e.GetPoliciesInDomain(stage.Domain)
	assert.NotEmpty(t, policies)

	// Count should match what service returns
	servicePolicies, err := stage.Service.GetPolicy(stage.Admin, stage.Domain)
	assert.NoError(t, err)
	// service filters by permission, enforcer returns raw; enforcer count >= service count
	assert.GreaterOrEqual(t, len(policies), len(servicePolicies))
}
