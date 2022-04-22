package caskin

// CreateRole
// if there does not exist the role then create a new one
// 1. current user has role's write permission
// 2. create a new role into metadata database
func (s *server) CreateRole(user User, domain Domain, role Role) error {
	return s.CreateObjectData(user, domain, role, ObjectTypeRole)
}

// RecoverRole
// if there exist the role but soft deleted then recover it
// 1. current user has role's write permission
// 2. recover the soft delete one role at metadata database
func (s *server) RecoverRole(user User, domain Domain, role Role) error {
	return s.RecoverObjectData(user, domain, role)
}

// DeleteRole
// if there exist the object
// 1. current user has role's write permission
// 1. delete role's g in the domain
// 2. delete role's p in the domain
// 3. soft delete one role in metadata database
func (s *server) DeleteRole(user User, domain Domain, role Role) error {
	if err := s.CheckDeleteObjectData(user, domain, role); err != nil {
		return err
	}
	role.SetDomainID(domain.GetID())
	if err := s.Enforcer.RemoveRoleInDomain(role, domain); err != nil {
		return err
	}
	return s.DB.DeleteByID(role, role.GetID())
}

// UpdateRole
// if there exist the role
// 1. current user has role's write permission and
// 2. update role's properties
func (s *server) UpdateRole(user User, domain Domain, role Role) error {
	return s.UpdateObjectData(user, domain, role, ObjectTypeRole)
}

// GetRole
// get role
// 1. current user has role's read permission
func (s *server) GetRole(user User, domain Domain) ([]Role, error) {
	roles, err := s.DB.GetRoleInDomain(domain)
	if err != nil {
		return nil, err
	}
	return Filter(s.Enforcer, user, domain, Read, roles), nil
}
