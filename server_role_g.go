package caskin

// AddRoleG
// from role >= to role
// 1. current user has from/to role's modify permission
// 2. add from to g in the domain
func (s *server) AddRoleG(user User, domain Domain, from, to Role) error {
	if err := s.CheckModifyObjectData(user, domain, from); err != nil {
		return err
	}
	if err := s.CheckModifyObjectData(user, domain, to); err != nil {
		return err
	}
	return s.Enforcer.AddParentForRoleInDomain(from, to, domain)
}

// RemoveRoleG
// from role >= to role
// 1. current user has from/to role's modify permission
// 2. remove from to g in the domain
func (s *server) RemoveRoleG(user User, domain Domain, from, to Role) error {
	if err := s.CheckModifyObjectData(user, domain, from); err != nil {
		return err
	}
	if err := s.CheckModifyObjectData(user, domain, to); err != nil {
		return err
	}
	return s.Enforcer.RemoveParentForRoleInDomain(from, to, domain)
}
