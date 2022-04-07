package caskin

func New(options *Options, opts ...Option) (IService, error) {
	options = options.newOptions(opts...)
	// set default caskin option
	if options.DefaultSeparator != "" {
		DefaultSeparator = options.DefaultSeparator
	}
	if options.DefaultSuperadminDomainID != 0 {
		DefaultSuperadminDomainID = options.DefaultSuperadminDomainID
	}
	if options.DefaultSuperadminDomainName != "" {
		DefaultSuperadminDomainName = options.DefaultSuperadminDomainName
	}
	if options.DefaultSuperadminRoleID != 0 {
		DefaultSuperadminRoleID = options.DefaultSuperadminRoleID
	}
	if options.DefaultSuperadminRoleName != "" {
		DefaultSuperadminRoleName = options.DefaultSuperadminRoleName
	}

	if options.Enforcer == nil {
		return nil, ErrInitializationNilEnforcer
	}
	if options.MetaDB == nil {
		return nil, ErrInitializationNilMetaDB
	}
	return &server{}, nil
}
