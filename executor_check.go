package caskin

import "fmt"

func (e *server) CheckObject(user User, domain Domain, one Object, action Action) error {
	if ok := Check(e.Enforcer, user, domain, one, action); ok {
		return nil
	}
	return fmt.Errorf("no %v permission", action)
}

func (e *server) CheckObjectData(user User, domain Domain, one ObjectData, action Action) error {
	if ok := Check(e.Enforcer, user, domain, one, action); ok {
		return nil
	}
	return fmt.Errorf("no %v permission", action)
}

func (e *server) SuperadminCheck(user User) error {
	ok, _ := e.Enforcer.IsSuperAdmin(user)
	if !ok {
		return ErrIsNotSuperAdmin
	}
	return nil
}

func (e *server) DBCreateCheck(item any) error {
	if err := e.DB.Take(item); err == nil {
		return ErrAlreadyExists
	}
	return nil
}

func (e *server) DBRecoverCheck(item any) error {
	if err := e.DB.Take(item); err == nil {
		return ErrAlreadyExists
	}
	if err := e.DB.TakeUnscoped(item); err != nil {
		return ErrNotExists
	}
	return nil
}

func (e *server) IDInterfaceDeleteCheck(item idInterface) error {
	return e.IDInterfaceValidAndExistsCheck(item)
}

func (e *server) IDInterfaceUpdateCheck(item, old idInterface) error {
	if err := isValid(item); err != nil {
		return err
	}
	old.SetID(item.GetID())
	if err := e.DB.Take(old); err != nil {
		return ErrNotExists
	}
	return nil
}

func (e *server) IDInterfaceValidAndExistsCheck(item idInterface) error {
	if err := isValid(item); err != nil {
		return err
	}
	if err := e.DB.Take(item); err != nil {
		return ErrNotExists
	}
	return nil
}

func (e *server) IDInterfaceGetCheck(item idInterface) error {
	return e.IDInterfaceValidAndExistsCheck(item)
}

func (e *server) IDInterfaceModifyCheck(item idInterface) error {
	return e.IDInterfaceValidAndExistsCheck(item)
}

func (e *server) ObjectManageCheck(user User, domain Domain, item Object) error {
	if err := e.IDInterfaceModifyCheck(item); err != nil {
		return err
	}
	return e.CheckObject(user, domain, item, Manage)
}

func (e *server) ObjectParentCheck(user User, domain Domain, object Object) error {
	if object.GetParentID() == 0 {
		return ErrCantOperateRootObject
	}
	parent := newByE(object)
	parent.SetID(object.GetParentID())
	if err := e.ObjectManageCheck(user, domain, parent); err != nil {
		return err
	}
	if parent.GetObjectType() != object.GetObjectType() {
		return ErrInValidObjectType
	}
	return nil
}

func (e *server) ObjectParentToDescendantCheck(domain Domain, object Object, old Object) error {
	if object.GetParentID() == 0 || object.GetParentID() == old.GetParentID() {
		return nil
	}
	to := newByE(object)
	to.SetID(object.GetParentID())
	if ok, _ := e.Enforcer.EnforceObject(object, to, domain); ok {
		return ErrParentToDescendant
	}
	return nil
}

func (e *server) ObjectUpdateCheck(user User, domain Domain, object Object) error {
	old := newByE(object)
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

func (e *server) RoleParentCheck(role Role) error {
	if role.GetParentID() == 0 {
		return nil
	}
	parent := newByE(role)
	parent.SetID(role.GetParentID())
	if err := e.IDInterfaceModifyCheck(parent); err != nil {
		return err
	}
	if role.GetObject().GetID() != parent.GetObject().GetID() {
		return ErrParentCanNotDiff
	}
	return nil
}

func (e *server) RoleParentToDescendantCheck(domain Domain, role Role, old Role) error {
	if role.GetParentID() == 0 || role.GetParentID() == old.GetParentID() {
		return nil
	}
	to := newByE(role)
	to.SetID(role.GetParentID())
	if ok, _ := e.Enforcer.EnforceRole(role, to, domain); ok {
		return ErrParentToDescendant
	}
	return nil
}

func (e *server) RoleUpdateCheck(user User, domain Domain, role Role) error {
	old := newByE(role)
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
