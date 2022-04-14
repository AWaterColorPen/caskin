package caskin

import (
	"github.com/ahmetb/go-linq/v3"
)

// GetUserRole
// 1. get all user
// 2. get all role which current user has read permission in current domain
// 3. get user to role 's g as UserRolePair in current domain
func (s *server) GetUserRole(user User, domain Domain) ([]*UserRolePair, error) {
	users, err := s.GetUserByDomain(domain)
	if err != nil {
		return nil, err
	}
	roles, err := s.GetRole(user, domain)
	if err != nil {
		return nil, err
	}

	rm := IDMap(roles)
	var list []*UserRolePair
	for _, v := range users {
		rs := s.Enforcer.GetRolesForUserInDomain(v, domain)
		for _, r := range rs {
			if role, ok := rm[r.GetID()]; ok {
				list = append(list, &UserRolePair{User: v, Role: role})
			}
		}
	}
	return list, nil
}

// GetUserRoleByUser
// 1. get role which current user has read permission in current domain
// 2. get user to role 's g as UserRolePair in current domain
func (s *server) GetUserRoleByUser(user User, domain Domain, byUser User) ([]*UserRolePair, error) {
	if err := s.IDInterfaceGetCheck(byUser); err != nil {
		return nil, err
	}
	roles, err := s.GetRole(user, domain)
	if err != nil {
		return nil, err
	}

	rm := IDMap(roles)
	var list []*UserRolePair
	rs := s.Enforcer.GetRolesForUserInDomain(byUser, domain)
	for _, r := range rs {
		if role, ok := rm[r.GetID()]; ok {
			list = append(list, &UserRolePair{User: byUser, Role: role})
		}
	}
	return list, nil
}

// GetUserRoleByRole
// 1. get role which current user has read permission in current domain
// 2. get user to role 's g as UserRolePair in current domain
func (s *server) GetUserRoleByRole(user User, domain Domain, byRole Role) ([]*UserRolePair, error) {
	if err := s.ObjectDataGetCheck(user, domain, byRole); err != nil {
		return nil, err
	}

	us := s.Enforcer.GetUsersForRoleInDomain(byRole, domain)
	uid := ID(us)
	users, err := s.DB.GetUserByID(uid)
	if err != nil {
		return nil, err
	}

	var list []*UserRolePair
	for _, v := range users {
		list = append(list, &UserRolePair{User: v, Role: byRole})
	}
	return list, nil
}

// ModifyUserRolePerUser
// if current user has role's write permission
// 1. modify user to role 's g in current domain
func (s *server) ModifyUserRolePerUser(user User, domain Domain, perUser User, input []*UserRolePair) error {
	if err := isValid(perUser); err != nil {
		return err
	}
	pairs := UserRolePairs(input)
	if err := pairs.IsValidWithUser(perUser); err != nil {
		return err
	}
	if err := s.DB.Take(perUser); err != nil {
		return ErrNotExists
	}

	rs := s.Enforcer.GetRolesForUserInDomain(perUser, domain)
	rid1 := ID(rs)
	rid2 := pairs.RoleID()

	// get all role data
	var rid []uint64
	rid = append(rid, rid1...)
	rid = append(rid, rid2...)
	linq.From(rid).Distinct().ToSlice(&rid)
	roles, err := s.DB.GetRoleByID(rid)
	if err != nil {
		return err
	}
	roles = Filter(s.Enforcer, user, domain, Write, roles)
	rm := IDMap(roles)

	// make source and target role id list
	var source, target []uint64
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
		r := rm[v]
		if err = s.Enforcer.AddRoleForUserInDomain(user, r, domain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		r := rm[v]
		if err = s.Enforcer.RemoveRoleForUserInDomain(user, r, domain); err != nil {
			return err
		}
	}

	return nil
}

// ModifyUserRolePerRole
// if current user has role's write permission
// 1. modify role's to user 's g in current domain
func (s *server) ModifyUserRolePerRole(user User, domain Domain, perRole Role, input []*UserRolePair) error {
	if err := s.ObjectDataModifyCheck(user, domain, perRole); err != nil {
		return err
	}
	pairs := UserRolePairs(input)
	if err := pairs.IsValidWithRole(perRole); err != nil {
		return err
	}

	us := s.Enforcer.GetUsersForRoleInDomain(perRole, domain)
	uid1 := ID(us)
	uid2 := pairs.UserID()

	// get all role data
	var uid []uint64
	uid = append(uid, uid1...)
	uid = append(uid, uid2...)
	linq.From(uid).Distinct().ToSlice(&uid)
	users, err := s.DB.GetUserByID(uid)
	if err != nil {
		return err
	}
	um := IDMap(users)

	// make source and target role id list
	var source, target []uint64
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
		u := um[v]
		if err = s.Enforcer.AddRoleForUserInDomain(u, perRole, domain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		u := um[v]
		if err = s.Enforcer.RemoveRoleForUserInDomain(u, perRole, domain); err != nil {
			return err
		}
	}

	return nil
}
