package manager

import (
	"fmt"

	"github.com/awatercolorpen/caskin"
)

var (
	ErrNoInitialization               = fmt.Errorf("it is no initialization")
	ErrInitializationNilDB            = fmt.Errorf("gorm DB is nil")
	ErrExtensionConfigurationConflict = fmt.Errorf("extension configuration conflict")
)

type Manager struct {
}

func NewManager(configuration *Configuration) (*Manager, error) {
	// set default caskin option
	if configuration.DefaultSeparator != "" {
		caskin.DefaultSeparator = configuration.DefaultSeparator
	}
	if configuration.DefaultSuperadminDomainID != 0 {
		caskin.DefaultSuperadminDomainID = configuration.DefaultSuperadminDomainID
	}
	if configuration.DefaultSuperadminDomainName != "" {
		caskin.DefaultSuperadminDomainName = configuration.DefaultSuperadminDomainName
	}
	if configuration.DefaultSuperadminRoleID != 0 {
		caskin.DefaultSuperadminRoleID = configuration.DefaultSuperadminRoleID
	}
	if configuration.DefaultSuperadminRoleName != "" {
		caskin.DefaultSuperadminRoleName = configuration.DefaultSuperadminRoleName
	}

	m := &Manager{}

	if configuration.Enforcer == nil {
		return nil, caskin.ErrInitializationNilEnforcer
	}
	if configuration.MetaDB == nil {
		return nil, caskin.ErrInitializationNilMetaDB
	}
	return m, nil
}
