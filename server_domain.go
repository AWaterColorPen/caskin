package caskin

// CreateDomain
// if there does not exist the domain then create a new one
// 1. no permission checking
// 2. create a new domain into metadata database
func (s *server) CreateDomain(domain Domain) error {
	if err := s.DBCreateCheck(domain); err != nil {
		return err
	}
	return s.DB.Create(domain)
}

// RecoverDomain
// if there exist the domain but soft deleted then recover it
// 1. no permission checking
// 2. recover the soft delete one domain at metadata database
func (s *server) RecoverDomain(domain Domain) error {
	if err := s.DBRecoverCheck(domain); err != nil {
		return err
	}
	return s.DB.Recover(domain)
}

// DeleteDomain
// if there exist the domain soft delete the domain
// 1. no permission checking
// 2. delete all user's g in the domain
// 3. don't delete any role's g or object's g2 in the domain
// 4. soft delete one domain in metadata database
func (s *server) DeleteDomain(domain Domain) error {
	if err := s.IDInterfaceDeleteCheck(domain); err != nil {
		return err
	}
	if err := s.Enforcer.RemoveUsersInDomain(domain); err != nil {
		return err
	}
	return s.DB.DeleteByID(domain, domain.GetID())
}

// UpdateDomain
// if there exist the domain update domain
// 1. no permission checking
// 2. just update domain's properties
func (s *server) UpdateDomain(domain Domain) error {
	old := newByE(domain)
	if err := s.IDInterfaceUpdateCheck(domain, old); err != nil {
		return err
	}
	return s.DB.Update(domain)
}

// GetDomain
// get all domain
// 1. no permission checking
func (s *server) GetDomain() ([]Domain, error) {
	return s.DB.GetAllDomain()
}
