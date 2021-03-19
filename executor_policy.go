package caskin

import "github.com/ahmetb/go-linq/v3"

// GetPolicyList
// 1. get all policies which current user has role and object's read permission in current domain
// 2. build role's tree
func (e *Executor) GetPolicyList() ([]*Policy, error) {
	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	roles, err := e.DB.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	out1 := e.filterWithNoError(currentUser, currentDomain, Read, roles)
	linq.From(out1).ToSlice(&roles)

	objects, err := e.DB.GetObjectInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	out2 := e.filterWithNoError(currentUser, currentDomain, Read, objects)
	linq.From(out2).ToSlice(&objects)
	om := Objects(objects).IDMap()

	var list []*Policy
	for _, v := range roles {
		policy := e.Enforcer.GetPoliciesForRoleInDomain(v, currentDomain)
		for _, p := range policy {
			if object, ok := om[p.Object.GetID()]; ok {
				list = append(list, &Policy{
					Role:   v,
					Object: object.(Object),
					Domain: currentDomain,
					Action: p.Action,
				})
			}
		}
	}

	return list, nil
}

// GetPolicyListByRole
// 1. get policy which current user has role and object's read permission in current domain
// 2. get user to role 's g as Policy in current domain
func (e *Executor) GetPolicyListByRole(role Role) ([]*Policy, error) {
	if err := e.ObjectDataGetCheck(role); err != nil {
		return nil, err
	}

	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	objects, err := e.DB.GetObjectInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	out := e.filterWithNoError(currentUser, currentDomain, Read, objects)
	linq.From(out).ToSlice(&objects)
	om := Objects(objects).IDMap()

	var list []*Policy
	policy := e.Enforcer.GetPoliciesForRoleInDomain(role, currentDomain)
	for _, p := range policy {
		if object, ok := om[p.Object.GetID()]; ok {
			list = append(list, &Policy{
				Role:   role,
				Object: object.(Object),
				Domain: currentDomain,
				Action: p.Action,
			})
		}
	}

	return list, nil
}

// GetPolicyListByObject
// 1. get policy which current user has role and object's read permission in current domain
// 2. get user to role 's g as Policy in current domain
func (e *Executor) GetPolicyListByObject(object Object) ([]*Policy, error) {
	if err := e.ObjectDataGetCheck(object); err != nil {
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

	var list []*Policy
	policy := e.Enforcer.GetPoliciesForObjectInDomain(object, currentDomain)
	for _, p := range policy {
		if role, ok := rm[p.Role.GetID()]; ok {
			list = append(list, &Policy{
				Role:   role.(Role),
				Object: object,
				Domain: currentDomain,
				Action: p.Action,
			})
		}
	}

	return list, nil
}

// ModifyPolicyListPerRole
// if current user has role and object's write permission
// 1. modify role to objects 's p in current domain
func (e *Executor) ModifyPolicyListPerRole(role Role, input []*Policy) error {
	if err := e.ObjectDataModifyCheck(role); err != nil {
		return err
	}

	list := PolicyList(input)
	if err := list.IsValidWithRole(role); err != nil {
		return err
	}

	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return err
	}

	policy := e.Enforcer.GetPoliciesForRoleInDomain(role, currentDomain)
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

	out := e.filterWithNoError(currentUser, currentDomain, Write, objects)
	linq.From(out).ToSlice(&objects)
	om := Objects(objects).IDMap()

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
		if err := e.Enforcer.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}
	for _, v := range remove {
		if err := e.Enforcer.RemovePolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}

	return nil
}

// ModifyPolicyListPerObject
// if current user has role and object's write permission
// 1. modify role to objects 's p in current domain
func (e *Executor) ModifyPolicyListPerObject(object Object, input []*Policy) error {
	if err := e.ObjectDataModifyCheck(object); err != nil {
		return err
	}

	list := PolicyList(input)
	if err := list.IsValidWithObject(object); err != nil {
		return err
	}

	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return err
	}

	policy := e.Enforcer.GetPoliciesForObjectInDomain(object, currentDomain)
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

	rs := e.filterWithNoError(currentUser, currentDomain, Write, roles)
	roles = []Role{}
	for _, v := range rs {
		roles = append(roles, v.(Role))
	}
	rm := Roles(roles).IDMap()

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
		if err := e.Enforcer.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}
	for _, v := range remove {
		if err := e.Enforcer.RemovePolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}

	return nil
}
