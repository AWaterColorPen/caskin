package caskin

// RoleCreate
// if there does not exist the role then create a new one
// 1. current user has role's write permission
// 2. create a new role into metadata database
// 3. set role to parent's g in the domain
func (s *server) CreateRole(user User, domain Domain, role Role) error {
	if err := s.ObjectDataCreateCheck(user, domain, role, ObjectTypeRole); err != nil {
		return err
	}
	if err := s.RoleParentCheck(role); err != nil {
		return err
	}
	role.SetDomainID(domain.GetID())
	if err := s.DB.Create(role); err != nil {
		return err
	}
	updater := s.DefaultRoleUpdater()
	return updater.Run(role, domain)
}

// RoleRecover
// if there exist the role but soft deleted then recover it
// 1. current user has role's write permission
// 2. recover the soft delete one role at metadata database
// 3. set role to parent's g in the domain
func (s *server) RecoverRole(user User, domain Domain, role Role) error {
	if err := s.ObjectDataRecoverCheck(user, domain, role); err != nil {
		return err
	}
	if err := s.RoleParentCheck(role); err != nil {
		return err
	}
	role.SetDomainID(domain.GetID())
	if err := s.DB.Recover(role); err != nil {
		return err
	}
	updater := s.DefaultRoleUpdater()
	return updater.Run(role, domain)
}

// RoleDelete
// if there exist the object
// 1. current user has role's write permission
// 1. delete role's g in the domain
// 2. delete role's p in the domain
// 3. soft delete one role in metadata database
// 4. delete all son of the role in the domain
func (s *server) DeleteRole(user User, domain Domain, role Role) error {
	if err := s.ObjectDataDeleteCheck(user, domain, role); err != nil {
		return err
	}
	if err := s.RoleParentCheck(role); err != nil {
		return err
	}
	role.SetDomainID(domain.GetID())
	deleter := NewTreeNodeDeleter(s.DefaultRoleChildrenGetFunc(), s.DefaultRoleDeleteFunc())
	return deleter.Run(role, domain)
}

// RoleUpdate
// if there exist the role
// 1. current user has role's write permission and
// 2. update role's properties
// 3. update role to parent's g in the domain
func (s *server) UpdateRole(user User, domain Domain, role Role) error {
	if err := s.RoleUpdateCheck(user, domain, role); err != nil {
		return err
	}
	if err := s.RoleParentCheck(role); err != nil {
		return err
	}
	role.SetDomainID(domain.GetID())
	if err := s.DB.Update(role); err != nil {
		return err
	}
	updater := s.DefaultRoleUpdater()
	return updater.Run(role, domain)
}

// RoleGet
// get role
// 1. current user has role's read permission
func (s *server) GetRole(user User, domain Domain) ([]Role, error) {
	roles, err := s.DB.GetRoleInDomain(domain)
	if err != nil {
		return nil, err
	}
	return Filter(s.Enforcer, user, domain, Read, roles), nil
}
