package caskin

import "github.com/ahmetb/go-linq/v3"

// GetAllRolesForUser
// 1. get all user which current user has read permission in current domain
// 2. get all role which current user has read permission in current domain
// 3. get user to roles 's g as RolesForUser in current domain
func (e *executor) GetAllRolesForUser() ([]*RolesForUser, error) {
	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	us := e.e.GetUsersInDomain(currentDomain)
	uid := getIDList(us)
	linq.From(uid).Distinct().ToSlice(&uid)
	users, err := e.mdb.GetUserByID(uid)
	if err != nil {
		return nil, err
	}
	users = e.filterWithNoError(currentUser, currentDomain, Read, users).([]User)

	roles, err := e.mdb.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	roles = e.filterWithNoError(currentUser, currentDomain, Read, roles).([]Role)
	rm := getIDMap(roles)

	var rus []*RolesForUser
	for _, v := range users {
		ru := &RolesForUser{User: v}
		rs := e.e.GetRolesForUserInDomain(v, currentDomain)
		for _, r := range rs {
			if role, ok := rm[r.GetID()]; ok {
				ru.Roles = append(ru.Roles, role.(Role))
			}
		}
		rus = append(rus, ru)
	}

	return rus, nil
}

// ModifyRolesForUser if current user has user and role's write permission
// 1. modify user to roles 's g in current domain
func (e *executor) ModifyRolesForUser(ru *RolesForUser) error {
	if err := isValid(ru.User); err != nil {
		return err
	}

	if err := e.mdb.TakeUser(ru.User); err != nil {
		return ErrNotExists
	}

	if err := e.check(Write, ru.User); err != nil {
		return err
	}

	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return err
	}

	user := ru.User
	rs := e.e.GetRolesForUserInDomain(user, currentDomain)
	rid1 := getIDList(rs)
	rid2 := getIDList(ru.Roles)

	// get all role data
	var rid []uint64
	rid = append(rid, rid1...)
	rid = append(rid, rid2...)
	linq.From(rid).Distinct().ToSlice(&rid)
	roles, err := e.mdb.GetRoleByID(rid)
	if err != nil {
		return err
	}
	roles = e.filterWithNoError(currentUser, currentDomain, Write, roles).([]Role)
	rm := getIDMap(roles)

	// make source and target role id list
	var source, target []interface{}
	for _, v := range rid1 {
		if _, ok := rm[v]; ok {
			source = append(source, v)
		}
	}
	for _, v := range rid2 {
		if _, ok := rm[v]; ok {
			target = append(target, v)
		}
	}

	// get diff to add and remove
	add, remove := Diff(source, target)
	for _, v := range add {
		r := rm[v.(uint64)]
		if err := e.e.AddRoleForUserInDomain(user, r.(Role), currentDomain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		r := rm[v.(uint64)]
		if err := e.e.RemoveRoleForUserInDomain(user, r.(Role), currentDomain); err != nil {
			return err
		}
	}

	return nil
}

// CreateDomain if there does not exist the domain, then create a new one
// 1. create a new domain into metadata database
// 2. initialize the new domain
func (e *executor) CreateUser(User User) error {
	return e.createOrRecoverUser(domain, e.mdb.CreateUser)
}

// RecoverUser if there exist the domain but soft deleted, then recover it
// 1. recover the soft delete one domain at metadata database
// 2. re initialize the recovering domain
func (e *executor) RecoverUser(domain User) error {
	return e.createOrRecoverUser(domain, e.mdb.RecoverUser)
}

// DeleteUser if user has domain's write permission
// 1. delete all user's g in the domain
// 2. don't delete any role's g or object's g2 in the domain
// 3. soft delete one domain in metadata database
func (e *executor) DeleteUser(domain User) error {
	fn := func(domain User) error {
		if err := e.e.RemoveUsersInUser(domain); err != nil {
			return err
		}
		return e.mdb.DeleteUserByID(domain.GetID())
	}

	return e.writeUser(domain, fn)
}

// UpdateUser if there exist the domain and user has domain's write permission
// 1. just update domain's properties
func (e *executor) UpdateUser(domain User) error {
	return e.writeUser(domain, e.mdb.UpdateUser)
}

// ReInitializeUser if there exist the domain and user has domain's write permission
// 1. just re initialize the domain
func (e *executor) ReInitializeUser(domain User) error {
	return e.writeUser(domain, e.initializeUser)
}

// GetAllUser if user has domain's read permission
// 1. get all domain
func (e *executor) GetAllUser() ([]User, error) {
	domains, err := e.mdb.GetAllUser()
	if err != nil {
		return nil, err
	}

	out, err := e.filter(Read, domains)
	if err != nil {
		return nil, err
	}

	return out.([]User), nil
}

func (e *executor) createOrRecoverUser(domain User, fn func(User) error) error {
	if err := e.mdb.TakeUser(domain); err == nil {
		return ErrAlreadyExists
	}

	if err := fn(domain); err != nil {
		return err
	}

	return e.initializeUser(domain)
}

func (e *executor) writeUser(domain User, fn func(User) error) error {
	if err := isValid(domain); err != nil {
		return err
	}

	if err := e.mdb.TakeUser(domain); err != nil {
		return ErrNotExists
	}

	if err := e.check(Write, domain); err != nil {
		return err
	}

	return fn(domain)
}
