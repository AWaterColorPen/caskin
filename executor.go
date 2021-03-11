package caskin

import "github.com/ahmetb/go-linq/v3"

type executor struct {
	e        ienforcer
	mdb      MetaDB
	provider CurrentProvider
	factory  EntryFactory
	options  *Options
}

func (e *executor) newObject() parentEntry {
	return e.factory.NewObject()
}

func (e *executor) newRole() parentEntry {
	return e.factory.NewRole()
}

func (e *executor) createCheck(item ObjectData) error {
	if err := e.mdb.Take(item); err == nil {
		return ErrAlreadyExists
	}
	return e.check(Write, item)
}

func (e *executor) deleteCheck(item objectDataEntry) error {
	if err := isValid(item); err != nil {
		return err
	}
	if err := e.mdb.Take(item); err != nil {
		return ErrNotExists
	}
	return e.check(Write, item)
}

func (e *executor) getOrModifyCheck(item objectDataEntry, actions ...Action) error {
	if err := isValid(item); err != nil {
		return err
	}

	if err := e.mdb.Take(item); err != nil {
		return ErrNotExists
	}

	for _, action := range actions {
		if err := e.check(action, item); err != nil {
			return err
		}
	}

	return nil
}

func (e *executor) objectDeleteFn() deleteFn {
	return func(p parentEntry, d Domain) error {
		if err := e.e.RemoveObjectInDomain(p.(Object), d); err != nil {
			return err
		}
		return e.mdb.DeleteObjectByID(p.GetID())
	}
}

func (e *executor) roleDeleteFn() deleteFn {
	return func(p parentEntry, d Domain) error {
		if err := e.e.RemoveRoleInDomain(p.(Role), d); err != nil {
			return err
		}
		return e.mdb.DeleteRoleByID(p.GetID())
	}
}

func (e *executor) objectChildrenFn() childrenFn {
	return e.childrenFn(func(p parentEntry, domain Domain) interface{} {
		return e.e.GetChildrenForObjectInDomain(p.(Object), domain)
	})
}

func (e *executor) roleChildrenFn() childrenFn {
	return e.childrenFn(func(p parentEntry, domain Domain) interface{} {
		return e.e.GetChildrenForRoleInDomain(p.(Role), domain)
	})
}

func (e *executor) childrenFn(fn func(parentEntry, Domain) interface{}) childrenFn {
	return func(p parentEntry, domain Domain) []parentEntry {
		var out []parentEntry
		children := fn(p, domain)
		linq.From(children).ToSlice(&out)
		return out
	}
}

func (e *executor) filter(action Action, source interface{}) ([]interface{}, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	return Filter(e.e, u, d, action, source), nil
}

// 原来的写法
func (e *executor) filter2(action Action, source interface{}) (interface{}, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	return Filter(e.e, u, d, action, source), nil
}

func (e *executor) filterWithNoError(user User, domain Domain, action Action, source interface{}) []interface{} {
	return Filter(e.e, user, domain, action, source)
}

// filterWithNoError2 original code
func (e *executor) filterWithNoError2(user User, domain Domain, action Action, source interface{}) interface{} {
	return Filter(e.e, user, domain, action, source)
}

func (e *executor) check(action Action, one ObjectData) error {
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

type objectDataEntry interface {
	entry
	ObjectData
}
