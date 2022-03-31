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

func (e *Executor) filterWithNoError(user User, domain Domain, action Action, source any) []any {
	return Filter(e.Enforcer, user, domain, action, source)
}

func (e *Executor) checkObjectData(one ObjectData, action Action) error {
	return e.checkInternal(func(enforcer IEnforcer, user User, domain Domain) bool {
		return ObjectDataCheck(e.Enforcer, user, domain, one, action)
	}, action)
}

func (e *Executor) checkObject(one Object, action Action) error {
	return e.checkInternal(func(enforcer IEnforcer, user User, domain Domain) bool {
		return ObjectCheck(e.Enforcer, user, domain, one, action)
	}, action)
}

func (e *Executor) checkInternal(fn func(IEnforcer, User, Domain) bool, action Action) error {
	u, d, err := e.provider.Get()
	if err != nil {
		return err
	}

	if ok := fn(e.Enforcer, u, d); !ok {
		switch action {
		case Read:
			return ErrNoReadPermission
		case Write:
			return ErrNoWritePermission
		default:
		}
	}

	return nil
}
