package caskin

type executor struct {
	e        ienforcer
	mdb      MetaDB
	provider CurrentProvider
	factory  EntryFactory
	options  *Options
}

func (e *executor) filter(action Action, source interface{}) (interface{}, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	return Filter(e.e, u, d, action, source), nil
}

func (e *executor) filterWithNoßError(user User, domain Domain, action Action, source interface{}) []interface{} {
	return Filter(e.e, user, domain, action, source)
}

// 原来的写法
func (e *executor) filterWithNoError2(user User, domain Domain, action Action, source interface{}) interface{} {
	return Filter(e.e, user, domain, action, source)
}

func (e *executor) check(action Action, one entry) error {
	u, d, err := e.provider.Get()
	if err != nil {
		return err
	}

	if ok := Check(e.e, u, d, action, one); !ok {
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

	if ok := Check(e.e, u, d, Write, one); !ok {
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

		if ok := Check(e.e, u, d, Write, toCheck); !ok {
			return ErrNoWritePermission
		}
	}

	return nil
}

type takeParentEntry func(uint64) (parentEntry, error)
