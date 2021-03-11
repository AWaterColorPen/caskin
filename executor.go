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

func (e *executor) createEntryCheck(item entry) error {
	if err := e.mdb.Take(item); err == nil {
		return ErrAlreadyExists
	}
	return nil
}

func (e *executor) recoverEntryCheck(item entry) error {
	if err := e.mdb.Take(item); err == nil {
		return ErrAlreadyExists
	}
	return e.mdb.TakeUnscoped(item)
}

func (e *executor) getOrModifyEntryCheck(item entry) error {
	if err := isValid(item); err != nil {
		return err
	}
	if err := e.mdb.Take(item); err != nil {
		return ErrNotExists
	}
	return nil
}

func (e *executor) deleteEntryCheck(item entry) error {
	return e.getOrModifyEntryCheck(item)
}

func (e *executor) getEntryCheck(item entry) error {
	return e.getOrModifyEntryCheck(item)
}

func (e *executor) modifyEntryCheck(item entry) error {
	return e.getOrModifyEntryCheck(item)
}

func (e *executor) updateEntryCheck(item entry, tmp entry) error {
	if err := isValid(item); err != nil {
		return err
	}

	tmp.SetID(item.GetID())
	if err := e.mdb.Take(tmp); err != nil {
		return ErrNotExists
	}

	return nil
}

func (e *executor) createObjectDataEntryCheck(item objectDataEntry) error {
	if err := e.createEntryCheck(item); err != nil {
		return err
	}
	return e.check(item, Write)
}

func (e *executor) recoverObjectDataEntryCheck(item objectDataEntry) error {
	if err := e.recoverEntryCheck(item); err != nil {
		return err
	}
	return e.check(item, Write)
}

func (e *executor) getOrModifyObjectDataEntryCheck(item objectDataEntry, actions ...Action) error {
	if err := e.getOrModifyEntryCheck(item); err != nil {
		return err
	}
	for _, action := range actions {
		if err := e.check(item, action); err != nil {
			return err
		}
	}
	return nil
}

func (e *executor) deleteObjectDataEntryCheck(item objectDataEntry) error {
	return e.getOrModifyObjectDataEntryCheck(item, Write)
}

func (e *executor) getObjectDataEntryCheck(item objectDataEntry) error {
	return e.getOrModifyObjectDataEntryCheck(item, Read)
}

func (e *executor) modifyObjectDataEntryCheck(item objectDataEntry) error {
	return e.getOrModifyObjectDataEntryCheck(item, Write)
}

func (e *executor) updateObjectDataEntryCheck(item objectDataEntry, tmp objectDataEntry) error {
	if err := e.updateEntryCheck(item, tmp); err != nil {
		return err
	}

	return e.check(tmp, Write)
}

func (e *executor) parentEntryCheck(item parentEntry, parentsFn parentsFn) error {
	_, domain, _ := e.provider.Get()
	parents := parentsFn(item, domain)
	for _, v := range parents {
		if err := e.mdb.Take(v); err != nil {
			return err
		}
		if err := e.check(v, Write); err != nil {
			return err
		}
	}
	return nil
}

func (e *executor) objectParentUpdater() *parentEntryUpdater {
	return &parentEntryUpdater{
		newEntry:  e.newObject,
		parentsFn: e.objectParentsFn(),
		addParentFn: func(p1 parentEntry, p2 parentEntry, domain Domain) error {
			return e.e.AddParentForObjectInDomain(p1.(Object), p2.(Object), domain)
		},
		deleteParentFn: func(p1 parentEntry, p2 parentEntry, domain Domain) error {
			return e.e.RemoveParentForObjectInDomain(p1.(Object), p2.(Object), domain)
		},
	}
}

func (e *executor) objectDeleteFn() deleteFn {
	return func(p parentEntry, d Domain) error {
		if err := e.e.RemoveObjectInDomain(p.(Object), d); err != nil {
			return err
		}
		return e.mdb.DeleteObjectByID(p.GetID())
	}
}

func (e *executor) objectChildrenFn() childrenFn {
	return e.childrenOrParentFn(func(p parentEntry, domain Domain) interface{} {
		return e.e.GetChildrenForObjectInDomain(p.(Object), domain)
	})
}

func (e *executor) objectParentsFn() parentsFn {
	return e.childrenOrParentFn(func(p parentEntry, domain Domain) interface{} {
		return e.e.GetParentsForObjectInDomain(p.(Object), domain)
	})
}

func (e *executor) roleParentUpdater() *parentEntryUpdater {
	return &parentEntryUpdater{
		newEntry:  e.newRole,
		parentsFn: e.roleParentsFn(),
		addParentFn: func(p1 parentEntry, p2 parentEntry, domain Domain) error {
			return e.e.AddParentForRoleInDomain(p1.(Role), p2.(Role), domain)
		},
		deleteParentFn: func(p1 parentEntry, p2 parentEntry, domain Domain) error {
			return e.e.RemoveParentForRoleInDomain(p1.(Role), p2.(Role), domain)
		},
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

func (e *executor) roleChildrenFn() childrenFn {
	return e.childrenOrParentFn(func(p parentEntry, domain Domain) interface{} {
		return e.e.GetChildrenForRoleInDomain(p.(Role), domain)
	})
}

func (e *executor) roleParentsFn() parentsFn {
	return e.childrenOrParentFn(func(p parentEntry, domain Domain) interface{} {
		return e.e.GetParentsForRoleInDomain(p.(Role), domain)
	})
}

func (e *executor) childrenOrParentFn(fn func(parentEntry, Domain) interface{}) childrenFn {
	return func(p parentEntry, domain Domain) []parentEntry {
		var out []parentEntry
		linq.From(fn(p, domain)).ToSlice(&out)
		return out
	}
}

func (e *executor) parentEntryFlowHandler(item parentEntry,
	check func(objectDataEntry) error,
	newEntry func() parentEntry,
	fn func(Domain) error) error {
	if err := check(item); err != nil {
		return err
	}

	if err := e.parentEntryCheck(item, singleParentsFunc(item, newEntry)); err != nil {
		return err
	}

	_, domain, err := e.provider.Get()
	if err != nil {
		return err
	}
	item.SetDomainID(domain.GetID())
	return fn(domain)
}

func (e *executor) filter(action Action, source interface{}) ([]interface{}, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}
	return Filter(e.e, u, d, action, source), nil
}

func (e *executor) filterWithNoError(user User, domain Domain, action Action, source interface{}) []interface{} {
	return Filter(e.e, user, domain, action, source)
}

func (e *executor) check(one ObjectData, action Action) error {
	u, d, err := e.provider.Get()
	if err != nil {
		return err
	}

	if ok := Check(e.e, u, d, one, action); !ok {
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

type objectDataEntry interface {
	entry
	ObjectData
}
