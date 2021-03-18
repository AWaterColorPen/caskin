package caskin

type Caskin struct {
	options *Options
}

func (c *Caskin) GetExecutor(provider CurrentProvider) *Executor {
	e := NewEnforcer(c.options.Enforcer, c.options.EntryFactory)
	return &Executor{
		Enforcer: e,
		factory:  c.options.EntryFactory,
		DB:       c.options.MetaDB,
		provider: provider,
		options:  c.options,
	}
}

func (c *Caskin) GetOptions() *Options {
	return c.options
}

func New(options *Options, opts ...Option) (*Caskin, error) {
	options = options.newOptions(opts...)
	if options.DomainCreator == nil {
		return nil, ErrInitializationNilDomainCreator
	}
	if options.Enforcer == nil {
		return nil, ErrInitializationNilEnforcer
	}
	if options.EntryFactory == nil {
		return nil, ErrInitializationNilEntryFactory
	}
	if options.MetaDB == nil {
		return nil, ErrInitializationNilMetaDB
	}
	if err := options.MetaDB.AutoMigrate(
		options.EntryFactory.NewUser(),
		options.EntryFactory.NewDomain(),
		options.EntryFactory.NewRole(),
		options.EntryFactory.NewObject()); err != nil {
		return nil, err
	}
	return &Caskin{options: options}, nil
}
