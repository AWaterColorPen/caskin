package caskin

import "github.com/ahmetb/go-linq/v3"

// GetAllPoliciesForRole
// 1. get all policies which current user has role and object's read permission in current domain
// 2. get role to objects 's p as ObjectsForRole in current domain
// 3. build role's tree
func (e *executor) GetAllPoliciesForRole() ([]*PoliciesForRole, error) {
	currentUser, currentDomain, err := e.provider.Get()
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

	objects, err := e.mdb.GetObjectInDomain(currentDomain)
	if err != nil {
		return nil, err
	}

	os := e.filterWithNoError(currentUser, currentDomain, Read, objects)
	objects = []Object{}
	for _, v := range os {
		objects = append(objects, v.(Object))
	}
	om := getIDMap(objects)

	e.e.GetPoliciesInDomain(currentDomain)
	var prs []*PoliciesForRole
	for _, v := range roles {
		if p, ok := tree[v.GetID()]; ok {
			v.SetParentID(p)
		}

		pr := &PoliciesForRole{Role: v}
		policy := e.e.GetPoliciesForRoleInDomain(v, currentDomain)
		for _, p := range policy {
			p.Object.GetID()
			if object, ok := om[p.Object.GetID()]; ok {
				pr.Policies = append(pr.Policies, &Policy{
					Role:   v,
					Object: object.(Object),
					Domain: currentDomain,
					Action: p.Action,
				})
			}
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

// ModifyPoliciesForRole
// if current user has user and role and object's write permission
// 1. modify role to objects 's p in current domain
func (e *executor) ModifyPoliciesForRole(pr *PoliciesForRole) error {
	if err := isValid(pr.Role); err != nil {
		return err
	}

	if err := e.mdb.TakeRole(pr.Role); err != nil {
		return ErrNotExists
	}

	if err := e.check(Write, pr.Role); err != nil {
		return err
	}

	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return err
	}

	role := pr.Role
	policy := e.e.GetPoliciesForRoleInDomain(role, currentDomain)
	var oid, oid1, oid2 []uint64
	for _, v := range policy {
		oid1 = append(oid1, v.Object.GetID())
	}
	for _, v := range pr.Policies {
		oid2 = append(oid2, v.Object.GetID())
	}
	oid = append(oid, oid1...)
	oid = append(oid, oid2...)
	linq.From(oid).Distinct().ToSlice(&oid)
	objects, err := e.mdb.GetObjectByID(oid)
	if err != nil {
		return err
	}
	os := e.filterWithNoError(currentUser, currentDomain, Write, objects)
	objects = []Object{}
	for _, v := range os {
		objects = append(objects, v.(Object))
	}
	om := getIDMap(objects)

	// make source and target role id list
	var source, target []*Policy
	for _, v := range policy {
		if _, ok := om[v.Object.GetID()]; ok {
			source = append(source, v)
		}
	}
	for _, v := range pr.Policies {
		if _, ok := om[v.Object.GetID()]; ok {
			target = append(target, v)
		}
	}

	// get diff to add and remove
	add, remove := DiffPolicy(source, target)
	for _, v := range add {
		if err := e.e.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}
	for _, v := range remove {
		if err := e.e.RemovePolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}

	return nil
}
