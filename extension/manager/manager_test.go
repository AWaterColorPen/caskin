package manager_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/extension/manager"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	config := &manager.Configuration{}
	assert.Equal(t, caskin.ErrInitializationNilDomainCreator, manager.Init(config))

	config.DomainCreator = example.NewDomainCreator
	assert.Equal(t, caskin.ErrInitializationNilEnforcer, manager.Init(config))
}
