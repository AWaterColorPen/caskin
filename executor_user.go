package caskin

// UserCreate
// if there does not exist the user then create a new one
// 1. no permission checking
// 2. create a new user into metadata database
func (e *Executor) UserCreate(user User) error {
	if err := e.DBCreateCheck(user); err != nil {
		return err
	}
	return e.DB.Create(user)
}

// UserRecover
// if there exist the user but soft deleted then recover it
// 1. no permission checking
// 2. recover the soft delete one user at metadata database
func (e *Executor) UserRecover(user User) error {
	if err := e.DBRecoverCheck(user); err != nil {
		return err
	}
	return e.DB.Recover(user)
}

// UserDelete
// if there exist the user delete user
// 1. no permission checking
// 2. delete all user's g in all domain
// 3. soft delete one user in metadata database
func (e *Executor) UserDelete(user User) error {
	if err := e.IDInterfaceDeleteCheck(user); err != nil {
		return err
	}
	domains, err := e.DB.GetAllDomain()
	if err != nil {
		return err
	}
	for _, v := range domains {
		if err = e.Enforcer.RemoveUserInDomain(user, v); err != nil {
			return err
		}
	}
	return e.DB.DeleteByID(user, user.GetID())
}

// UserUpdate
// if there exist the user update user
// 1. no permission checking
// 1. just update user's properties
func (e *Executor) UserUpdate(user User) error {
	if err := e.IDInterfaceUpdateCheck(user); err != nil {
		return err
	}
	return e.DB.Update(user)
}
