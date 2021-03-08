package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
)

func TestExecutorUser_GetUserRolePair(t *testing.T) {
	stage, _ := getStage(t)
	provider := &example.Provider{
		User: stage.MemberUser,
		Domain: stage.Domain,
	}
	executor := stage.Caskin.GetExecutor(provider)

	list1, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list1, 0)

	provider.User = stage.AdminUser
	list2, err := executor.GetUserRolePair()
	assert.NoError(t, err)
	assert.Len(t, list2, 2)
}
