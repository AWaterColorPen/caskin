package caskin

// SetCurrent
// set current user and current domain for ICurrentService
func (s *server) SetCurrent(user User, domain Domain) IService {
	return &server{
		Enforcer:      s.Enforcer,
		DB:            s.DB,
		Dictionary:    s.Dictionary,
		CurrentUser:   user,
		CurrentDomain: domain,
	}
}

// CreateObjectDataWithCurrent
// for ICurrentService
func (s *server) CreateObjectDataWithCurrent(item ObjectData, ty ObjectType) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.CreateObjectData(s.CurrentUser, s.CurrentDomain, item, ty)
}

// RecoverObjectDataWithCurrent
// for ICurrentService
func (s *server) RecoverObjectDataWithCurrent(item ObjectData) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.RecoverObjectData(s.CurrentUser, s.CurrentDomain, item)
}

// DeleteObjectDataWithCurrent
// for ICurrentService
func (s *server) DeleteObjectDataWithCurrent(item ObjectData) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.DeleteObjectData(s.CurrentUser, s.CurrentDomain, item)
}

// UpdateObjectDataWithCurrent
// for ICurrentService
func (s *server) UpdateObjectDataWithCurrent(item ObjectData, ty ObjectType) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.UpdateObjectData(s.CurrentUser, s.CurrentDomain, item, ty)
}

// CheckCreateObjectDataWithCurrent
// for ICurrentService
func (s *server) CheckCreateObjectDataWithCurrent(item ObjectData, ty ObjectType) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.CheckCreateObjectData(s.CurrentUser, s.CurrentDomain, item, ty)
}

// CheckRecoverObjectDataWithCurrent
// for ICurrentService
func (s *server) CheckRecoverObjectDataWithCurrent(item ObjectData) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.CheckRecoverObjectData(s.CurrentUser, s.CurrentDomain, item)
}

// CheckDeleteObjectDataWithCurrent
// for ICurrentService
func (s *server) CheckDeleteObjectDataWithCurrent(item ObjectData) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.CheckDeleteObjectData(s.CurrentUser, s.CurrentDomain, item)
}

// CheckWriteObjectDataWithCurrent
// for ICurrentService
func (s *server) CheckWriteObjectDataWithCurrent(item ObjectData, ty ObjectType) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.CheckWriteObjectData(s.CurrentUser, s.CurrentDomain, item, ty)
}

// CheckUpdateObjectDataWithCurrent
// for ICurrentService
func (s *server) CheckUpdateObjectDataWithCurrent(item ObjectData, ty ObjectType) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.CheckUpdateObjectData(s.CurrentUser, s.CurrentDomain, item, ty)
}

// CheckModifyObjectDataWithCurrent
// for ICurrentService
func (s *server) CheckModifyObjectDataWithCurrent(item ObjectData) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.CheckModifyObjectData(s.CurrentUser, s.CurrentDomain, item)
}

// CheckGetObjectDataWithCurrent
// for ICurrentService
func (s *server) CheckGetObjectDataWithCurrent(item ObjectData) error {
	if s.CurrentUser == nil || s.CurrentDomain == nil {
		return ErrInValidCurrent
	}
	return s.CheckGetObjectData(s.CurrentUser, s.CurrentDomain, item)
}
