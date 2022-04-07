package caskin

import (
	"github.com/casbin/casbin/v2"
)

var (
	DefaultSuperadminRoleID     uint64 = 0
	DefaultSuperadminDomainID   uint64 = 0
	DefaultSuperadminRoleName          = "superadmin_role"
	DefaultSuperadminDomainName        = "superadmin_domain"

	DefaultSeparator = "$$"
)

type Option func(*Options)

// Options configuration for caskin
type Options struct {
	// default caskin option
	DefaultSeparator            string `json:"default_separator"`
	DefaultSuperadminDomainID   uint64 `json:"default_superadmin_domain_id"`
	DefaultSuperadminDomainName string `json:"default_superadmin_domain_name"`
	DefaultSuperadminRoleID     uint64 `json:"default_superadmin_role_id"`
	DefaultSuperadminRoleName   string `json:"default_superadmin_role_name"`
	DefaultNoPermissionObject   string `json:"default_no_permission_object"`
	// options for implementations of the interface
	Enforcer casbin.IEnforcer `json:"-"`
	MetaDB   MetaDB           `json:"-"`
}

func (o *Options) newOptions(opts ...Option) *Options {
	for _, v := range opts {
		v(o)
	}
	return o
}

// EnforcerOption set the casbin.IEnforcer for the options
func EnforcerOption(enforcer casbin.IEnforcer) Option {
	return func(o *Options) {
		o.Enforcer = enforcer
	}
}

// MetaDBOption set the MetaDB for the options
func MetaDBOption(metaDB MetaDB) Option {
	return func(o *Options) {
		o.MetaDB = metaDB
	}
}
