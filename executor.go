package caskin

import "github.com/ahmetb/go-linq/v3"

type executor struct {
	e        ienforcer
	mdb      MetaDB
	provider CurrentProvider
	factory  EntryFactory
	options  *Options
}

func (e *executor) newObject() treeNodeEntry {
	return e.factory.NewObject()
}

func (e *executor) newRole() treeNodeEntry {
	return e.factory.NewRole()
}

func (e *executor) createCheck(item interface{}) error {
	if err := e.mdb.Take(item); err == nil {
		return ErrAlreadyExists
	}
	return nil
}

func (e *executor) recoverCheck(item entry) error {
	if err := isValid(item); err != nil {
		return err
	}
	if err := e.mdb.Take(item); err == nil {
		return ErrAlreadyExists
	}
	if err := e.mdb.TakeUnscoped(item); err != nil {
		return ErrNotExists
	}
	return nil
}

func (e *executor) getOrModifyEntryCheck(item idInterface) error {
	if err := isValid(item); err != nil {
		return err
	}
	if err := e.mdb.Take(item); err != nil {
		return ErrNotExists
	}
	return nil
}

func (e *executor) deleteEntryCheck(item idInterface) error {
	return e.getOrModifyEntryCheck(item)
}

func (e *executor) getEntryCheck(item idInterface) error {
	return e.getOrModifyEntryCheck(item)
}

func (e *executor) modifyEntryCheck(item idInterface) error {
	return e.getOrModifyEntryCheck(item)
}

func (e *executor) updateEntryCheck(item idInterface, tmp entry) error {
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
	if err := e.createCheck(item); err != nil {
		return err
	}
	return e.check(item, Write)
}

func (e *executor) recoverObjectDataEntryCheck(item objectDataEntry) error {
	if err := e.recoverCheck(item); err != nil {
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

func (e *executor) parentEntryCheck(item treeNodeEntry, parentsFn parentsFn) error {
	user, domain, _ := e.provider.Get()
	parents := parentsFn(item, domain)
	for _, v := range parents {
		// special logic: normal user can't operate root object
		if v.GetID() == 0 {
			_, ok1 := item.(Object)
			ok2, _ := e.e.IsSuperAdmin(user)
			if ok1 && !ok2 {
				return ErrEmptyParentIdOrNotSuperadmin
			}
			return nil
		}

		if err := e.mdb.Take(v); err != nil {
			return err
		}
		if err := e.check(v, Write); err != nil {
			return err
		}
		// their object type should be same
		if u, ok := v.(Object); ok {
			w := item.(Object)
			if u.GetObjectType() != w.GetObjectType() {
				return ErrInValidObjectType
			}
		}
		if _, ok := v.(Object); !ok {
			u := v.(Role)
			w := item.(Role)
			if err := isValidFamily(w, u, e.mdb.Take); err != nil {
				return err
			}
		}

		if err := e.mdb.Take(v); err != nil {
			return err
		}
		if err := e.check(v, Write); err != nil {
			return err
		}
		// their object type should be same
		if u, ok := v.(Object); ok {
			w := item.(Object)
			if u.GetObjectType() != w.GetObjectType() {
				return ErrInValidObjectType
			}
		}
		// role is ObjectData, their object type should be same
		if _, ok := v.(Object); !ok {
			u := v.(Role)
			w := item.(Role)
			if err := isValidFamily(w, u, e.mdb.Take); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *executor) objectParentUpdater() *parentEntryUpdater {
	return &parentEntryUpdater{
		newEntry:  e.newObject,
		parentsFn: e.objectParentsFn(),
		addParentFn: func(p1 treeNodeEntry, p2 treeNodeEntry, domain Domain) error {
			return e.e.AddParentForObjectInDomain(p1.(Object), p2.(Object), domain)
		},
		deleteParentFn: func(p1 treeNodeEntry, p2 treeNodeEntry, domain Domain) error {
			return e.e.RemoveParentForObjectInDomain(p1.(Object), p2.(Object), domain)
		},
	}
}

func (e *executor) objectDeleteFn() deleteFn {
	return func(p treeNodeEntry, d Domain) error {
		if err := e.e.RemoveObjectInDomain(p.(Object), d); err != nil {
			return err
		}
		return e.mdb.DeleteObjectByID(p.GetID())
	}
}

func (e *executor) objectChildrenFn() childrenFn {
	return e.childrenOrParentFn(func(p treeNodeEntry, domain Domain) interface{} {
		return e.e.GetChildrenForObjectInDomain(p.(Object), domain)
	})
}

func (e *executor) objectParentsFn() parentsFn {
	return e.childrenOrParentFn(func(p treeNodeEntry, domain Domain) interface{} {
		return e.e.GetParentsForObjectInDomain(p.(Object), domain)
	})
}

func (e *executor) roleParentUpdater() *parentEntryUpdater {
	return &parentEntryUpdater{
		newEntry:  e.newRole,
		parentsFn: e.roleParentsFn(),
		addParentFn: func(p1 treeNodeEntry, p2 treeNodeEntry, domain Domain) error {
			return e.e.AddParentForRoleInDomain(p1.(Role), p2.(Role), domain)
		},
		deleteParentFn: func(p1 treeNodeEntry, p2 treeNodeEntry, domain Domain) error {
			return e.e.RemoveParentForRoleInDomain(p1.(Role), p2.(Role), domain)
		},
	}
}

func (e *executor) roleDeleteFn() deleteFn {
	return func(p treeNodeEntry, d Domain) error {
		if err := e.e.RemoveRoleInDomain(p.(Role), d); err != nil {
			return err
		}
		return e.mdb.DeleteRoleByID(p.GetID())
	}
}

func (e *executor) roleChildrenFn() childrenFn {
	return e.childrenOrParentFn(func(p treeNodeEntry, domain Domain) interface{} {
		return e.e.GetChildrenForRoleInDomain(p.(Role), domain)
	})
}

func (e *executor) roleParentsFn() parentsFn {
	return e.childrenOrParentFn(func(p treeNodeEntry, domain Domain) interface{} {
		return e.e.GetParentsForRoleInDomain(p.(Role), domain)
	})
}

func (e *executor) childrenOrParentFn(fn func(treeNodeEntry, Domain) interface{}) childrenFn {
	return func(p treeNodeEntry, domain Domain) []treeNodeEntry {
		var out []treeNodeEntry
		linq.From(fn(p, domain)).ToSlice(&out)
		return out
	}
}

func (e *executor) parentEntryFlowHandler(item treeNodeEntry,
	check func(objectDataEntry) error,
	newEntry func() treeNodeEntry,
	fn func(Domain) error) error {
	if err := check(item); err != nil {
		return err
	}

	user, domain, err := e.provider.Get()
	if err != nil {
		return err
	}

	if item.GetParentID() != 0 {
		if err := e.parentEntryCheck(item, singleParentsFunc(item, newEntry)); err != nil {
			return err
		}
	} else {
		if ok, _ := e.e.IsSuperAdmin(user); !ok {
			return ErrCanNotOperateRootObjectWithoutSuperadmin
		}
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

