package web_feature

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/caskin"
)

// GetFeatureRelation
// 1. get all feature to backend and frontend 's relations, not inheritance relations
func (e *Executor) GetFeatureRelation() (Relations, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	dump, err := e.Dump()
	if err != nil {
		return nil, err
	}
	return caskin.SortedInheritanceRelations(dump.FeatureRelation), nil
}

// GetFeatureRelationByFeature
// 1. get one feature to backend and frontend 's relation, not inheritance relation
func (e *Executor) GetFeatureRelationByFeature(object caskin.Object) (Relation, error) {
	relation, _, err := e.featureRelationPerFeatureInternal(object)
	return relation, err
}

// ModifyFeatureRelationPerFeature
// 1. modify one feature to backend and frontend 's relation, not inheritance relation
func (e *Executor) ModifyFeatureRelationPerFeature(object caskin.Object, relation Relation) error {
	old, dump, err := e.featureRelationPerFeatureInternal(object)
	if err != nil {
		return err
	}
	relation = dump.InitSingleFeatureRelation(relation)

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

func (e *Executor) featureRelationPerFeatureInternal(object caskin.Object) (relation Relation, dump *Dump, err error) {
	if err = e.operationPermissionCheck(); err != nil {
		return
	}
	object.SetObjectType(ObjectTypeFeature)
	if err = e.e.ObjectDataModifyCheck(object); err != nil {
		return
	}
	children := e.e.Enforcer.GetChildrenForObjectInDomain(object, e.operationDomain)
	relation = caskin.InheritanceRelation{}
	for _, v := range children {
		relation = append(relation, v.GetID())
	}
	dump, err = e.Dump()
	if err != nil {
		return
	}
	relation = dump.InitSingleFeatureRelation(relation)
	relation = caskin.SortedInheritanceRelation(relation)
	return
}