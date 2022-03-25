package web_feature

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/caskin"
)

func (e *Executor) NormalDomainGetFeatureObject() ([]caskin.Object, error) {
	objects, err := e.e.DB.GetObjectInDomain(e.operationDomain, ObjectTypeFeature)
	if err != nil {
		return nil, err
	}
	return e.e.FilterObject(objects, caskin.Read)
}

func (e *Executor) NormalDomainGetPolicyList() ([]*caskin.Policy, error) {
	provider := e.e.GetCurrentProvider()
	_, currentDomain, err := provider.Get()
	if err != nil {
		return nil, err
	}

	roles, err := e.e.DB.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	out1, err := e.e.FilterObjectData(roles, caskin.Read)
	if err != nil {
		return nil, err
	}
	linq.From(out1).ToSlice(&roles)

	objects, err := e.e.DB.GetObjectInDomain(e.operationDomain, ObjectTypeFeature)
	if err != nil {
		return nil, err
	}
	objects, err = e.e.FilterObject(objects, caskin.Read)
	if err != nil {
		return nil, err
	}
	om := caskin.IDMap(objects)

	var list []*caskin.Policy
	for _, v := range roles {
		policy := e.e.Enforcer.GetPoliciesForRoleInDomain(v, currentDomain)
		for _, p := range policy {
			if object, ok := om[p.Object.GetID()]; ok {
				list = append(list, &caskin.Policy{
					Role:   v,
					Object: object,
					Domain: currentDomain,
					Action: p.Action,
				})
			}
		}
	}

	return list, nil
}

func (e *Executor) NormalDomainGetPolicyListByRole(role caskin.Role) ([]*caskin.Policy, error) {
	if err := e.e.ObjectDataGetCheck(role); err != nil {
		return nil, err
	}

	provider := e.e.GetCurrentProvider()
	_, currentDomain, err := provider.Get()
	if err != nil {
		return nil, err
	}

	objects, err := e.e.DB.GetObjectInDomain(e.operationDomain, ObjectTypeFeature)
	if err != nil {
		return nil, err
	}
	objects, err = e.e.FilterObject(objects, caskin.Read)
	if err != nil {
		return nil, err
	}
	om := caskin.IDMap(objects)

	var list []*caskin.Policy
	policy := e.e.Enforcer.GetPoliciesForRoleInDomain(role, currentDomain)
	for _, p := range policy {
		if object, ok := om[p.Object.GetID()]; ok {
			list = append(list, &caskin.Policy{
				Role:   role,
				Object: object,
				Domain: currentDomain,
				Action: p.Action,
			})
		}
	}

	return list, nil
}

func (e *Executor) NormalDomainModifyPolicyListPerRole(role caskin.Role, input []*caskin.Policy) error {
	if err := e.e.ObjectDataModifyCheck(role); err != nil {
		return err
	}

	list := caskin.PolicyList(input)
	if err := list.IsValidWithRole(role); err != nil {
		return err
	}

	provider := e.e.GetCurrentProvider()
	_, currentDomain, err := provider.Get()
	if err != nil {
		return err
	}

	policy := e.e.Enforcer.GetPoliciesForRoleInDomain(role, currentDomain)
	objects, err := e.e.DB.GetObjectInDomain(e.operationDomain, ObjectTypeFeature)
	if err != nil {
		return err
	}
	objects, err = e.e.FilterObject(objects, caskin.Read)
	if err != nil {
		return err
	}
	om := caskin.IDMap(objects)

	// make source and target role id list
	var source, target []*caskin.Policy
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

	if err := isPolicyListValidAction(target); err != nil {
		return err
	}

	// get diff to add and remove
	add, remove := caskin.DiffPolicy(source, target)
	for _, v := range add {
		if err := e.e.Enforcer.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}
	for _, v := range remove {
		if err := e.e.Enforcer.RemovePolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}

	return nil
}
