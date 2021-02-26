package caskin

import "math"

var (
	DefaultSuperadminRoleID   uint64 = math.MaxInt32
	DefaultSuperadminDomainID uint64 = math.MaxInt32

	DefaultSuperadminRoleName   = "superadmin_role"
	DefaultSuperadminDomainName = "superadmin_domain"
	// default
	DefaultSeparator = ","
)

type Option struct {
	// option of superadmin
	SuperAdminOption *SuperAdminOption `json:"super_admin_option"`

	// create new domain's function
	DomainCreator DomainCreator
}

type SuperAdminOption struct {
	// default is false
	Enable             bool `json:"enable"`
	// if there is superadmin domain and role record in metadata database.
	// default is false
	RealSuperadminInDB bool `json:"real_superadmin_in_db"`
	// provide superadmin Role
	Role               func() Role
	// provide superadmin Domain
	Domain             func() Domain
}

type DomainCreator func(Domain) ([]Role, []Object, []*Policy)

func (o *Option) IsEnableSuperAdmin() bool {
	return o.SuperAdminOption != nil && o.SuperAdminOption.Enable
}

func (o *Option) GetSuperAdminRole() Role {
	if !o.IsEnableSuperAdmin() {
		return nil
	}

	if o.SuperAdminOption.Role != nil {
		return o.SuperAdminOption.Role()
	}

	return &sampleSuperadminRole{}
}

func (o *Option) GetSuperAdminDomain() Domain {
	if !o.IsEnableSuperAdmin() {
		return nil
	}

	if o.SuperAdminOption.Domain != nil {
		return o.SuperAdminOption.Domain()
	}

	return &sampleSuperAdminDomain{}
}
