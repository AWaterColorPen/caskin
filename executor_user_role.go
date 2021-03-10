package caskin

import (
	"github.com/ahmetb/go-linq/v3"
)

// GetUserRolePair
// 1. get all user
// 2. get all role which current user has read permission in current domain
// 3. get user to roles 's g as UserRolePair in current domain
func (e *executor) GetUserRolePair() ([]*UserRolePair, error) {
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

	rs := e.e.GetRolesInDomain(currentDomain)
	tree := getTree(rs)
	roles, err := e.mdb.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	r := e.filterWithNoError(currentUser, currentDomain, Read, roles)
	roles = []Role{}
	for _, v := range r {
		roles = append(roles, v.(Role))
	}
	rm := getIDMap(roles)
	for _, v := range roles {
		if p, ok := tree[v.GetID()]; ok {
			v.SetParentID(p)
		}
	}

	var list []*UserRolePair
	for _, v := range users {
		rs := e.e.GetRolesForUserInDomain(v, currentDomain)
		for _, r := range rs {
			if role, ok := rm[r.GetID()]; ok {
				list = append(list, &UserRolePair{User: v, Role: role.(Role)})
			}
		}
	}

	return list, nil
}

// GetUserRolePairByUser
// 1. get role which current user has read permission in current domain
// 2. get user to role 's g as UserRolePair in current domain
func (e *executor) GetUserRolePairByUser(user User) ([]*UserRolePair, error) {
	if err := isValid(user); err != nil {
		return nil, err
	}

	if err := e.mdb.Take(user); err != nil {
		return nil, err
	}

	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	roles, err := e.mdb.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	out := e.filterWithNoError(currentUser, currentDomain, Read, roles)
	roles = []Role{}
	for _, v := range out {
		roles = append(roles, v.(Role))
	}
	rm := getIDMap(roles)

	var list []*UserRolePair
	rs := e.e.GetRolesForUserInDomain(user, currentDomain)
	for _, r := range rs {
		if role, ok := rm[r.GetID()]; ok {
			list = append(list, &UserRolePair{User: user, Role: role.(Role)})
		}
	}

	return list, nil
}

// GetUserRolePairByRole
// 1. get role which current user has read permission in current domain
// 2. get user to role 's g as UserRolePair in current domain
func (e *executor) GetUserRolePairByRole(role Role) ([]*UserRolePair, error) {
	if err := e.getOrModifyCheck(role, Read); err != nil {
		return nil, err
	}

	_, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	us := e.e.GetUsersForRoleInDomain(role, currentDomain)
	oid := getIDList(us)
	users, err := e.mdb.GetUserByID(oid)
	if err != nil {
		return nil, err
	}

	var list []*UserRolePair
	for _, v := range users {
		list = append(list, &UserRolePair{User: v, Role: role})
	}

	return list, nil
}

// ModifyUserRolePairPerUser
// if current user has role's write permission
// 1. modify user to roles 's g in current domain
func (e *executor) ModifyUserRolePairPerUser(user User, input []*UserRolePair) error {
	if err := isValid(user); err != nil {
		return err
	}

	pairs := UserRolePairs(input)
	if err := pairs.IsValidWithUser(user); err != nil {
		return err
	}

	if err := e.mdb.Take(user); err != nil {
		return ErrNotExists
	}

	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return err
	}

	rs := e.e.GetRolesForUserInDomain(user, currentDomain)
	rid1 := getIDList(rs)
	rid2 := pairs.RoleID()

	// get all role data
	var rid []uint64
	rid = append(rid, rid1...)
	rid = append(rid, rid2...)
	linq.From(rid).Distinct().ToSlice(&rid)
	roles, err := e.mdb.GetRoleByID(rid)
	if err != nil {
		return err
	}
	r := e.filterWithNoError(currentUser, currentDomain, Write, roles)
	roles = []Role{}
	for _, v := range r {
		roles = append(roles, v.(Role))

	}
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

// ModifyUserRolePairPerRole
// if current user has role's write permission
// 1. modify role's to user 's g in current domain
func (e *executor) ModifyUserRolePairPerRole(role Role, input []*UserRolePair) error {
	if err := e.getOrModifyCheck(role, Write); err != nil {
		return err
	}

	pairs := UserRolePairs(input)
	if err := pairs.IsValidWithRole(role); err != nil {
		return err
	}

	_, currentDomain, err := e.provider.Get()
	if err != nil {
		return err
	}

	us := e.e.GetUsersForRoleInDomain(role, currentDomain)
	uid1 := getIDList(us)
	uid2 := pairs.UserID()

	// get all role data
	var uid []uint64
	uid = append(uid, uid1...)
	uid = append(uid, uid2...)
	linq.From(uid).Distinct().ToSlice(&uid)
	users, err := e.mdb.GetUserByID(uid)
	if err != nil {
		return err
	}
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
		if err := e.e.RemoveRoleForUserInDomain(u.(User), role, currentDomain); err != nil {
			return err
		}
	}

	return nil
}
