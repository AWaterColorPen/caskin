package caskin

import "github.com/ahmetb/go-linq/v3"

// CreateUser
// if there does not exist the user
// then create a new one without permission checking
// 1. create a new user into metadata database
func (e *Executor) CreateUser(user User) error {
	if err := e.DBCreateCheck(user); err != nil {
		return err
	}
	return e.DB.Create(user)
}

// RecoverUser
// if there exist the user but soft deleted
// then recover it without permission checking
// 1. recover the soft delete one user at metadata database
func (e *Executor) RecoverUser(user User) error {
	if err := e.DBRecoverCheck(user); err != nil {
		return err
	}
	return e.DB.Recover(user)
}

// DeleteUser
// if there exist the user
// delete user without permission checking
// 1. delete all user's g in all domain
// 2. soft delete one user in metadata database
func (e *Executor) DeleteUser(user User) error {
	if err := e.IDInterfaceDeleteCheck(user); err != nil {
		return err
	}
	domains, err := e.DB.GetAllDomain()
	if err != nil {
		return err
	}
	for _, v := range domains {
		if err := e.Enforcer.RemoveUserInDomain(user, v); err != nil {
			return err
		}
	}
	return e.DB.DeleteByID(user, user.GetID())
}

// UpdateUser
// if there exist the user
// update user without permission checking
// 1. just update user's properties
func (e *Executor) UpdateUser(user User) error {
	tmp := e.factory.NewUser()
	if err := e.IDInterfaceUpdateCheck(user, tmp); err != nil {
		return err
	}
	return e.DB.Update(user)
}

func (e *Executor) GetUsers() (Users, error) {
	_, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	us := e.Enforcer.GetUsersInDomain(currentDomain)
	uid := getIDList(us)
	linq.From(uid).Distinct().ToSlice(&uid)
	users, err := e.DB.GetUserByID(uid)
	if err != nil {
		return nil, err
	}
	return users, nil

}
