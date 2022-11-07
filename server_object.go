package caskin

// CreateObject
// if there does not exist the object then create a new one
// 1. current user has manage permission of object's parent
// 2. create a new object into metadata database
// 3. set object to parent's g2 in the domain
func (s *server) CreateObject(user User, domain Domain, object Object) error {
	if err := s.DBCreateCheck(object); err != nil {
		return err
	}
	if err := s.ObjectParentCheck(user, domain, object); err != nil {
		return err
	}
	if err := s.ObjectHierarchyCheck(domain, object); err != nil {
		return err
	}
	object.SetDomainID(domain.GetID())
	if err := s.DB.Create(object); err != nil {
		return err
	}
	updater := defaultObjectUpdater(s.Enforcer)
	return updater.Run(object, domain)
}

// RecoverObject
// if there exist the object but soft deleted then recover it
// 1. current user has manage permission of object's parent
// 2. recover the soft delete one object at metadata database
// 3. set object to parent's g2 in the domain
func (s *server) RecoverObject(user User, domain Domain, object Object) error {
	if err := s.DBRecoverCheck(object); err != nil {
		return err
	}
	if err := s.ObjectParentCheck(user, domain, object); err != nil {
		return err
	}
	object.SetDomainID(domain.GetID())
	if err := s.DB.Recover(object); err != nil {
		return err
	}
	updater := defaultObjectUpdater(s.Enforcer)
	return updater.Run(object, domain)
}

// DeleteObject
// if there exist the object
// 1. current user has manage permission of object's parent
// 2. delete object's g2 in the domain
// 3. delete object's p in the domain
// 4. soft delete one object in metadata database
// 5. delete all son of the object in the domain
func (s *server) DeleteObject(user User, domain Domain, object Object) error {
	if err := s.IDInterfaceDeleteCheck(object); err != nil {
		return err
	}
	if err := s.ObjectParentCheck(user, domain, object); err != nil {
		return err
	}
	object.SetDomainID(domain.GetID())
	deleter := defaultObjectDeleter(s.Enforcer, s.DB)
	return deleter.Run(object, domain)
}

// UpdateObject
// if there exist the object
// 1. current user has manage permission of object's parent to change parent_id
//    current user has manage permission of object to change properties
// 2. update object's properties
// 3. update object to parent's g2 in the domain
func (s *server) UpdateObject(user User, domain Domain, object Object) error {
	if err := s.ObjectUpdateCheck(user, domain, object); err != nil {
		return err
	}
	if err := s.ObjectParentCheck(user, domain, object); err != nil {
		return err
	}
	if err := s.ObjectHierarchyCheck(domain, object); err != nil {
		return err
	}
	object.SetDomainID(domain.GetID())
	if err := s.DB.Update(object); err != nil {
		return err
	}
	updater := defaultObjectUpdater(s.Enforcer)
	return updater.Run(object, domain)
}

// GetObject
// get choose object
// 1. current user has permission of object
//    manage permission for admin to manage
//    read/write permission for customer to get directory
// 2. get object by type
func (s *server) GetObject(user User, domain Domain, action Action, ty ...ObjectType) ([]Object, error) {
	objects, err := s.DB.GetObjectInDomain(domain, ty...)
	if err != nil {
		return nil, err
	}
	return Filter(s.Enforcer, user, domain, action, objects), nil
}

// GetObjectHierarchyLevel
// get input object's hierarchy_level
// 1. current user has manage permission of object
func (s *server) GetObjectHierarchyLevel(user User, domain Domain, object Object) (int, error) {
	if err := s.ObjectManageCheck(user, domain, object); err != nil {
		return 0, err
	}
	hierarchyLevel := objectHierarchyBFS(domain, object, s.Enforcer.GetParentsForObjectInDomain)
	return hierarchyLevel, nil
}

func defaultObjectUpdater(e IEnforcer) *objectUpdater {
	return NewObjectUpdater(e.GetParentsForObjectInDomain, e.AddParentForObjectInDomain, e.RemoveParentForObjectInDomain)
}

func defaultObjectDeleter(e IEnforcer, db MetaDB) *objectDeleter {
	fn1 := func(p Object, d Domain) error {
		if err := e.RemoveObjectInDomain(p, d); err != nil {
			return err
		}
		return db.DeleteByID(p, p.GetID())
	}
	fn2 := func(p Object, domain Domain) []Object {
		os := e.GetChildrenForObjectInDomain(p, domain)
		om := IDMap(os)
		os2, _ := db.GetObjectInDomain(domain, p.GetObjectType())
		for _, v := range os2 {
			if v.GetParentID() != p.GetID() {
				continue
			}
			if _, ok := om[v.GetID()]; !ok {
				om[v.GetID()] = v
				os = append(os, v)
			}
		}
		return os
	}
	return NewObjectDeleter(fn2, fn1)
}
