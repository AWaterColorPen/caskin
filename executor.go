package caskin

type executor struct {
	e        ienforcer
	mdb      MetaDB
	provider CurrentUserProvider
	factory  EntryFactory
	option   *Option
}

func (e *executor) filter(action Action, source interface{}) (interface{}, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	return Filter(e.e, u, d, action, e.factory.NewObject, source), nil
}

func (e *executor) filterWithNoError(user User, domain Domain, action Action, source interface{}) interface{} {
	return Filter(e.e, user, domain, action, e.factory.NewObject, source)
}

func (e *executor) check(action Action, one entry) error {
	u, d, err := e.provider.Get()
	if err != nil {
		return err
	}

	if ok := Check(e.e, u, d, action, e.factory.NewObject, one); !ok {
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

func (e *executor) checkParentEntryWrite(one parentEntry, take takeParentEntry) error {
	u, d, err := e.provider.Get()
	if err != nil {
		return err
	}

	if ok := Check(e.e, u, d, Write, e.factory.NewObject, one); !ok {
		return ErrNoWritePermission
	}

	for _, v := range []uint64{
		one.GetID(),
		one.GetParentID(),
	} {
		if v == 0 {
			continue
		}

		toCheck, err := take(v)
		if err != nil {
			return err
		}

		if ok := Check(e.e, u, d, Write, e.factory.NewObject, toCheck); !ok {
			return ErrNoWritePermission
		}
	}

	return nil
}

type takeParentEntry func(uint64) (parentEntry, error)
