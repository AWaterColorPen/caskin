package caskin

// AddSuperadminUser
// if user is superadmin
// 1. add the user as superadmin role in superadmin domain
func (e *executor) AddSuperadminUser(user User) error {
	return e.writeSuperadminUser(user, e.e.AddRoleForUserInDomain)
}

// DeleteSuperadminUser
// if user is superadmin
// 1. delete the user from superadmin role in superadmin domain
func (e *executor) DeleteSuperadminUser(user User) error {
	return e.writeSuperadminUser(user, e.e.RemoveRoleForUserInDomain)
}

// GetAllSuperadminUser
// if user is superadmin
// 1. get all superadmin user
func (e *executor) GetAllSuperadminUser() ([]User, error) {
	if !e.option.IsEnableSuperAdmin() {
		return nil, ErrSuperAdminIsNoEnabled
	}

	domain := e.option.GetSuperAdminDomain()
	role := e.option.GetSuperAdminRole()
	us := e.e.GetUsersForRoleInDomain(role, domain)
	id := getIDList(us)
	return e.mdb.GetUserByID(id)
}

func (e *executor) writeSuperadminUser(user User, fn func(User, Role, Domain) error) error {
	if !e.option.IsEnableSuperAdmin() {
		return ErrSuperAdminIsNoEnabled
	}

	if err := isValid(user); err != nil {
		return err
	}

	if err := e.mdb.TakeUser(user); err != nil {
		return err
	}

	domain := e.option.GetSuperAdminDomain()
	role := e.option.GetSuperAdminRole()
	return fn(user, role, domain)
}
