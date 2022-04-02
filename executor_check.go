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

func (e *Executor) IDInterfaceUpdateCheck(item, old idInterface) error {
	if err := isValid(item); err != nil {
		return err
	}
	old.SetID(item.GetID())
	if err := e.DB.Take(old); err != nil {
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

func (e *Executor) ObjectManageCheck(user User, domain Domain, item Object) error {
	if err := e.IDInterfaceModifyCheck(item); err != nil {
		return err
	}
	return e.CheckObject(user, domain, item, Manage)
}

func (e *Executor) ObjectParentCheck(user User, domain Domain, object Object) error {
	if object.GetParentID() == 0 {
		return ErrCantOperateRootObject
	}
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

func (e *Executor) ObjectParentToDescendantCheck(domain Domain, object Object, old Object) error {
	if object.GetParentID() == 0 || object.GetParentID() == old.GetParentID() {
		return nil
	}
	to := createByE(object)
	to.SetID(object.GetParentID())
	if ok, _ := e.Enforcer.EnforceObject(object, to, domain); ok {
		return ErrParentToDescendant
	}
	return nil
}

func (e *Executor) ObjectUpdateCheck(user User, domain Domain, object Object) error {
	old := createByE(object)
	if err := e.IDInterfaceUpdateCheck(object, old); err != nil {
		return err
	}
	if object.GetID() == object.GetParentID() {
		return ErrParentCanNotBeItself
	}
	if object.GetObjectType() != "" && object.GetObjectType() != old.GetObjectType() {
		return ErrCantChangeObjectType
	}
	if err := e.ObjectManageCheck(user, domain, object); err != nil {
		return err
	}
	return e.ObjectParentToDescendantCheck(domain, object, old)
}

func (e *Executor) RoleParentCheck(role Role) error {
	if role.GetParentID() == 0 {
		return nil
	}
	parent := createByE(role)
	parent.SetID(role.GetParentID())
	if err := e.IDInterfaceModifyCheck(parent); err != nil {
		return err
	}
	if role.GetObject().GetID() != parent.GetObject().GetID() {
		return ErrParentCanNotDiff
	}
	return nil
}

func (e *Executor) RoleParentToDescendantCheck(domain Domain, role Role, old Role) error {
	if role.GetParentID() == 0 || role.GetParentID() == old.GetParentID() {
		return nil
	}
	to := createByE(role)
	to.SetID(role.GetParentID())
	if ok, _ := e.Enforcer.EnforceRole(role, to, domain); ok {
		return ErrParentToDescendant
	}
	return nil
}

func (e *Executor) RoleUpdateCheck(user User, domain Domain, role Role) error {
	old := createByE(role)
	if err := e.IDInterfaceUpdateCheck(role, old); err != nil {
		return err
	}
	if role.GetID() == role.GetParentID() {
		return ErrParentCanNotBeItself
	}
	if err := e.ObjectDataWriteCheck(user, domain, old, ObjectTypeRole); err != nil {
		return err
	}
	if role.GetObject().GetID() != old.GetObject().GetID() {
		return e.ObjectDataWriteCheck(user, domain, role, ObjectTypeRole)
	}
	return e.RoleParentToDescendantCheck(domain, role, old)
}
