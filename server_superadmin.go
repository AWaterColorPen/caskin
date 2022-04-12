package caskin

// AddSuperadmin
// add the user as superadmin role in superadmin domain
// 1. no permission checking
func (s *server) AddSuperadmin(user User) error {
	return s.writeSuperadminUser(user, s.Enforcer.AddRoleForUserInDomain)
}

// DeleteSuperadmin
// delete a user from superadmin
// 1. no permission checking
func (s *server) DeleteSuperadmin(user User) error {
	return s.writeSuperadminUser(user, s.Enforcer.RemoveRoleForUserInDomain)
}

// GetSuperadmin
// get all superadmin user
// 1. no permission checking
func (s *server) GetSuperadmin() ([]User, error) {
	us := s.Enforcer.GetUsersForRoleInDomain(GetSuperadminRole(), GetSuperadminDomain())
	id := ID(us)
	return s.DB.GetUserByID(id)
}

func (s *server) writeSuperadminUser(user User, fn func(User, Role, Domain) error) error {
	if err := s.IDInterfaceValidAndExistsCheck(user); err != nil {
		return err
	}
	return fn(user, GetSuperadminRole(), GetSuperadminDomain())
}
