package caskin

import (
	"github.com/casbin/casbin/v2"
)

var (
	DefaultSuperadminRoleID   uint64 = 0
	DefaultSuperadminDomainID uint64 = 0

	DefaultSuperadminRoleName   = "superadmin_role"
	DefaultSuperadminDomainName = "superadmin_domain"
	DefaultNoPermissionObject   = "no_permission_object"

	// default
	DefaultSeparator = "$$"
)

type Option func(*Options)

// Options configuration for caskin
type Options struct {
	// options of superadmin
	SuperadminDisable bool          `json:"superadmin_disable"` // default is false
	SuperadminRole    RoleFactory   `json:"-"`                  // provide superadmin Role
	SuperadminDomain  DomainFactory `json:"-"`                  // provide superadmin Domain

	// options for implementations of the interface
	DomainCreator DomainCreator    `json:"-"`
	Enforcer      casbin.IEnforcer `json:"-"`
	EntryFactory  EntryFactory     `json:"-"`
	MetaDB        MetaDB           `json:"-"`
}

func (o *Options) newOptions(opts ...Option) *Options {
	for _, v := range opts {
		v(o)
	}
	return o
}

func (o *Options) IsDisableSuperadmin() bool {
	return o.SuperadminDisable
}

func (o *Options) GetSuperadminRole() Role {
	if o.IsDisableSuperadmin() {
		return nil
	}

	if o.SuperadminRole != nil {
		return o.SuperadminRole()
	}

	return &SampleSuperadminRole{
		ID:   DefaultSuperadminRoleID,
		Name: DefaultSuperadminRoleName,
	}
}

func (o *Options) GetSuperadminDomain() Domain {
	if o.IsDisableSuperadmin() {
		return nil
	}

	if o.SuperadminDomain != nil {
		return o.SuperadminDomain()
	}

	return &SampleSuperadminDomain{
		ID:   DefaultSuperadminDomainID,
		Name: DefaultSuperadminDomainName,
	}
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
