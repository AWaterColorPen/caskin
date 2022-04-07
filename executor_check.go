package caskin

import "fmt"

func (s *server) CheckObject(user User, domain Domain, one Object, action Action) error {
	if ok := Check(s.Enforcer, user, domain, one, action); ok {
		return nil
	}
	return fmt.Errorf("no %v permission", action)
}

func (s *server) CheckObjectData(user User, domain Domain, one ObjectData, action Action) error {
	if ok := Check(s.Enforcer, user, domain, one, action); ok {
		return nil
	}
	return fmt.Errorf("no %v permission", action)
}

func (s *server) SuperadminCheck(user User) error {
	ok, _ := s.Enforcer.IsSuperAdmin(user)
	if !ok {
		return ErrIsNotSuperAdmin
	}
	return nil
}

func (s *server) DBCreateCheck(item any) error {
	if err := s.DB.Take(item); err == nil {
		return ErrAlreadyExists
	}
	return nil
}

func (s *server) DBRecoverCheck(item any) error {
	if err := s.DB.Take(item); err == nil {
		return ErrAlreadyExists
	}
	if err := s.DB.TakeUnscoped(item); err != nil {
		return ErrNotExists
	}
	return nil
}

func (s *server) IDInterfaceDeleteCheck(item idInterface) error {
	return s.IDInterfaceValidAndExistsCheck(item)
}

func (s *server) IDInterfaceUpdateCheck(item, old idInterface) error {
	if err := isValid(item); err != nil {
		return err
	}
	old.SetID(item.GetID())
	if err := s.DB.Take(old); err != nil {
		return ErrNotExists
	}
	return nil
}

func (s *server) IDInterfaceValidAndExistsCheck(item idInterface) error {
	if err := isValid(item); err != nil {
		return err
	}
	if err := s.DB.Take(item); err != nil {
		return ErrNotExists
	}
	return nil
}

func (s *server) IDInterfaceGetCheck(item idInterface) error {
	return s.IDInterfaceValidAndExistsCheck(item)
}

func (s *server) IDInterfaceModifyCheck(item idInterface) error {
	return s.IDInterfaceValidAndExistsCheck(item)
}

func (s *server) ObjectManageCheck(user User, domain Domain, item Object) error {
	if err := s.IDInterfaceModifyCheck(item); err != nil {
		return err
	}
	return s.CheckObject(user, domain, item, Manage)
}

func (s *server) ObjectParentCheck(user User, domain Domain, object Object) error {
	if object.GetParentID() == 0 {
		return ErrCantOperateRootObject
	}
	parent := newByE(object)
	parent.SetID(object.GetParentID())
	if err := s.ObjectManageCheck(user, domain, parent); err != nil {
		return err
	}
	if parent.GetObjectType() != object.GetObjectType() {
		return ErrInValidObjectType
	}
	return nil
}

func (s *server) ObjectParentToDescendantCheck(domain Domain, object Object, old Object) error {
	if object.GetParentID() == 0 || object.GetParentID() == old.GetParentID() {
		return nil
	}
	to := newByE(object)
	to.SetID(object.GetParentID())
	if ok, _ := s.Enforcer.EnforceObject(object, to, domain); ok {
		return ErrParentToDescendant
	}
	return nil
}

func (s *server) ObjectUpdateCheck(user User, domain Domain, object Object) error {
	old := newByE(object)
	if err := s.IDInterfaceUpdateCheck(object, old); err != nil {
		return err
	}
	if object.GetID() == object.GetParentID() {
		return ErrParentCanNotBeItself
	}
	if object.GetObjectType() != "" && object.GetObjectType() != old.GetObjectType() {
		return ErrCantChangeObjectType
	}
	if err := s.ObjectManageCheck(user, domain, object); err != nil {
		return err
	}
	return s.ObjectParentToDescendantCheck(domain, object, old)
}

func (s *server) RoleParentCheck(role Role) error {
	if role.GetParentID() == 0 {
		return nil
	}
	parent := newByE(role)
	parent.SetID(role.GetParentID())
	if err := s.IDInterfaceModifyCheck(parent); err != nil {
		return err
	}
	if role.GetObject().GetID() != parent.GetObject().GetID() {
		return ErrParentCanNotDiff
	}
	return nil
}

func (s *server) RoleParentToDescendantCheck(domain Domain, role Role, old Role) error {
	if role.GetParentID() == 0 || role.GetParentID() == old.GetParentID() {
		return nil
	}
	to := newByE(role)
	to.SetID(role.GetParentID())
	if ok, _ := s.Enforcer.EnforceRole(role, to, domain); ok {
		return ErrParentToDescendant
	}
	return nil
}

func (s *server) RoleUpdateCheck(user User, domain Domain, role Role) error {
	old := newByE(role)
	if err := s.IDInterfaceUpdateCheck(role, old); err != nil {
		return err
	}
	if role.GetID() == role.GetParentID() {
		return ErrParentCanNotBeItself
	}
	if err := s.ObjectDataWriteCheck(user, domain, old, ObjectTypeRole); err != nil {
		return err
	}
	if role.GetObject().GetID() != old.GetObject().GetID() {
		return s.ObjectDataWriteCheck(user, domain, role, ObjectTypeRole)
	}
	return s.RoleParentToDescendantCheck(domain, role, old)
}
