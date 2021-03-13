package caskin

// CreateUser
// if there does not exist the user
// then create a new one without permission checking
// 1. create a new user into metadata database
func (e *Executor) CreateUser(user User) error {
	return e.createOrRecoverUser(user, e.db.Create)
}

// RecoverUser
// if there exist the user but soft deleted
// then recover it without permission checking
// 1. recover the soft delete one user at metadata database
func (e *Executor) RecoverUser(user User) error {
	return e.createOrRecoverUser(user, e.db.Recover)
}

// DeleteUser
// if there exist the user
// delete user without permission checking
// 1. delete all user's g in all domain
// 2. soft delete one user in metadata database
func (e *Executor) DeleteUser(user User) error {
	fn := func(interface{}) error {
		domains, err := e.db.GetAllDomain()
		if err != nil {
			return err
		}
		for _, v := range domains {
			if err := e.e.RemoveUserInDomain(user, v); err != nil {
				return err
			}
		}
		return e.db.DeleteByID(user, user.GetID())
	}

	return e.writeUser(user, fn)
}

// UpdateUser
// if there exist the user
// update user without permission checking
// 1. just update user's properties
func (e *Executor) UpdateUser(user User) error {
	return e.writeUser(user, e.db.Update)
}

func (e *Executor) createOrRecoverUser(user User, fn func(interface{}) error) error {
	if err := e.db.Take(user); err == nil {
		return ErrAlreadyExists
	}

	return fn(user)
}

func (e *Executor) writeUser(user User, fn func(interface{}) error) error {
	if err := isValid(user); err != nil {
		return err
	}

	tmp := e.factory.NewUser()
	tmp.SetID(user.GetID())

	if err := e.db.Take(tmp); err != nil {
		return ErrNotExists
	}

	return fn(user)
}
