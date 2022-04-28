package caskin

import (
	"github.com/ahmetb/go-linq/v3"
)

// GetPolicy
// get all policies
// 1. current user has role and object's read permission in current domain
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
// 2. get role to object 's p as Policy in current domain
func (s *server) GetPolicyByRole(user User, domain Domain, byRole Role) ([]*Policy, error) {
	if err := s.CheckGetObjectData(user, domain, byRole); err != nil {
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

// ModifyPolicyPerRole
// if current user has role and object's write permission
// 1. modify role to object 's p in current domain
// 2. policy required object and action
func (s *server) ModifyPolicyPerRole(user User, domain Domain, perRole Role, input []*Policy) error {
	if err := s.CheckModifyObjectData(user, domain, perRole); err != nil {
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
		v.Role, v.Domain = perRole, domain
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
