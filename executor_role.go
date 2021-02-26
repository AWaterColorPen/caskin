package caskin

import "github.com/ahmetb/go-linq/v3"

// GetAllUsersForRole
// 1. get all user which current user has read permission in current domain
// 2. get all role which current user has read permission in current domain
// 3. get role to users 's g as UsersForRole in current domain
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

	roles, err := e.mdb.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	roles = e.filterWithNoError(currentUser, currentDomain, Read, roles).([]Role)

	var urs []*UsersForRole
	for _, v := range roles {
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
