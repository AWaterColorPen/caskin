package caskin_test

import (
	"encoding/json"
	"fmt"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutorDomain(t *testing.T) {
	stage, err := getStage(t)
	assert.NoError(t, err)
	provider := example.Provider{
		User:   stage.Domain,
		Domain: stage.SuperadminUser,
	}
	executor := stage.Caskin.GetExecutor(provider)

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

	bytes, _ = json.Marshal(domain)
	fmt.Println(string(bytes))
}