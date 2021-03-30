package manager

import (
	"github.com/awatercolorpen/caskin"
	"github.com/casbin/casbin/v2"
)

type Configuration struct {
	// default caskin option
	DefaultSeparator            string `json:"default_separator"              yaml:"default_separator"`
	DefaultSuperadminDomainID   uint64 `json:"default_superadmin_domain_id"   yaml:"default_superadmin_domain_id"`
	DefaultSuperadminDomainName string `json:"default_superadmin_domain_name" yaml:"default_superadmin_domain_name"`
	DefaultSuperadminRoleID     uint64 `json:"default_superadmin_role_id"     yaml:"default_superadmin_role_id"`
	DefaultSuperadminRoleName   string `json:"default_superadmin_role_name"   yaml:"default_superadmin_role_name"`
	DefaultNoPermissionObject   string `json:"default_no_permission_object"   yaml:"default_no_permission_object"`

	// default caskin web_feature option
	DefaultBackendRootPath            string `json:"default_backend_root_path"              yaml:"default_backend_root_path"`
	DefaultBackendRootMethod          string `json:"default_backend_root_method"            yaml:"default_backend_root_method"`
	DefaultBackendRootDescription     string `json:"default_backend_root_description"       yaml:"default_backend_root_description"`
	DefaultBackendRootGroup           string `json:"default_backend_root_group"             yaml:"default_backend_root_group"`
	DefaultFeatureRootName            string `json:"default_feature_root_name"              yaml:"default_feature_root_name"`
	DefaultFeatureRootDescription     string `json:"default_feature_root_description"       yaml:"default_feature_root_description"`
	DefaultFeatureRootGroup           string `json:"default_feature_root_group"             yaml:"default_feature_root_group"`
	DefaultFrontendRootKey            string `json:"default_frontend_root_key"              yaml:"default_frontend_root_key"`
	DefaultFrontendRootType           string `json:"default_frontend_root_type"             yaml:"default_frontend_root_type"`
	DefaultFrontendRootDescription    string `json:"default_frontend_root_description"      yaml:"default_frontend_root_description"`
	DefaultFrontendRootGroup          string `json:"default_frontend_root_group"            yaml:"default_frontend_root_group"`
	DefaultSuperRootName              string `json:"default_super_root_name"                yaml:"default_super_root_name"`
	DefaultWebFeatureVersionTableName string `json:"default_web_feature_version_table_name" yaml:"default_web_feature_version_table_name"`

	// default caskin domain creator option
	DefaultDomainCreatorObjectTableName string `json:"default_domain_creator_object_table_name" yaml:"default_domain_creator_object_table_name"`
	DefaultDomainCreatorRoleTableName   string `json:"default_domain_creator_role_table_name"   yaml:"default_domain_creator_role_table_name"`
	DefaultDomainCreatorPolicyTableName string `json:"default_domain_creator_policy_table_name" yaml:"default_domain_creator_policy_table_name"`

	// configurations for superadmin
	SuperadminDisable bool `json:"superadmin_disable" yaml:"superadmin_disable"`

	// configurations for web_feature
	WebFeatureCacheDisable bool `json:"web_feature_cache_disable" yaml:"web_feature_cache_disable"`

	// implementations of the caskin interface
	DomainCreator caskin.DomainCreator
	Enforcer      casbin.IEnforcer
	EntryFactory  caskin.EntryFactory
	MetaDB        caskin.MetaDB

	// implementations of the caskin superadmin interface
	SuperadminDomain caskin.DomainFactory
	SuperadminRole   caskin.RoleFactory
}
