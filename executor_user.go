package caskin

// UserCreate
// if there does not exist the user then create a new one
// 1. no permission checking
// 2. create a new user into metadata database
func (s *server) UserCreate(user User) error {
	if err := s.DBCreateCheck(user); err != nil {
		return err
	}
	return s.DB.Create(user)
}

// UserRecover
// if there exist the user but soft deleted then recover it
// 1. no permission checking
// 2. recover the soft delete one user at metadata database
func (s *server) UserRecover(user User) error {
	if err := s.DBRecoverCheck(user); err != nil {
		return err
	}
	return s.DB.Recover(user)
}

// UserDelete
// if there exist the user delete user
// 1. no permission checking
// 2. delete all user's g in all domain
// 3. soft delete one user in metadata database
func (s *server) UserDelete(user User) error {
	if err := s.IDInterfaceDeleteCheck(user); err != nil {
		return err
	}
	domains, err := s.DB.GetAllDomain()
	if err != nil {
		return err
	}
	for _, v := range domains {
		if err = s.Enforcer.RemoveUserInDomain(user, v); err != nil {
			return err
		}
	}
	return s.DB.DeleteByID(user, user.GetID())
}

// UserUpdate
// if there exist the user update user
// 1. no permission checking
// 1. just update user's properties
func (s *server) UserUpdate(user User) error {
	if err := s.IDInterfaceUpdateCheck(user, newByE(user)); err != nil {
		return err
	}
	return s.DB.Update(user)
}
