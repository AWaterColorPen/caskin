package caskin

type Caskin struct {
	options *Options
}

//
// func (c *Caskin) GetExecutor(provider CurrentProvider) *server {
// 	e := NewEnforcer(c.options.Enforcer, c.options.Factory)
// 	return &server{
// 		Enforcer: e,
// 		DB:       c.options.MetaDB,
// 		provider: provider,
// 		options:  c.options,
// 	}
// }

// func (c *Caskin) GetOptions() *Options {
// 	return c.options
// }

func New(options *Options, opts ...Option) (*Caskin, error) {
	options = options.newOptions(opts...)
	if options.Enforcer == nil {
		return nil, ErrInitializationNilEnforcer
	}
	if options.MetaDB == nil {
		return nil, ErrInitializationNilMetaDB
	}
	if err := options.MetaDB.AutoMigrate(); err != nil {
		return nil, err
	}
	return &Caskin{options: options}, nil
}
