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
	if err := s.ObjectManageCheck(user, domain, old); err != nil {
		return err
	}
	return s.ObjectParentToDescendantCheck(domain, object, old)
}

func (s *server) ObjectHierarchyCheck(domain Domain, object Object) error {
	parent := newByE(object)
	parent.SetID(object.GetParentID())
	hierarchyLevel1 := objectHierarchyBFS(domain, parent, s.Enforcer.GetParentsForObjectInDomain)
	hierarchyLevel2 := objectHierarchyBFS(domain, object, s.Enforcer.GetChildrenForObjectInDomain)
	if hierarchyLevel := hierarchyLevel1 + hierarchyLevel2; hierarchyLevel > 10 {
		return fmt.Errorf("max directory depth %v is too large than 10", hierarchyLevel)
	}
	return nil
}

func objectHierarchyBFS(domain Domain, object Object, fn func(Object, Domain) []Object) int {
	hierarchyLevel := 0
	m := map[uint64]int{object.GetID(): 1}
	queue := []uint64{object.GetID()}
	for i := 0; i < len(queue); i++ {
		node := newByE(object)
		node.SetID(queue[i])
		for _, v := range fn(node, domain) {
			if _, ok := m[v.GetID()]; ok {
				continue
			}
			m[v.GetID()] = m[node.GetID()] + 1
			queue = append(queue, v.GetID())
		}
		if hierarchyLevel < m[node.GetID()] {
			hierarchyLevel = m[node.GetID()]
		}
	}
	return hierarchyLevel
}
