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
	DefaultNoPermissionObject   = "no_permission_object"

	// default
	DefaultSeparator = "$$"
)

// SuperadminOption option of superadmin
type SuperadminOption struct {
	Disable bool          `json:"disable"` // default is false
	Role    func() Role   `json:"-"`       // provide superadmin Role
	Domain  func() Domain `json:"-"`       // provide superadmin Domain
}

type Option func(*Options)

// Options configuration for caskin
type Options struct {
	SuperadminOption *SuperadminOption `json:"super_admin_option"` // options for configuration

	// options for implementations of the interface
	DomainCreator DomainCreator
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

func (o *Options) IsDisableSuperAdmin() bool {
	return o.SuperadminOption != nil && o.SuperadminOption.Disable
}

func (o *Options) GetSuperadminRole() Role {
	if o.IsDisableSuperAdmin() {
		return nil
	}

	if o.SuperadminOption != nil && o.SuperadminOption.Role != nil {
		return o.SuperadminOption.Role()
	}

	return &SampleSuperadminRole{}
}

func (o *Options) GetSuperadminDomain() Domain {
	if o.IsDisableSuperAdmin() {
		return nil
	}

	if o.SuperadminOption != nil && o.SuperadminOption.Domain != nil {
		return o.SuperadminOption.Domain()
	}

	return &SampleSuperAdminDomain{}
}

// DomainCreatorOption set the DomainCreator for the options
func DomainCreatorOption(creator DomainCreator) Option {
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
