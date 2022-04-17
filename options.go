package caskin

import "github.com/casbin/casbin/v2"

var (
	DefaultSuperadminRoleName   = "superadmin_role"
	DefaultSuperadminDomainName = "superadmin_domain"
)

type Option func(*Options)

// Options configuration for caskin
type Options struct {
	// default caskin option
	DefaultSuperadminDomainName string            `json:"default_superadmin_domain_name"`
	DefaultSuperadminRoleName   string            `json:"default_superadmin_role_name"`
	Dictionary                  *DictionaryOption `json:"dictionary"`
	DB                          *DBOption         `json:"db"`
	// options for implementations of the interface
	Enforcer casbin.IEnforcer `json:"-"`
}

func (o *Options) newOptions(opts ...Option) *Options {
	for _, v := range opts {
		v(o)
	}
	return o
}
