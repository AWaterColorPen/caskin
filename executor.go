package caskin

type Executor struct {
	Enforcer IEnforcer
	DB       MetaDB
	provider CurrentProvider
	factory  EntryFactory
	options  *Options
}

func (e *Executor) GetCurrentProvider() CurrentProvider {
	return e.provider
}

func (e *Executor) newObject() TreeNodeEntry {
	return e.factory.NewObject()
}

func (e *Executor) newRole() TreeNodeEntry {
	return e.factory.NewRole()
}

func (e *Executor) filter(action Action, source interface{}) ([]interface{}, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}
	return Filter(e.Enforcer, u, d, action, source), nil
}

func (e *Executor) filterWithNoError(user User, domain Domain, action Action, source interface{}) []interface{} {
	return Filter(e.Enforcer, user, domain, action, source)
}

func (e *Executor) checkObjectData(one ObjectData, action Action) error {
	return e.checkInternal(func(enforcer IEnforcer, user User, domain Domain) bool {
		return CheckObjectData(e.Enforcer, user, domain, one, action)
	}, action)
}

func (e *Executor) checkObject(one Object, action Action) error {
	return e.checkInternal(func(enforcer IEnforcer, user User, domain Domain) bool {
		return CheckObject(e.Enforcer, user, domain, one, action)
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
