package manager

import (
	"github.com/awatercolorpen/caskin"
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type Configuration struct {
	// default caskin option
	DefaultSeparator            string `json:"default_separator"              yaml:"default_separator"`
	DefaultSuperadminDomainID   uint64 `json:"default_superadmin_domain_id"   yaml:"default_superadmin_domain_id"`
	DefaultSuperadminDomainName string `json:"default_superadmin_domain_name" yaml:"default_superadmin_domain_name"`
	DefaultSuperadminRoleID     uint64 `json:"default_superadmin_role_id"     yaml:"default_superadmin_role_id"`
	DefaultSuperadminRoleName   string `json:"default_superadmin_role_name"   yaml:"default_superadmin_role_name"`
	DefaultNoPermissionObject   string `json:"default_no_permission_object"   yaml:"default_no_permission_object"`

	// dependencies
	DB *gorm.DB

	// implementations of the caskin interface
	DomainCreator caskin.DomainCreator
	Enforcer      casbin.IEnforcer
	MetaDB        caskin.MetaDB
}
