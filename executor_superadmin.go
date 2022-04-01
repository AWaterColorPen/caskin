package caskin

// SuperadminUserAdd
// add the user as superadmin role in superadmin domain
// 1. no permission checking
func (e *Executor) SuperadminUserAdd(user User) error {
	return e.writeSuperadminUser(user, e.Enforcer.AddRoleForUserInDomain)
}

// SuperadminUserDelete
// delete a user from superadmin
// 1. no permission checking
func (e *Executor) SuperadminUserDelete(user User) error {
	return e.writeSuperadminUser(user, e.Enforcer.RemoveRoleForUserInDomain)
}

// SuperadminUserGet
// get all superadmin user
// 1. no permission checking
func (e *Executor) SuperadminUserGet() ([]User, error) {
	domain := GetSuperadminDomain()
	role := GetSuperadminRole()
	us := e.Enforcer.GetUsersForRoleInDomain(role, domain)
	id := ID(us)
	return e.DB.GetUserByID(id)
}

func (e *Executor) writeSuperadminUser(user User, fn func(User, Role, Domain) error) error {
	if err := e.IDInterfaceValidAndExistsCheck(user); err != nil {
		return err
	}
	domain := GetSuperadminDomain()
	role := GetSuperadminRole()
	return fn(user, role, domain)
}
