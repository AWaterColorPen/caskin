package caskin

func (e *Executor) DBCreateCheck(item interface{}) error {
	if err := e.db.Take(item); err == nil {
		return ErrAlreadyExists
	}
	return nil
}

func (e *Executor) DBRecoverCheck(item interface{}) error {
	if err := e.db.Take(item); err == nil {
		return ErrAlreadyExists
	}
	if err := e.db.TakeUnscoped(item); err != nil {
		return ErrNotExists
	}
	return nil
}

func (e *Executor) IDInterfaceDeleteCheck(item idInterface) error {
	return e.IDInterfaceValidAndExistsCheck(item)
}

func (e *Executor) IDInterfaceUpdateCheck(item idInterface, tmp idInterface) error {
	if err := isValid(item); err != nil {
		return err
	}
	tmp.SetID(item.GetID())
	if err := e.db.Take(tmp); err != nil {
		return ErrNotExists
	}
	return nil
}

func (e *Executor) IDInterfaceValidAndExistsCheck(item idInterface) error {
	if err := isValid(item); err != nil {
		return err
	}
	if err := e.db.Take(item); err != nil {
		return ErrNotExists
	}
	return nil
}

func (e *Executor) IDInterfaceGetCheck(item idInterface) error {
	return e.IDInterfaceValidAndExistsCheck(item)
}

func (e *Executor) IDInterfaceModifyCheck(item idInterface) error {
	return e.IDInterfaceValidAndExistsCheck(item)
}

func (e *Executor) createObjectDataEntryCheck(item ObjectData) error {
	if err := e.DBCreateCheck(item); err != nil {
		return err
	}
	return e.check(item, Write)
}

func (e *Executor) recoverObjectDataEntryCheck(item ObjectData) error {
	if err := e.DBRecoverCheck(item); err != nil {
		return err
	}
	return e.check(item, Write)
}

func (e *Executor) getOrModifyObjectDataEntryCheck(item ObjectData, actions ...Action) error {
	if err := e.IDInterfaceValidAndExistsCheck(item); err != nil {
		return err
	}
	for _, action := range actions {
		if err := e.check(item, action); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) deleteObjectDataEntryCheck(item ObjectData) error {
	return e.getOrModifyObjectDataEntryCheck(item, Write)
}

func (e *Executor) getObjectDataEntryCheck(item ObjectData) error {
	return e.getOrModifyObjectDataEntryCheck(item, Read)
}

func (e *Executor) modifyObjectDataEntryCheck(item ObjectData) error {
	return e.getOrModifyObjectDataEntryCheck(item, Write)
}

func (e *Executor) updateObjectDataEntryCheck(item ObjectData, tmp ObjectData) error {
	if err := e.IDInterfaceUpdateCheck(item, tmp); err != nil {
		return err
	}

	return e.check(tmp, Write)
}

func (e *Executor) treeNodeParentCheck(takenItem treeNodeEntry, newEntry func() treeNodeEntry) error {
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
		if err := isValidFamily(w, u, e.db.Take); err != nil {
			return err
		}
	}

	return nil
}
