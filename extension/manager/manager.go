package manager

import (
	"fmt"
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/domain_creator"
	"github.com/awatercolorpen/caskin/extension/web_feature_old"
)

var (
	ErrNoInitialization               = fmt.Errorf("it is no initialization")
	ErrInitializationNilDB            = fmt.Errorf("gorm DB is nil")
	ErrExtensionConfigurationConflict = fmt.Errorf("extension configuration conflict")
)

type Manager struct {
	domainCreatorFactory *domain_creator.Factory
	webFeature           *web_feature_old.WebFeature
	caskin               *caskin.Caskin
}

func (m *Manager) GetDomainCreatorFactory() (*domain_creator.Factory, error) {
	if m.domainCreatorFactory == nil {
		return nil, ErrNoInitialization
	}

	return m.domainCreatorFactory, nil
}

func (m *Manager) GetWebFeature() (*web_feature_old.WebFeature, error) {
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

	// set default caskin web_feature_old option
	if configuration.DefaultBackendRootPath != "" {
		web_feature_old.DefaultBackendRootPath = configuration.DefaultBackendRootPath
	}
	if configuration.DefaultBackendRootMethod != "" {
		web_feature_old.DefaultBackendRootMethod = configuration.DefaultBackendRootMethod
	}
	if configuration.DefaultBackendRootDescription != "" {
		web_feature_old.DefaultBackendRootDescription = configuration.DefaultBackendRootDescription
	}
	if configuration.DefaultBackendRootGroup != "" {
		web_feature_old.DefaultBackendRootGroup = configuration.DefaultBackendRootGroup
	}
	if configuration.DefaultFeatureRootName != "" {
		web_feature_old.DefaultFeatureRootName = configuration.DefaultFeatureRootName
	}
	if configuration.DefaultFeatureRootDescription != "" {
		web_feature_old.DefaultFeatureRootDescription = configuration.DefaultFeatureRootDescription
	}
	if configuration.DefaultFeatureRootGroup != "" {
		web_feature_old.DefaultFeatureRootGroup = configuration.DefaultFeatureRootGroup
	}
	if configuration.DefaultFrontendRootKey != "" {
		web_feature_old.DefaultFrontendRootKey = configuration.DefaultFrontendRootKey
	}
	if configuration.DefaultFrontendRootType != "" {
		web_feature_old.DefaultFrontendRootType = web_feature_old.FrontendType(configuration.DefaultFrontendRootType)
	}
	if configuration.DefaultFrontendRootDescription != "" {
		web_feature_old.DefaultFrontendRootDescription = configuration.DefaultFrontendRootDescription
	}
	if configuration.DefaultFrontendRootGroup != "" {
		web_feature_old.DefaultFrontendRootGroup = configuration.DefaultFrontendRootGroup
	}
	if configuration.DefaultSuperRootName != "" {
		web_feature_old.DefaultSuperRootName = configuration.DefaultSuperRootName
	}
	if configuration.DefaultWebFeatureVersionTableName != "" {
		web_feature_old.DefaultWebFeatureVersionTableName = configuration.DefaultWebFeatureVersionTableName
	}

	// set default caskin domain creator option
	if configuration.DefaultDomainCreatorObjectTableName != "" {
		domain_creator.DefaultDomainCreatorObjectTableName = configuration.DefaultDomainCreatorObjectTableName
	}
	if configuration.DefaultDomainCreatorRoleTableName != "" {
		domain_creator.DefaultDomainCreatorRoleTableName = configuration.DefaultDomainCreatorRoleTableName
	}
	if configuration.DefaultDomainCreatorPolicyTableName != "" {
		domain_creator.DefaultDomainCreatorPolicyTableName = configuration.DefaultDomainCreatorPolicyTableName
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
		SuperadminDisable: configuration.SuperadminDisable,
		SuperadminRole:    configuration.SuperadminRole,
		SuperadminDomain:  configuration.SuperadminDomain,
		DomainCreator:     configuration.DomainCreator,
		Enforcer:          configuration.Enforcer,
		EntryFactory:      configuration.EntryFactory,
		MetaDB:            configuration.MetaDB,
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

