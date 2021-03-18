package web_feature

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/caskin"
)

// GetFeatureRelation
// 1. get all feature to backend and frontend 's relations, not inheritance relations
func (e *Executor) GetFeatureRelation() (FeatureRelations, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	all := e.allWebFeatureRelation(e.operationDomain)
	pair, err := e.GetFeature()
	if err != nil {
		return nil, err
	}
	relations := filterInheritancesToFeatureRelations(all, pair)
	relations = caskin.SortedInheritanceRelations(relations)
	return relations, nil
}

// GetFeatureRelationByFeature
// 1. get one feature to backend and frontend 's relation, not inheritance relation
func (e *Executor) GetFeatureRelationByFeature(object caskin.Object) (FeatureRelation, error) {
	relation, _, err := e.featureRelationInternalHelpFunc(object)
	return relation, err
}

// ModifyFeatureRelationPerFeature
// 1. modify one feature to backend and frontend 's relation, not inheritance relation
func (e *Executor) ModifyFeatureRelationPerFeature(object caskin.Object, relation FeatureRelation) error {
	old, pair, err := e.featureRelationInternalHelpFunc(object)
	if err != nil {
		return err
	}
	// TODO avoid id not in feature backend frontend
	relation = filterInheritanceToFeatureRelation(relation, pair)

	var source, target []interface{}
	linq.From(old).ToSlice(&source)
	linq.From(relation).ToSlice(&target)
	add, remove := caskin.Diff(source, target)
	for _, v := range add {
		o := e.objectFactory()
		o.SetID(v.(uint64))
		if err := e.e.Enforcer.AddParentForObjectInDomain(o, object, e.operationDomain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		o := e.objectFactory()
		o.SetID(v.(uint64))
		if err := e.e.Enforcer.RemoveParentForObjectInDomain(o, object, e.operationDomain); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) featureRelationInternalHelpFunc(object caskin.Object) (FeatureRelation, []*caskin.CustomizedDataPair, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, nil, err
	}
	object.SetObjectType(ObjectTypeFeature)
	if err := e.e.ObjectDataModifyCheck(object); err != nil {
		return nil, nil, err
	}
	children := e.e.Enforcer.GetChildrenForObjectInDomain(object, e.operationDomain)
	relation := caskin.InheritanceRelation{}
	for _, v := range children {
		relation = append(relation, v.GetID())
	}
	pair, err := e.GetFeature()
	if err != nil {
		return nil, nil, err
	}
	relation = filterInheritanceToFeatureRelation(relation, pair)
	relation = caskin.SortedInheritanceRelation(relation)
	return relation, pair, nil
}

func (e *Executor) allWebFeatureRelation(domain caskin.Domain) caskin.InheritanceRelations {
	queue := []caskin.Object{GetFeatureRootObject(), GetBackendRootObject(), GetFrontendRootObject()}
	inQueue := map[uint64]bool{}
	for _, v := range queue {
		inQueue[v.GetID()] = true
	}

	m := caskin.InheritanceRelations{}
	for i := 0; i < len(queue); i++ {
		m[queue[i].GetID()] = caskin.InheritanceRelation{}
		ll := e.e.Enforcer.GetChildrenForObjectInDomain(queue[i], domain)
		for _, v := range ll {
			if _, ok := inQueue[v.GetID()]; !ok {
				queue = append(queue, v)
				inQueue[v.GetID()] = true
			}
			m[queue[i].GetID()] = append(m[queue[i].GetID()], v.GetID())
		}
	}

	return m
}

func filterInheritancesToFeatureRelations(relations caskin.InheritanceRelations, pair []*caskin.CustomizedDataPair) FeatureRelations {
	om := map[interface{}]bool{}
	for _, v := range pair {
		if v.Object.GetName() == GetFeatureRootObject().GetName() {
			continue
		}
		om[v.Object.GetID()] = true
	}
	m := FeatureRelations{}
	for k, v := range relations {
		if _, ok := om[k]; ok {
			m[k] = filterInheritanceToFeatureRelation(v, pair)
		}
	}
	return m
}

func filterInheritanceToFeatureRelation(relation caskin.InheritanceRelation, pair []*caskin.CustomizedDataPair) FeatureRelation {
	om := map[interface{}]bool{
		GetFeatureRootObject().GetID():  true,
		GetFrontendRootObject().GetID(): true,
		GetBackendRootObject().GetID():  true,
	}
	for _, v := range pair {
		om[v.Object.GetID()] = true
	}
	m := FeatureRelation{}
	for _, v := range relation {
		if _, ok := om[v]; !ok {
			m = append(m, v)
		}
	}
	return m
}
