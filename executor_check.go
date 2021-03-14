package caskin

func (e *Executor) DBCreateCheck(item interface{}) error {
	if err := e.DB.Take(item); err == nil {
		return ErrAlreadyExists
	}
	return nil
}

func (e *Executor) DBRecoverCheck(item interface{}) error {
	if err := e.DB.Take(item); err == nil {
		return ErrAlreadyExists
	}
	if err := e.DB.TakeUnscoped(item); err != nil {
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
	if err := e.DB.Take(tmp); err != nil {
		return ErrNotExists
	}
	return nil
}

func (e *Executor) IDInterfaceValidAndExistsCheck(item idInterface) error {
	if err := isValid(item); err != nil {
		return err
	}
	if err := e.DB.Take(item); err != nil {
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

func (e *Executor) ObjectDataCreateCheck(item ObjectData) error {
	if err := e.DBCreateCheck(item); err != nil {
		return err
	}
	return e.check(item, Write)
}

func (e *Executor) ObjectDataRecoverCheck(item ObjectData) error {
	if err := e.DBRecoverCheck(item); err != nil {
		return err
	}
	return e.check(item, Write)
}

func (e *Executor) ObjectDataDeleteCheck(item ObjectData) error {
	if err := e.IDInterfaceDeleteCheck(item); err != nil {
		return err
	}
	return e.check(item, Write)
}

func (e *Executor) ObjectDataUpdateCheck(item ObjectData, tmp ObjectData) error {
	if err := e.IDInterfaceUpdateCheck(item, tmp); err != nil {
		return err
	}
	return e.check(tmp, Write)
}

func (e *Executor) ObjectDataGetCheck(item ObjectData) error {
	if err := e.IDInterfaceGetCheck(item); err != nil {
		return err
	}
	return e.check(item, Read)
}

func (e *Executor) ObjectDataModifyCheck(item ObjectData) error {
	if err := e.IDInterfaceModifyCheck(item); err != nil {
		return err
	}
	return e.check(item, Write)
}

func (e *Executor) treeNodeParentCheck(takenItem treeNodeEntry, newEntry func() treeNodeEntry) error {
	if isRoot(takenItem) {
		return nil
	}
	parent := newEntry()
	parent.SetID(takenItem.GetParentID())
	if err := e.ObjectDataModifyCheck(parent); err != nil {
		return err
	}
	return isValidFamily(takenItem, parent, e.DB.Take)
}

func (e *Executor) treeNodeEntryCheckFlow(item treeNodeEntry, check func(ObjectData) error, newEntry func() treeNodeEntry) error {
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
	return nil
}

func (e *Executor) objectTreeNodeParentCheck(takenObject Object) error {
	if isRoot(takenObject) {
		return e.rootObjectPermissionCheck()
	}
	parent := e.factory.NewObject()
	parent.SetID(takenObject.GetParentID())
	if err := e.ObjectDataModifyCheck(parent); err != nil {
		return err
	}
	if parent.GetObjectType() != takenObject.GetObjectType() {
		return ErrInValidObjectType
	}
	return nil
}

func (e *Executor) objectCheckFlow(object Object, check func(ObjectData) error) error {
	if err := check(object); err != nil {
		return err
	}

	if err := e.objectTreeNodeParentCheck(object); err != nil {
		return err
	}

	_, domain, err := e.provider.Get()
	if err != nil {
		return err
	}

	object.SetDomainID(domain.GetID())
	return nil
}

func (e *Executor) IsSuperadminCheck() error {
	user, _, err := e.provider.Get()
	if err != nil {
		return err
	}
	ok, _ := e.e.IsSuperAdmin(user)
	if !ok {
		return ErrIsNotSuperAdmin
	}
	return nil
}

func (e *Executor) rootObjectPermissionCheck() error {
	if err := e.IsSuperadminCheck(); err != nil {
		return ErrEmptyParentIdOrNotSuperadmin
	}
	return nil
}
