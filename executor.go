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

func (e *executor) recoverCheck(item interface{}) error {
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

func (e *executor) treeNodeParentCheck(takenItem treeNodeEntry, newEntry func() treeNodeEntry) error {
	user, _, _ := e.provider.Get()

	// special logic: normal user can't operate root object
	if isObjectRoot(takenItem) {
		ok, _ := e.e.IsSuperAdmin(user)
		if !ok {
			return ErrCanNotOperateRootObjectWithoutSuperadmin
		}

		return nil
	}

	pid := takenItem.GetParentID()
	parent := newEntry()
	parent.SetID(pid)

	if err := e.getOrModifyObjectDataEntryCheck(parent, Write); err != nil {
		return err
	}

	// TODO hanshu
	// special logic: their object type should be same
	if u, ok := parent.(Object); ok {
		w := takenItem.(Object)
		if u.GetObjectType() != w.GetObjectType() {
			return ErrInValidObjectType
		}
	}

	// TODO hanshu
	// special logic:
	if _, ok := parent.(Object); !ok {
		u := parent.(Role)
		w := takenItem.(Role)
		if err := isValidFamily(w, u, e.mdb.Take); err != nil {
			return err
		}
	}

	return nil
}

func (e *executor) objectParentUpdater() *parentEntryUpdater {
	return &parentEntryUpdater{
		newEntry:  e.newObject,
		parentGetFn: e.objectParentsFn(),
		parentAddFn: func(p1 treeNodeEntry, p2 treeNodeEntry, domain Domain) error {
			return e.e.AddParentForObjectInDomain(p1.(Object), p2.(Object), domain)
		},
		parentDelFn: func(p1 treeNodeEntry, p2 treeNodeEntry, domain Domain) error {
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
	return e.childrenOrParentGetFn(func(p treeNodeEntry, domain Domain) interface{} {
		return e.e.GetChildrenForObjectInDomain(p.(Object), domain)
	})
}

func (e *executor) objectParentsFn() parentGetFn {
	return e.childrenOrParentGetFn(func(p treeNodeEntry, domain Domain) interface{} {
		return e.e.GetParentsForObjectInDomain(p.(Object), domain)
	})
}

func (e *executor) roleParentUpdater() *parentEntryUpdater {
	return &parentEntryUpdater{
		newEntry:  e.newRole,
		parentGetFn: e.roleParentsFn(),
		parentAddFn: func(p1 treeNodeEntry, p2 treeNodeEntry, domain Domain) error {
			return e.e.AddParentForRoleInDomain(p1.(Role), p2.(Role), domain)
		},
		parentDelFn: func(p1 treeNodeEntry, p2 treeNodeEntry, domain Domain) error {
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
	return e.childrenOrParentGetFn(func(p treeNodeEntry, domain Domain) interface{} {
		return e.e.GetChildrenForRoleInDomain(p.(Role), domain)
	})
}

func (e *executor) roleParentsFn() parentGetFn {
	return e.childrenOrParentGetFn(func(p treeNodeEntry, domain Domain) interface{} {
		return e.e.GetParentsForRoleInDomain(p.(Role), domain)
	})
}

func (e *executor) childrenOrParentGetFn(fn func(treeNodeEntry, Domain) interface{}) childrenFn {
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

	if err := e.treeNodeParentCheck(item, newEntry); err != nil {
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
