package caskin_test

import (
    "testing"

    "github.com/awatercolorpen/caskin/example"
    "github.com/stretchr/testify/assert"

)

func TestExecutorUser(t *testing.T) {
    stage, _ := getStage(t)
    provider := example.Provider{
        User:   stage.SuperadminUser,
        Domain: stage.Domain,
    }
    executor := stage.Caskin.GetExecutor(provider)
    assert.Error(t, executor.DeleteUser(&example.User{ID: 4}))
    executor.
    domain := &example.Domain{Name: "domain_02"}
    assert.NoError(t, executor.CreateUser(domain))
    assert.NoError(t, executor.DeleteDomain(domain))
    domains1, err := executor.GetAllDomain()
    assert.NoError(t, err)
    assert.Len(t, domains1, 1)

    assert.NoError(t, executor.RecoverDomain(domain))
    domains, err := executor.GetAllDomain()
    assert.NoError(t, err)
    assert.Len(t, domains, 2)

    domain.Name = "domain_02_new_name"
    assert.NoError(t, executor.UpdateDomain(domain))

    domain0 := &example.Domain{ID:3, Name: "domain_03"}
    assert.Error(t, executor.UpdateDomain(domain0))
}
