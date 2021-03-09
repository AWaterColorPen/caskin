package caskin_test

import (
	"encoding/json"
	"fmt"
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutorPolicy(t *testing.T) {
	stage, _ := newStage(t)
	provider := &example.Provider{
		User:   stage.SuperadminUser,
		Domain: stage.Domain,
	}
	executor := stage.Caskin.GetExecutor(provider)

	policiesForRoleList, err := executor.GetAllPoliciesForRole()
	assert.NoError(t, err)
	assert.Len(t, policiesForRoleList, 2)
	//assert.Len(t, policiesForRoleList[0].Policies, 6)
	bytes, _ := json.Marshal(policiesForRoleList)
	fmt.Println(string(bytes))

	roles, err := executor.GetRoles()
	assert.NoError(t, err)
	bytes, _ = json.Marshal(roles)
	fmt.Println(string(bytes))

	objects, err := executor.GetObjects(caskin.ObjectTypeObject, caskin.ObjectTypeRole, caskin.ObjectTypeDefault)
	assert.NoError(t, err)

	domain := stage.Domain
	policiesForRole := &caskin.PoliciesForRole{
		Role: roles[0],
		Policies: []*caskin.Policy{
			{roles[0], objects[0], domain, caskin.Read},
			{roles[0], objects[1], domain, caskin.Read},
			{roles[0], objects[1], domain, caskin.Write},
			{roles[0], objects[2], domain, caskin.Read},
			{roles[0], objects[2], domain, caskin.Write},
		},
	}

	assert.NoError(t, executor.ModifyPoliciesForRole(policiesForRole))

	policiesForRoleList, err = executor.GetAllPoliciesForRole()
	assert.NoError(t, err)
	assert.Len(t, policiesForRoleList, 2)
	//assert.Len(t, policiesForRoleList[0].Policies, 5)
	bytes, _ = json.Marshal(policiesForRoleList)
	fmt.Println(string(bytes))
}
