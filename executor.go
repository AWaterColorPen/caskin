package caskin

type Executor struct {
	Enforcer IEnforcer
	DB       MetaDB
	provider CurrentProvider
	factory  Factory
	options  *Options
}

func (e *Executor) GetCurrentProvider() CurrentProvider {
	return e.provider
}
