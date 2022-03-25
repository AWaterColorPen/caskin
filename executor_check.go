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

func (e *Executor) IsSuperadminAndSuperdomainCheck() error {
	// TODO
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

func (e *Executor) ObjectDataWriteCheck(item ObjectData, ty ObjectType) error {
	if err := e.checkObjectData(item, Write); err != nil {
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
	return e.ObjectDataWriteCheck(item, ty)
}

func (e *Executor) ObjectDataRecoverCheck(item ObjectData) error {
	if err := e.DBRecoverCheck(item); err != nil {
		return err
	}
	return e.checkObjectData(item, Write)
}

func (e *Executor) ObjectDataDeleteCheck(item ObjectData) error {
	if err := e.IDInterfaceDeleteCheck(item); err != nil {
		return err
	}
	return e.checkObjectData(item, Write)
}

func (e *Executor) ObjectDataUpdateCheck(item ObjectData, tmp ObjectData, ty ObjectType) error {
	if err := e.IDInterfaceUpdateCheck(item, tmp); err != nil {
		return err
	}
	if err := e.ObjectDataWriteCheck(tmp, ty); err != nil {
		return err
	}
	return e.ObjectDataWriteCheck(item, ty)
}

func (e *Executor) ObjectDataGetCheck(item ObjectData) error {
	if err := e.IDInterfaceGetCheck(item); err != nil {
		return err
	}
	return e.checkObjectData(item, Read)
}

func (e *Executor) ObjectDataModifyCheck(item ObjectData) error {
	if err := e.IDInterfaceModifyCheck(item); err != nil {
		return err
	}
	return e.checkObjectData(item, Write)
}

func (e *Executor) TreeNodeEntryUpdateCheck(item TreeNodeEntry, tmp1 TreeNodeEntry, tmp2 TreeNodeEntry, ty ObjectType) error {
	if err := e.ObjectDataUpdateCheck(item, tmp1, ty); err != nil {
		return err
	}
	if item.GetID() == item.GetParentID() {
		return ErrParentCanNotBeItself
	}
	if err := e.TreeNodeEntryParentCheck(tmp1, tmp2); err != nil {
		return err
	}
	return e.TreeNodeEntryParentToDescendantCheck(item, tmp1, tmp2)
}

func (e *Executor) TreeNodeEntryParentCheck(item TreeNodeEntry, parent TreeNodeEntry) error {
	if isRoot(item) {
		return nil
	}
	parent.SetID(item.GetParentID())
	if err := e.ObjectDataModifyCheck(parent); err != nil {
		return err
	}
	return isValidFamily(item, parent, e.DB.Take)
}

func (e *Executor) ObjectTreeNodeUpdateCheck(item Object, tmp Object) error {
	if err := e.ObjectDataUpdateCheck(item, tmp, ObjectTypeObject); err != nil {
		return err
	}
	if item.GetID() != 0 && item.GetID() == item.GetParentID() {
		return ErrParentCanNotBeItself
	}
	if item.GetObjectType() != tmp.GetObjectType() {
		return ErrCantChangeObjectType
	}
	if err := e.ObjectTreeNodeParentCheck(tmp); err != nil {
		return err
	}
	return e.ObjectTreeNodeParentToDescendantCheck(item, tmp)
}

func (e *Executor) ObjectTreeNodeParentCheck(object Object) error {
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

func (e *Executor) ObjectTreeNodeParentToDescendantCheck(object Object, tmp Object) error {
	if object.GetParentID() == 0 || object.GetParentID() == tmp.GetParentID() {
		return nil
	}
	to := e.factory.NewObject()
	to.SetID(object.GetParentID())
	domain := e.factory.NewDomain()
	domain.SetID(tmp.GetDomainID())
	if ok, _ := e.Enforcer.EnforceObject(object, to, domain); ok {
		return ErrParentToDescendant
	}
	return nil
}

func (e *Executor) TreeNodeEntryParentToDescendantCheck(item TreeNodeEntry, tmp TreeNodeEntry, parent TreeNodeEntry) error {
	if item.GetParentID() == 0 || item.GetParentID() == tmp.GetParentID() {
		return nil
	}
	parent.SetID(item.GetParentID())
	domain := e.factory.NewDomain()
	domain.SetID(tmp.GetDomainID())
	if ok, _ := e.Enforcer.EnforceRole(item.(Role), parent.(Role), domain); ok {
		return ErrParentToDescendant
	}
	return nil
}
