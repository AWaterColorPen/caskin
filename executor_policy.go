package caskin

import "github.com/ahmetb/go-linq/v3"

// PolicyGet
// get all policies
// 1. current user has role and object's read permission in current domain
// 2. build role's tree
func (e *baseService) PolicyGet(user User, domain Domain) ([]*Policy, error) {
	roles, err := e.RoleGet(user, domain)
	if err != nil {
		return nil, err
	}
	objects, err := e.ObjectGet(user, domain, Manage)
	if err != nil {
		return nil, err
	}

	om := IDMap(objects)
	var list []*Policy
	for _, v := range roles {
		policy := e.Enforcer.GetPoliciesForRoleInDomain(v, domain)
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

// PolicyByRoleGet
// 1. get policy which current user has role and object's read permission in current domain
// 2. get user to role 's g as Policy in current domain
func (e *baseService) PolicyByRoleGet(user User, domain Domain, byRole Role) ([]*Policy, error) {
	if err := e.ObjectDataGetCheck(user, domain, byRole); err != nil {
		return nil, err
	}
	objects, err := e.ObjectGet(user, domain, Manage)
	if err != nil {
		return nil, err
	}

	om := IDMap(objects)
	var list []*Policy
	policy := e.Enforcer.GetPoliciesForRoleInDomain(byRole, domain)
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

// PolicyByObjectGet
// 1. get policy which current user has role and object's read permission in current domain
// 2. get user to role 's g as Policy in current domain
func (e *baseService) PolicyByObjectGet(user User, domain Domain, byObject Object) ([]*Policy, error) {
	if err := e.ObjectManageCheck(user, domain, byObject); err != nil {
		return nil, err
	}
	roles, err := e.RoleGet(user, domain)
	if err != nil {
		return nil, err
	}

	rm := IDMap(roles)
	var list []*Policy
	policy := e.Enforcer.GetPoliciesForObjectInDomain(byObject, domain)
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

// PolicyPerRoleModify
// if current user has role and object's write permission
// 1. modify role to object 's p in current domain
func (e *baseService) PolicyPerRoleModify(user User, domain Domain, perRole Role, input []*Policy) error {
	if err := e.ObjectDataModifyCheck(user, domain, perRole); err != nil {
		return err
	}
	list := PolicyList(input)
	if err := list.IsValidWithRole(perRole); err != nil {
		return err
	}

	policy := e.Enforcer.GetPoliciesForRoleInDomain(perRole, domain)
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
	objects, err := e.DB.GetObjectByID(oid)
	if err != nil {
		return err
	}
	objects = Filter(e.Enforcer, user, domain, Manage, objects)
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
		if err = e.Enforcer.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}
	for _, v := range remove {
		if err = e.Enforcer.RemovePolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}

	return nil
}

// PolicyPerObjectModify
// if current user has role and object's write permission
// 1. modify role to object 's p in current domain
func (e *baseService) PolicyPerObjectModify(user User, domain Domain, perObject Object, input []*Policy) error {
	if err := e.ObjectManageCheck(user, domain, perObject); err != nil {
		return err
	}
	list := PolicyList(input)
	if err := list.IsValidWithObject(perObject); err != nil {
		return err
	}

	policy := e.Enforcer.GetPoliciesForObjectInDomain(perObject, domain)
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
	roles, err := e.DB.GetRoleByID(rid)
	if err != nil {
		return err
	}
	roles = Filter(e.Enforcer, user, domain, Write, roles)
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
		if err = e.Enforcer.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}
	for _, v := range remove {
		if err = e.Enforcer.RemovePolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}
	return nil
}
