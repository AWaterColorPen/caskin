package caskin

func (s *server) ObjectDataWriteCheck(user User, domain Domain, item ObjectData, ty ObjectType) error {
	if err := s.CheckObjectData(user, domain, item, Write); err != nil {
		return err
	}
	o := item.GetObject()
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
	if item.GetObject().GetID() != old.GetObject().GetID() {
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
