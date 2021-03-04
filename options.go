package caskin

import (
	"math"

	"github.com/casbin/casbin/v2"
)

var (
	DefaultSuperadminRoleID   uint64 = math.MaxInt32
	DefaultSuperadminDomainID uint64 = math.MaxInt32

	DefaultSuperadminRoleName   = "superadmin_role"
	DefaultSuperadminDomainName = "superadmin_domain"

	// default
	DefaultSeparator = "$$"
)

// SuperAdminOption option of superadmin
type SuperAdminOption struct {
	// default is false
	Enable bool `json:"enable"`
	// if there is superadmin domain and role record in metadata database.
	// default is false
	RealSuperadminInDB bool `json:"real_superadmin_in_db"`
	// provide superadmin Role
	Role func() Role
	// provide superadmin Domain
	Domain func() Domain
}

type Option func(*Options)

// Options configuration for caskin
type Options struct {
	// options for configuration
	SuperAdminOption *SuperAdminOption `json:"super_admin_option"`

	// options for implementations of the interface
	DomainCreator Creator
	Enforcer      casbin.IEnforcer
	EntryFactory  EntryFactory
	MetaDB        MetaDB
}

func (o *Options) newOptions(opts ...Option) *Options {
	for _, v := range opts {
		v(o)
	}
	return o
}

func (o *Options) IsEnableSuperAdmin() bool {
	return o.SuperAdminOption != nil && o.SuperAdminOption.Enable
}

func (o *Options) GetSuperAdminRole() Role {
	if !o.IsEnableSuperAdmin() {
		return nil
	}

	if o.SuperAdminOption.Role != nil {
		return o.SuperAdminOption.Role()
	}

	return &sampleSuperadminRole{}
}

func (o *Options) GetSuperAdminDomain() Domain {
	if !o.IsEnableSuperAdmin() {
		return nil
	}

	if o.SuperAdminOption.Domain != nil {
		return o.SuperAdminOption.Domain()
	}

	return &sampleSuperAdminDomain{}
}

// DomainCreatorOption set the DomainCreator for the options
func DomainCreatorOption(creator Creator) Option {
	return func(o *Options) {
		o.DomainCreator = creator
	}
}

// EnforcerOption set the casbin.IEnforcer for the options
func EnforcerOption(enforcer casbin.IEnforcer) Option {
	return func(o *Options) {
		o.Enforcer = enforcer
	}
}

// EntryFactoryOption set the EntryFactory for the options
func EntryFactoryOption(entryFactory EntryFactory) Option {
	return func(o *Options) {
		o.EntryFactory = entryFactory
	}
}

// MetaDBOption set the MetaDB for the options
func MetaDBOption(metaDB MetaDB) Option {
	return func(o *Options) {
		o.MetaDB = metaDB
	}
}
