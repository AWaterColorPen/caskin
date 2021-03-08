package caskin_test

import (
	"encoding/json"
	"fmt"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutorDomain(t *testing.T) {
	env, err := getInitializedTestCaskin(t)
	assert.NoError(t, err)
	executor := env.caskinClient.GetExecutor(env.provider)

	domain := &example.Domain{
		Name:      "test_domain_02",
	}
	assert.NoError(t, executor.CreateDomain(domain))

	assert.NoError(t, executor.DeleteDomain(domain))

	assert.NoError(t, executor.RecoverDomain(domain))

	//allDomain, err := executor.GetAllDomain()
	assert.NoError(t, err)
	bytes, _ := json.Marshal(domain)
	fmt.Println(string(bytes))

	domain.Name = "test_domain_022"
	assert.NoError(t, executor.UpdateDomain(domain))
}