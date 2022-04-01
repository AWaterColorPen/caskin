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

func GetSuperadminRole() Role {
	return &SampleSuperadminRole{
		ID:   DefaultSuperadminRoleID,
		Name: DefaultSuperadminRoleName,
	}
}

func GetSuperadminDomain() Domain {
	return &SampleSuperadminDomain{
		ID:   DefaultSuperadminDomainID,
		Name: DefaultSuperadminDomainName,
	}
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
