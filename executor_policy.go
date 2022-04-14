package caskin

import "github.com/ahmetb/go-linq/v3"

// GetPolicy
// get all policies
// 1. current user has role and object's read permission in current domain
// 2. build role's tree
func (s *server) GetPolicy(user User, domain Domain) ([]*Policy, error) {
	roles, err := s.GetRole(user, domain)
	if err != nil {
		return nil, err
	}
	objects, err := s.GetObject(user, domain, Manage)
	if err != nil {
		return nil, err
	}

	om := IDMap(objects)
	var list []*Policy
	for _, v := range roles {
		policy := s.Enforcer.GetPoliciesForRoleInDomain(v, domain)
		for _, p := range policy {
			if object, ok := om[p.Object.GetID()]; ok {
				list = append(list, &Policy{
					Role:   v,
					Object: object,
					Domain: domain,
					Action: p.Action,
				})
			}
		}
	}

	return list, nil
}

// GetPolicyByRole
// 1. get policy which current user has role and object's read permission in current domain
// 2. get user to role 's g as Policy in current domain
func (s *server) GetPolicyByRole(user User, domain Domain, byRole Role) ([]*Policy, error) {
	if err := s.ObjectDataGetCheck(user, domain, byRole); err != nil {
		return nil, err
	}
	objects, err := s.GetObject(user, domain, Manage)
	if err != nil {
		return nil, err
	}

	om := IDMap(objects)
	var list []*Policy
	policy := s.Enforcer.GetPoliciesForRoleInDomain(byRole, domain)
	for _, p := range policy {
		if object, ok := om[p.Object.GetID()]; ok {
			list = append(list, &Policy{
				Role:   byRole,
				Object: object,
				Domain: domain,
				Action: p.Action,
			})
		}
	}

	return list, nil
}

// GetPolicyByObject
// 1. get policy which current user has role and object's read permission in current domain
// 2. get user to role 's g as Policy in current domain
func (s *server) GetPolicyByObject(user User, domain Domain, byObject Object) ([]*Policy, error) {
	if err := s.ObjectManageCheck(user, domain, byObject); err != nil {
		return nil, err
	}
	roles, err := s.GetRole(user, domain)
	if err != nil {
		return nil, err
	}

	rm := IDMap(roles)
	var list []*Policy
	policy := s.Enforcer.GetPoliciesForObjectInDomain(byObject, domain)
	for _, p := range policy {
		if role, ok := rm[p.Role.GetID()]; ok {
			list = append(list, &Policy{
				Role:   role,
				Object: byObject,
				Domain: domain,
				Action: p.Action,
			})
		}
	}

	return list, nil
}

// ModifyPolicyPerRole
// if current user has role and object's write permission
// 1. modify role to object 's p in current domain
func (s *server) ModifyPolicyPerRole(user User, domain Domain, perRole Role, input []*Policy) error {
	if err := s.ObjectDataModifyCheck(user, domain, perRole); err != nil {
		return err
	}
	list := PolicyList(input)
	if err := list.IsValidWithRole(perRole); err != nil {
		return err
	}

	policy := s.Enforcer.GetPoliciesForRoleInDomain(perRole, domain)
	var oid, oid1, oid2 []uint64
	for _, v := range policy {
		oid1 = append(oid1, v.Object.GetID())
	}
	for _, v := range input {
		oid2 = append(oid2, v.Object.GetID())
	}
	oid = append(oid, oid1...)
	oid = append(oid, oid2...)
	linq.From(oid).Distinct().ToSlice(&oid)
	objects, err := s.DB.GetObjectByID(oid)
	if err != nil {
		return err
	}
	objects = Filter(s.Enforcer, user, domain, Manage, objects)
	om := IDMap(objects)

	// make source and target role id list
	var source, target []*Policy
	for _, v := range policy {
		if _, ok := om[v.Object.GetID()]; ok {
			source = append(source, v)
		}
	}
	for _, v := range input {
		if _, ok := om[v.Object.GetID()]; ok {
			target = append(target, v)
		}
	}

	// get diff to add and remove
	add, remove := DiffPolicy(source, target)
	for _, v := range add {
		if err = s.Enforcer.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}
	for _, v := range remove {
		if err = s.Enforcer.RemovePolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}

	return nil
}

// ModifyPolicyPerObject
// if current user has role and object's write permission
// 1. modify role to object 's p in current domain
func (s *server) ModifyPolicyPerObject(user User, domain Domain, perObject Object, input []*Policy) error {
	if err := s.ObjectManageCheck(user, domain, perObject); err != nil {
		return err
	}
	list := PolicyList(input)
	if err := list.IsValidWithObject(perObject); err != nil {
		return err
	}

	policy := s.Enforcer.GetPoliciesForObjectInDomain(perObject, domain)
	var rid, rid1, rid2 []uint64
	for _, v := range policy {
		rid1 = append(rid1, v.Role.GetID())
	}
	for _, v := range input {
		rid2 = append(rid2, v.Role.GetID())
	}
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
	var source, target []*Policy
	for _, v := range policy {
		if _, ok := rm[v.Role.GetID()]; ok {
			source = append(source, v)
		}
	}
	for _, v := range input {
		if _, ok := rm[v.Role.GetID()]; ok {
			target = append(target, v)
		}
	}

	// get diff to add and remove
	add, remove := DiffPolicy(source, target)
	for _, v := range add {
		if err = s.Enforcer.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}
	for _, v := range remove {
		if err = s.Enforcer.RemovePolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}
	return nil
}
