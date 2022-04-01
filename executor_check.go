package caskin

import "fmt"

func (e *Executor) CheckObject(user User, domain Domain, one Object, action Action) error {
	if ok := Check(e.Enforcer, user, domain, one, action); ok {
		return nil
	}
	return fmt.Errorf("no %v permission", action)
}

func (e *Executor) CheckObjectData(user User, domain Domain, one ObjectData, action Action) error {
	if ok := Check(e.Enforcer, user, domain, one, action); ok {
		return nil
	}
	return fmt.Errorf("no %v permission", action)
}

func (e *Executor) SuperadminCheck(user User) error {
	ok, _ := e.Enforcer.IsSuperAdmin(user)
	if !ok {
		return ErrIsNotSuperAdmin
	}
	return nil
}

func (e *Executor) SuperadminAndSuperdomainCheck(user User, domain Domain) error {
	// TODO
	return nil
}

func (e *Executor) DBCreateCheck(item any) error {
	if err := e.DB.Take(item); err == nil {
		return ErrAlreadyExists
	}
	return nil
}

func (e *Executor) DBRecoverCheck(item any) error {
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

func (e *Executor) IDInterfaceUpdateCheck(item idInterface) error {
	if err := isValid(item); err != nil {
		return err
	}
	tmp := createByE(item)
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

func (e *Executor) ObjectCreateCheck(item Object, ty ObjectType) error {
	if err := e.DBCreateCheck(item); err != nil {
		return err
	}
	return e.ObjectDataWriteCheck(item, ty)
}

func (e *Executor) ObjectManageCheck(user User, domain Domain, item Object) error {
	if err := e.IDInterfaceModifyCheck(item); err != nil {
		return err
	}
	return e.CheckObject(user, domain, item, Manage)
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
	if err := e.IDInterfaceUpdateCheck(item); err != nil {
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

func (e *Executor) TreeNodeEntryUpdateCheck(item TreeNodeEntry, ty ObjectType) error {
	tmp1 := createByE(item)
	if err := e.ObjectDataUpdateCheck(item, tmp1, ty); err != nil {
		return err
	}
	if item.GetID() == item.GetParentID() {
		return ErrParentCanNotBeItself
	}
	if err := e.TreeNodeEntryParentCheck(tmp1); err != nil {
		return err
	}
	return e.TreeNodeEntryParentToDescendantCheck(item, tmp1, tmp2)
}

func (e *Executor) TreeNodeEntryParentCheck(item TreeNodeEntry) error {
	if isRoot(item) {
		return nil
	}
	parent := createByE(item)
	parent.SetID(item.GetParentID())
	if err := e.ObjectDataModifyCheck(parent); err != nil {
		return err
	}
	return e.TreeNodeEntryIsValidFamily(item, parent)
}

func (e *Executor) TreeNodeEntryIsValidFamily(data1, data2 ObjectData) error {
	o1 := data1.GetObject()
	o2 := data2.GetObject()
	if err := e.DB.Take(o1); err != nil {
		return ErrInValidParentObject
	}
	if err := e.DB.Take(o2); err != nil {
		return ErrInValidParentObject
	}
	if o1.GetObjectType() != o2.GetObjectType() {
		return ErrInValidParentObject
	}
	return nil
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
	if err := e.ObjectParentCheck(tmp); err != nil {
		return err
	}
	return e.ObjectTreeNodeParentToDescendantCheck(item, tmp)
}

func (e *Executor) ObjectParentCheck(user User, domain Domain, object Object) error {
	parent := createByE(object)
	parent.SetID(object.GetParentID())
	if err := e.ObjectManageCheck(user, domain, parent); err != nil {
		return err
	}
	if parent.GetObjectType() != object.GetObjectType() {
		return ErrInValidObjectType
	}
	return nil
}

func (e *Executor) rootObjectPermissionCheck() error {
	if err := e.SuperadminCheck(); err != nil {
		return ErrEmptyParentIdOrNotSuperadmin
	}
	return nil
}

func (e *Executor) ObjectTreeNodeParentToDescendantCheck(object Object, tmp Object) error {
	if object.GetParentID() == 0 || object.GetParentID() == tmp.GetParentID() {
		return nil
	}
	to := createByE(object)
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
