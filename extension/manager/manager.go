package manager

import (
	"fmt"
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/domain_creator"
	"github.com/awatercolorpen/caskin/extension/web_feature"
)

var (
	ErrNoInitialization               = fmt.Errorf("it is no initialization")
	ErrInitializationNilDB            = fmt.Errorf("gorm DB is nil")
	ErrExtensionConfigurationConflict = fmt.Errorf("extension configuration conflict")
)

type Manager struct {
	domainCreatorFactory *domain_creator.Factory
	webFeature           *web_feature.WebFeature
	caskin               *caskin.Caskin
}

func (m *Manager) GetDomainCreatorFactory() (*domain_creator.Factory, error) {
	if m.domainCreatorFactory == nil {
		return nil, ErrNoInitialization
	}

	return m.domainCreatorFactory, nil
}

func (m *Manager) GetWebFeature() (*web_feature.WebFeature, error) {
	if m.webFeature == nil {
		return nil, ErrNoInitialization
	}

	return m.webFeature, nil
}

func (m *Manager) GetCaskin() (*caskin.Caskin, error) {
	if m.caskin == nil {
		return nil, ErrNoInitialization
	}

	return m.caskin, nil
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
	if configuration.DefaultNoPermissionObject != "" {
		caskin.DefaultNoPermissionObject = configuration.DefaultNoPermissionObject
	}

	m := &Manager{}
	// initialize prefix extension
	if extension := configuration.Extension; extension != nil {
		if extension.DomainCreator != nil {
			if domainCreatorFactory, err := m.extensionDomainCreator(configuration); err != nil {
				return nil, err
			} else {
				m.domainCreatorFactory = domainCreatorFactory
				configuration.DomainCreator = domainCreatorFactory.NewCreator
			}
		}
	}

	if configuration.DomainCreator == nil {
		return nil, caskin.ErrInitializationNilDomainCreator
	}
	if configuration.Enforcer == nil {
		return nil, caskin.ErrInitializationNilEnforcer
	}
	if configuration.EntryFactory == nil {
		return nil, caskin.ErrInitializationNilEntryFactory
	}
	if configuration.MetaDB == nil {
		return nil, caskin.ErrInitializationNilMetaDB
	}

	// initialize caskin
	ckOptions := &caskin.Options{
		SuperadminRole:   configuration.SuperadminRole,
		SuperadminDomain: configuration.SuperadminDomain,
		DomainCreator:    configuration.DomainCreator,
		Enforcer:         configuration.Enforcer,
		EntryFactory:     configuration.EntryFactory,
		MetaDB:           configuration.MetaDB,
	}

	if ck, err := caskin.New(ckOptions); err != nil {
		return nil, err
	} else {
		m.caskin = ck
	}

	// initialize suffix extension
	if extension := configuration.Extension; extension != nil {
		if extension.WebFeature != nil {
			if webFeature, err := m.extensionWebFeature(configuration); err != nil {
				return nil, err
			} else {
				m.webFeature = webFeature
			}
		}
	}

	return m, nil
}
