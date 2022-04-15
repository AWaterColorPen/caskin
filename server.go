package caskin

type server struct {
	Enforcer   IEnforcer
	DB         MetaDB
	Dictionary IDictionary
}

type currentServer struct {
	server
	CurrentUser   User
	CurrentDomain Domain
}

func New(options *Options, opts ...Option) (IService, error) {
	options = options.newOptions(opts...)
	// set default caskin option
	if options.DefaultSuperadminDomainName != "" {
		DefaultSuperadminDomainName = options.DefaultSuperadminDomainName
	}
	if options.DefaultSuperadminRoleName != "" {
		DefaultSuperadminRoleName = options.DefaultSuperadminRoleName
	}

	if options.Enforcer == nil {
		return nil, ErrInitializationNilEnforcer
	}
	dictionary, err := NewDictionary(options.Dictionary)
	if err != nil {
		return nil, err
	}
	db, err := options.DB.NewDB()
	if err != nil {
		return nil, err
	}
	return &server{
		Enforcer:   NewEnforcer(options.Enforcer, DefaultFactory()),
		DB:         DefaultFactory().MetadataDB(db),
		Dictionary: dictionary,
	}, nil
}
