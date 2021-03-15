package caskin

// AddSuperadminUser
// add the user as superadmin role in superadmin domain without permission checking
func (e *Executor) AddSuperadminUser(user User) error {
	return e.writeSuperadminUser(user, e.e.AddRoleForUserInDomain)
}

// DeleteSuperadminUser
// delete a user from superadmin without permission checking
func (e *Executor) DeleteSuperadminUser(user User) error {
	return e.writeSuperadminUser(user, e.e.RemoveRoleForUserInDomain)
}

// GetAllSuperadminUser
// get all superadmin user without permission checking
func (e *Executor) GetAllSuperadminUser() ([]User, error) {
	if !e.options.IsEnableSuperAdmin() {
		return nil, ErrSuperAdminIsNoEnabled
	}

	domain := e.options.GetSuperadminDomain()
	role := e.options.GetSuperadminRole()
	us := e.e.GetUsersForRoleInDomain(role, domain)
	id := getIDList(us)
	return e.DB.GetUserByID(id)
}

func (e *Executor) writeSuperadminUser(user User, fn func(User, Role, Domain) error) error {
	if !e.options.IsEnableSuperAdmin() {
		return ErrSuperAdminIsNoEnabled
	}
	if err := e.IDInterfaceValidAndExistsCheck(user); err != nil {
		return err
	}
	domain := e.options.GetSuperadminDomain()
	role := e.options.GetSuperadminRole()
	return fn(user, role, domain)
}
