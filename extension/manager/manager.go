package manager

import (
    "fmt"
    "sync"

    "github.com/awatercolorpen/caskin"
    "github.com/awatercolorpen/caskin/extension/domain_creator"
    "github.com/awatercolorpen/caskin/extension/web_feature"
)

var (
	instance *manager
	once     sync.Once
)

type manager struct {
}

func newManager(configuration *Configuration) (*manager, error) {
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

    // set default caskin web_feature option
    if configuration.DefaultBackendRootPath != "" {
        web_feature.DefaultBackendRootPath = configuration.DefaultBackendRootPath
    }
    if configuration.DefaultBackendRootMethod != "" {
        web_feature.DefaultBackendRootMethod = configuration.DefaultBackendRootMethod
    }
    if configuration.DefaultBackendRootDescription != "" {
        web_feature.DefaultBackendRootDescription = configuration.DefaultBackendRootDescription
    }
    if configuration.DefaultBackendRootGroup != "" {
        web_feature.DefaultBackendRootGroup = configuration.DefaultBackendRootGroup
    }
    if configuration.DefaultFeatureRootName != "" {
        web_feature.DefaultFeatureRootName = configuration.DefaultFeatureRootName
    }
    if configuration.DefaultFeatureRootDescription != "" {
        web_feature.DefaultFeatureRootDescription = configuration.DefaultFeatureRootDescription
    }
    if configuration.DefaultFeatureRootGroup != "" {
        web_feature.DefaultFeatureRootGroup = configuration.DefaultFeatureRootGroup
    }
    if configuration.DefaultFrontendRootKey != "" {
        web_feature.DefaultFrontendRootKey = configuration.DefaultFrontendRootKey
    }
    if configuration.DefaultFrontendRootType != "" {
        web_feature.DefaultFrontendRootType = web_feature.FrontendType(configuration.DefaultFrontendRootType)
    }
    if configuration.DefaultFrontendRootDescription != "" {
        web_feature.DefaultFrontendRootDescription = configuration.DefaultFrontendRootDescription
    }
    if configuration.DefaultFrontendRootGroup != "" {
        web_feature.DefaultFrontendRootGroup = configuration.DefaultFrontendRootGroup
    }
    if configuration.DefaultSuperRootName != "" {
        web_feature.DefaultSuperRootName = configuration.DefaultSuperRootName
    }
    if configuration.DefaultWebFeatureVersionTableName != "" {
        web_feature.DefaultWebFeatureVersionTableName = configuration.DefaultWebFeatureVersionTableName
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

    return nil, nil
}

func Init(configuration *Configuration) (err error) {
    if instance != nil {
        return fmt.Errorf("can't duplicate initialization")
    }

    once.Do(func() {
        instance, err = newManager(configuration)
    })
    return nil
}
