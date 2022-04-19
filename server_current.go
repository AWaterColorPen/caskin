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
