package caskin

// CreateObjectData
// if there does not exist the item then create a new one
// 1. current user has item's write permission
// 2. create a new item into database
func (s *server) CreateObjectData(user User, domain Domain, item ObjectData, ty ObjectType) error {
	if err := s.ObjectDataCreateCheck(user, domain, item, ty); err != nil {
		return err
	}
	item.SetDomainID(domain.GetID())
	return s.DB.Create(item)
}

// RecoverObjectData
// if there exist the item but soft deleted then recover it
// 1. current user has item's write permission
// 2. recover the soft delete one item at database
func (s *server) RecoverObjectData(user User, domain Domain, item ObjectData) error {
	if err := s.ObjectDataRecoverCheck(user, domain, item); err != nil {
		return err
	}
	return s.DB.Recover(item)
}

// DeleteObjectData
// if there exist the item
// 1. current user has item's write permission
// 3. soft delete one item in database
func (s *server) DeleteObjectData(user User, domain Domain, item ObjectData) error {
	if err := s.ObjectDataDeleteCheck(user, domain, item); err != nil {
		return err
	}
	item.SetDomainID(domain.GetID())
	return s.DB.DeleteByID(item, item.GetID())
}

// UpdateObjectData
// if there exist the item
// 1. current user has item's write permission and
// 2. update item's properties
func (s *server) UpdateObjectData(user User, domain Domain, item ObjectData, ty ObjectType) error {
	if err := s.ObjectDataUpdateCheck(user, domain, item, ty); err != nil {
		return err
	}
	item.SetDomainID(domain.GetID())
	return s.DB.Update(item)
}

// GetObjectData
// get items
// 1. current user has item's read permission
// func (s *server) GetObjectData(user User, domain Domain, item ObjectData) ([]ObjectData, error) {
//
// 	GetInDomain[item](s.DB, domain)
// 	roles, err := s.DB.GetRoleInDomain(domain)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return Filter(s.Enforcer, user, domain, Read, roles), nil
// }

func (s *server) ObjectDataWriteCheck(user User, domain Domain, item ObjectData, ty ObjectType) error {
	if err := s.CheckObjectData(user, domain, item, Write); err != nil {
		return err
	}
	o := DefaultFactory().NewObject()
	o.SetID(item.GetObjectID())
	if err := s.DB.Take(o); err != nil {
		return ErrInValidObject
	}
	if o.GetObjectType() != ty {
		return ErrInValidObjectType
	}
	return nil
}

func (s *server) ObjectDataCreateCheck(user User, domain Domain, item ObjectData, ty ObjectType) error {
	if err := s.DBCreateCheck(item); err != nil {
		return err
	}
	return s.ObjectDataWriteCheck(user, domain, item, ty)
}

func (s *server) ObjectDataRecoverCheck(user User, domain Domain, item ObjectData) error {
	if err := s.DBRecoverCheck(item); err != nil {
		return err
	}
	return s.CheckObjectData(user, domain, item, Write)
}

func (s *server) ObjectDataDeleteCheck(user User, domain Domain, item ObjectData) error {
	if err := s.IDInterfaceDeleteCheck(item); err != nil {
		return err
	}
	return s.CheckObjectData(user, domain, item, Write)
}

func (s *server) ObjectDataUpdateCheck(user User, domain Domain, item ObjectData, ty ObjectType) error {
	old := newByE(item)
	if err := s.IDInterfaceUpdateCheck(item, old); err != nil {
		return err
	}
	if err := s.ObjectDataWriteCheck(user, domain, old, ty); err != nil {
		return err
	}
	if item.GetObjectID() != old.GetObjectID() {
		return s.ObjectDataWriteCheck(user, domain, item, ty)
	}
	return nil
}

func (s *server) ObjectDataModifyCheck(user User, domain Domain, item ObjectData) error {
	if err := s.IDInterfaceModifyCheck(item); err != nil {
		return err
	}
	return s.CheckObjectData(user, domain, item, Write)
}

func (s *server) ObjectDataGetCheck(user User, domain Domain, item ObjectData) error {
	if err := s.IDInterfaceGetCheck(item); err != nil {
		return err
	}
	return s.CheckObjectData(user, domain, item, Read)
}
