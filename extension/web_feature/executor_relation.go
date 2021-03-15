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
	children := e.e.Enforcer.GetChildrenForObjectInDomain(object, e.operationDomain)
	old := caskin.InheritanceRelation{}
	for _, v := range children {
		relation = append(old, v.GetID())
	}

	add, remove := caskin.Diff(old, relation)
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

func (e *Executor) SyncRelationsToAllDomain() error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	relations, err := e.GetRelation()
	if err != nil {
		return err
	}
	domains, err := e.e.GetAllDomain()
	if err != nil {
		return err
	}
	for _, v := range domains {
		if err := e.syncRelationsToOneDomain(relations, v); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) SyncRelationsToOneDomain(domain caskin.Domain) error {
	relations, err := e.GetRelation()
	if err != nil {
		return err
	}
	return e.syncRelationsToOneDomain(relations, domain)
}

func (e *Executor) syncRelationsToOneDomain(relations caskin.InheritanceRelations, domain caskin.Domain) error {
	version := ""
	if versioned, ok := domain.(VersionedDomain); ok && versioned.GetVersion() == version {
		return nil
	}



	if versioned, ok := domain.(VersionedDomain); ok {
		versioned.SetVersion(version)
		return e.e.DB.Update(domain)
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
