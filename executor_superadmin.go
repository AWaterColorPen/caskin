package caskin

// AddSuperadminUser
// add the user as superadmin role in superadmin domain without permission checking
func (e *executor) AddSuperadminUser(user User) error {
	return e.writeSuperadminUser(user, e.e.AddRoleForUserInDomain)
}

// DeleteSuperadminUser
// delete a user from superadmin without permission checking
func (e *executor) DeleteSuperadminUser(user User) error {
	return e.writeSuperadminUser(user, e.e.RemoveRoleForUserInDomain)
}

// GetAllSuperadminUser
// get all superadmin user without permission checking
func (e *executor) GetAllSuperadminUser() ([]User, error) {
	if !e.options.IsEnableSuperAdmin() {
		return nil, ErrSuperAdminIsNoEnabled
	}

	domain := e.options.GetSuperadminDomain()
	role := e.options.GetSuperadminRole()
	us := e.e.GetUsersForRoleInDomain(role, domain)
	id := getIDList(us)
	return e.db.GetUserByID(id)
}

func (e *executor) writeSuperadminUser(user User, fn func(User, Role, Domain) error) error {
	if !e.options.IsEnableSuperAdmin() {
		return ErrSuperAdminIsNoEnabled
	}

	if err := isValid(user); err != nil {
		return err
	}

	if err := e.db.Take(user); err != nil {
		return err
	}

	domain := e.options.GetSuperadminDomain()
	role := e.options.GetSuperadminRole()
	return fn(user, role, domain)
}
