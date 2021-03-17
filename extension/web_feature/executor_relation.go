package web_feature

import "github.com/awatercolorpen/caskin"

func (e *Executor) GetRelation() (caskin.InheritanceRelations, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	relations := e.e.Enforcer.GetObjectInheritanceRelationInDomain(e.operationDomain)
	pair, err := e.GetFeature()
	if err != nil {
		return nil, err
	}
	relations = e.filterInheritanceRelations(relations, pair)
	relations = caskin.SortedInheritanceRelations(relations)
	return relations, nil
}

func (e *Executor) GetRelationByFeature(feature *Feature, object caskin.Object) (caskin.InheritanceRelation, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	if !caskin.CustomizedDataEqualObject(feature, object) {
		return nil, caskin.ErrCustomizedDataIsNotBelongToObject
	}
	if err := e.e.ObjectDataModifyCheck(object); err != nil {
		return nil, err
	}
	children := e.e.Enforcer.GetChildrenForObjectInDomain(object, e.operationDomain)
	relation := caskin.InheritanceRelation{}
	for _, v := range children {
		relation = append(relation, v.GetID())
	}
	return relation, nil
}

func (e *Executor) ModifyRelationPerFeature(feature *Feature, object caskin.Object, relation caskin.InheritanceRelation) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	if !caskin.CustomizedDataEqualObject(feature, object) {
		return caskin.ErrCustomizedDataIsNotBelongToObject
	}
	if err := e.e.ObjectDataModifyCheck(object); err != nil {
		return err
	}
	children := e.e.Enforcer.GetChildrenForObjectInDomain(object, e.operationDomain)
	old := caskin.InheritanceRelation{}
	for _, v := range children {
		old = append(old, v.GetID())
	}

	add, remove := caskin.Diff(old, relation)
	for _, v := range add {
		o := e.objectFactory()
		o.SetID(toUint64(v))
		if err := e.e.Enforcer.AddParentForObjectInDomain(o, object, e.operationDomain); err != nil {
			return err
		}
	}

	for _, v := range remove {
		o := e.objectFactory()
		o.SetID(toUint64(v))
		if err := e.e.Enforcer.RemoveParentForObjectInDomain(o, object, e.operationDomain); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) filterInheritanceRelations(relations caskin.InheritanceRelations, pair []*caskin.CustomizedDataPair) caskin.InheritanceRelations {
	om := map[interface{}]bool{}
	for _, v := range pair {
		om[v.Object.GetID()] = true
	}
	m := caskin.InheritanceRelations{}
	for k, v := range relations {
		if _, ok := om[k]; ok {
			m[k] = v
		}
	}
	return m
}

func toUint64(v interface{}) uint64 {
	switch u := v.(type) {
	case uint64:
		return u
	case int:
		return uint64(u)
	case int64:
		return uint64(u)
	case uint32:
		return uint64(u)
	default:
		return 0
	}
}