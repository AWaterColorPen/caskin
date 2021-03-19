package caskin

import (
	"github.com/ahmetb/go-linq/v3"
)

// GetUserRolePair
// 1. get all user
// 2. get all role which current user has read permission in current domain
// 3. get user to roles 's g as UserRolePair in current domain
func (e *Executor) GetUserRolePair() ([]*UserRolePair, error) {
	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	us := e.Enforcer.GetUsersInDomain(currentDomain)
	uid := Users(us).ID()
	linq.From(uid).Distinct().ToSlice(&uid)
	users, err := e.DB.GetUserByID(uid)
	if err != nil {
		return nil, err
	}

	// rs := e.Enforcer.GetRolesInDomain(currentDomain)
	// tree := getTree(rs)
	roles, err := e.DB.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	rs := e.filterWithNoError(currentUser, currentDomain, Read, roles)
	linq.From(rs).ToSlice(&roles)
	rm := Roles(roles).IDMap()
	// for _, v := range roles {
	// 	if p, ok := tree[v.GetID()]; ok {
	// 		v.SetParentID(p)
	// 	}
	// }

	var list []*UserRolePair
	for _, v := range users {
		rs := e.Enforcer.GetRolesForUserInDomain(v, currentDomain)
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
func (e *Executor) GetUserRolePairByUser(user User) ([]*UserRolePair, error) {
	if err := e.IDInterfaceGetCheck(user); err != nil {
		return nil, err
	}

	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	roles, err := e.DB.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	out := e.filterWithNoError(currentUser, currentDomain, Read, roles)
	linq.From(out).ToSlice(&roles)
	rm := Roles(roles).IDMap()

	var list []*UserRolePair
	rs := e.Enforcer.GetRolesForUserInDomain(user, currentDomain)
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
func (e *Executor) GetUserRolePairByRole(role Role) ([]*UserRolePair, error) {
	if err := e.ObjectDataGetCheck(role); err != nil {
		return nil, err
	}

	_, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	us := e.Enforcer.GetUsersForRoleInDomain(role, currentDomain)
	uid := Users(us).ID()
	users, err := e.DB.GetUserByID(uid)
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
func (e *Executor) ModifyUserRolePairPerUser(user User, input []*UserRolePair) error {
	if err := isValid(user); err != nil {
		return err
	}

	pairs := UserRolePairs(input)
	if err := pairs.IsValidWithUser(user); err != nil {
		return err
	}

	if err := e.DB.Take(user); err != nil {
		return ErrNotExists
	}

	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return err
	}

	rs := e.Enforcer.GetRolesForUserInDomain(user, currentDomain)
	rid1 := Roles(rs).ID()
	rid2 := pairs.RoleID()

	// get all role data
	var rid []uint64
	rid = append(rid, rid1...)
	rid = append(rid, rid2...)
	linq.From(rid).Distinct().ToSlice(&rid)
	roles, err := e.DB.GetRoleByID(rid)
	if err != nil {
		return err
	}
	out := e.filterWithNoError(currentUser, currentDomain, Write, roles)
	linq.From(out).ToSlice(&roles)
	rm := Roles(roles).IDMap()

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
		if err := e.Enforcer.AddRoleForUserInDomain(user, r.(Role), currentDomain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		r := rm[v.(uint64)]
		if err := e.Enforcer.RemoveRoleForUserInDomain(user, r.(Role), currentDomain); err != nil {
			return err
		}
	}

	return nil
}

// ModifyUserRolePairPerRole
// if current user has role's write permission
// 1. modify role's to user 's g in current domain
func (e *Executor) ModifyUserRolePairPerRole(role Role, input []*UserRolePair) error {
	if err := e.ObjectDataModifyCheck(role); err != nil {
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

	us := e.Enforcer.GetUsersForRoleInDomain(role, currentDomain)
	uid1 := Users(us).ID()
	uid2 := pairs.UserID()

	// get all role data
	var uid []uint64
	uid = append(uid, uid1...)
	uid = append(uid, uid2...)
	linq.From(uid).Distinct().ToSlice(&uid)
	users, err := e.DB.GetUserByID(uid)
	if err != nil {
		return err
	}
	um := Users(users).IDMap()

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
		if err := e.Enforcer.AddRoleForUserInDomain(u.(User), role, currentDomain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		u := um[v.(uint64)]
		if err := e.Enforcer.RemoveRoleForUserInDomain(u.(User), role, currentDomain); err != nil {
			return err
		}
	}

	return nil
}
