package caskin

func (e *Executor) IsSuperadminCheck() error {
	user, _, err := e.provider.Get()
	if err != nil {
		return err
	}
	ok, _ := e.Enforcer.IsSuperAdmin(user)
	if !ok {
		return ErrIsNotSuperAdmin
	}
	return nil
}

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

func (e *Executor) objectDataWriteCheck(item ObjectData, ty ObjectType) error {
	if err := e.check(item, Write); err != nil {
		return err
	}
	o := item.GetObject()
	if err := e.DB.Take(o); err != nil {
		return ErrInValidObject
	}
	if o.GetObjectType() != ty {
		return ErrInValidObjectType
	}
	return nil
}

func (e *Executor) ObjectDataCreateCheck(item ObjectData, ty ObjectType) error {
	if err := e.DBCreateCheck(item); err != nil {
		return err
	}
	return e.objectDataWriteCheck(item, ty)
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

func (e *Executor) ObjectDataUpdateCheck(item ObjectData, tmp ObjectData, ty ObjectType) error {
	if err := e.IDInterfaceUpdateCheck(item, tmp); err != nil {
		return err
	}
	return e.objectDataWriteCheck(item, ty)
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

func (e *Executor) treeNodeEntryUpdateCheck(item treeNodeEntry, tmp1 treeNodeEntry, tmp2 treeNodeEntry, ty ObjectType) error {
	if err := e.ObjectDataUpdateCheck(item, tmp1, ty); err != nil {
		return err
	}
	if item.GetID() == item.GetParentID() {
		return ErrParentCanNotBeItself
	}
	return e.treeNodeEntryParentCheck(tmp1, tmp2)
}

func (e *Executor) treeNodeEntryParentCheck(item treeNodeEntry, parent treeNodeEntry) error {
	if isRoot(item) {
		return nil
	}
	parent.SetID(item.GetParentID())
	if err := e.ObjectDataModifyCheck(parent); err != nil {
		return err
	}
	return isValidFamily(item, parent, e.DB.Take)
}

func (e *Executor) objectTreeNodeUpdateCheck(item Object, tmp Object) error {
	if err := e.ObjectDataUpdateCheck(item, tmp, ObjectTypeObject); err != nil {
		return err
	}
	if item.GetID() != 0 && item.GetID() == item.GetParentID() {
		return ErrParentCanNotBeItself
	}
	if item.GetObjectType() != tmp.GetObjectType() {
		return ErrCantChangeObjectType
	}
	return e.objectTreeNodeParentCheck(tmp)
}

func (e *Executor) objectTreeNodeParentCheck(object Object) error {
	if isRoot(object) {
		return e.rootObjectPermissionCheck()
	}
	parent := e.factory.NewObject()
	parent.SetID(object.GetParentID())
	if err := e.ObjectDataModifyCheck(parent); err != nil {
		return err
	}
	if parent.GetObjectType() != object.GetObjectType() {
		return ErrInValidObjectType
	}
	return nil
}

func (e *Executor) rootObjectPermissionCheck() error {
	if err := e.IsSuperadminCheck(); err != nil {
		return ErrEmptyParentIdOrNotSuperadmin
	}
	return nil
}
