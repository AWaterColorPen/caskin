package caskin

// SuperadminAdd
// add the user as superadmin role in superadmin domain
// 1. no permission checking
func (s *server) SuperadminAdd(user User) error {
	return s.writeSuperadminUser(user, s.Enforcer.AddRoleForUserInDomain)
}

// SuperadminDelete
// delete a user from superadmin
// 1. no permission checking
func (s *server) SuperadminDelete(user User) error {
	return s.writeSuperadminUser(user, s.Enforcer.RemoveRoleForUserInDomain)
}

// SuperadminGet
// get all superadmin user
// 1. no permission checking
func (s *server) SuperadminGet() ([]User, error) {
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
