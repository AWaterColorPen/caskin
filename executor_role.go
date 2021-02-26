package caskin

import "github.com/ahmetb/go-linq/v3"

// GetAllUsersForRole
// 1. get all user which current user has read permission in current domain
// 2. get all role which current user has read permission in current domain
// 3. get role to users 's g as UsersForRole in current domain
// 4. build role's tree
func (e *executor) GetAllUsersForRole() ([]*UsersForRole, error) {
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
	um := getIDMap(users)

	rs := e.e.GetRolesInDomain(currentDomain)
	tree := getTree(rs)
	roles, err := e.mdb.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	roles = e.filterWithNoError(currentUser, currentDomain, Read, roles).([]Role)

	var urs []*UsersForRole
	for _, v := range roles {
		if p, ok := tree[v.GetID()]; ok {
			v.SetParentID(p)
		}

		ur := &UsersForRole{Role: v}
		uus := e.e.GetUsersForRoleInDomain(v, currentDomain)
		for _, u := range uus {
			if user, ok := um[u.GetID()]; ok {
				ur.Users = append(ur.Users, user.(User))
			}
		}
		urs = append(urs, ur)
	}

	return urs, nil
}

// ModifyUsersForRole if current user has user and role's write permission
// 1. modify role to users 's g in current domain
func (e *executor) ModifyUsersForRole(ur *UsersForRole) error {
	if err := isValid(ur.Role); err != nil {
		return err
	}

	if err := e.mdb.TakeRole(ur.Role); err != nil {
		return ErrNotExists
	}

	if err := e.check(Write, ur.Role); err != nil {
		return err
	}

	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return err
	}

	role := ur.Role
	us := e.e.GetUsersForRoleInDomain(role, currentDomain)
	uid1 := getIDList(us)
	uid2 := getIDList(ur.Users)

	// get all role data
	var uid []uint64
	uid = append(uid, uid1...)
	uid = append(uid, uid2...)
	linq.From(uid).Distinct().ToSlice(&uid)
	users, err := e.mdb.GetUserByID(uid)
	if err != nil {
		return err
	}
	users = e.filterWithNoError(currentUser, currentDomain, Write, users).([]User)
	um := getIDMap(users)

	// make source and target role id list
	var source, target []interface{}
	for _, v := range uid1 {
		if _, ok := um[v]; ok {
			source = append(source, v)
		}
	}
	for _, v := range uid2 {
		if _, ok := um[v]; ok {
			target = append(target, v)
		}
	}

	// get diff to add and remove
	add, remove := Diff(source, target)
	for _, v := range add {
		u := um[v.(uint64)]
		if err := e.e.AddRoleForUserInDomain(u.(User), role, currentDomain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		u := um[v.(uint64)]
		if err := e.e.RemoveRoleForUserInDomain(u.(User), role,  currentDomain); err != nil {
			return err
		}
	}

	return nil
}

// CreateRole if there does not exist the role, then create a new one
// 1. create a new role into metadata database
func (e *executor) CreateRole(role Role) error {
	return e.createOrRecoverRole(role, e.mdb.CreateRole)
}

// RecoverRole if there exist the role but soft deleted, then recover it
// 1. recover the soft delete one role at metadata database
func (e *executor) RecoverRole(role Role) error {
	return e.createOrRecoverRole(role, e.mdb.RecoverRole)
}

// DeleteRole if current user has role's write permission
// 1. delete all user's g in the domain
// 2. don't delete any role's g or object's g2 in the domain
// 3. soft delete one domain in metadata database
func (e *executor) DeleteRole(domain Domain) error {
	fn := func(domain Domain) error {
		if err := e.e.RemoveUsersInDomain(domain); err != nil {
			return err
		}
		return e.mdb.DeleteDomainByID(domain.GetID())
	}

	return e.writeDomain(domain, fn)
}

// UpdateRole if there exist the domain and user has domain's write permission
// 1. just update domain's properties
func (e *executor) UpdateRole(domain Domain) error {
	return e.writeDomain(domain, e.mdb.UpdateDomain)
}

func (e *executor) createOrRecoverRole(role Role, fn func(Role) error) error {
	if err := e.mdb.TakeRole(role); err == nil {
		return ErrAlreadyExists
	}

	take := func(id uint64, domain Domain) (parentEntry, error) {
		r := e.factory.NewRole()
		r.SetDomainID(domain.GetID())
		err := e.mdb.TakeRole(r)
		return r, err
	}

	if err := e.checkParentEntryWrite(role, take); err != nil {
		return err
	}

	return fn(role)
}
