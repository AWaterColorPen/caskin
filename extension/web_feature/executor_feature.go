package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

func (e *Executor) CreateFeature(feature *Feature, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	if err := isEmptyObject(object); err != nil {
		return err
	}
	e.setFeatureRoot(object)
	return e.e.CreateObjectWithCustomizedData(feature, object)
}

func (e *Executor) RecoverFeature(feature *Feature, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	if err := isEmptyObject(object); err != nil {
		return err
	}
	return e.e.RecoverObjectWithCustomizedData(feature, object)
}

func (e *Executor) DeleteFeature(object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}

	object.SetObjectType((&Feature{}).GetObjectType())
	if err := e.e.ObjectDataDeleteCheck(object); err != nil {
		return err
	}

	if err := e.e.ObjectTreeNodeParentCheck(object); err != nil {
		return err
	}

	provider := e.e.GetCurrentProvider()
	_, domain, _ := provider.Get()
	object.SetDomainID(domain.GetID())
	deleter := caskin.NewTreeNodeEntryDeleter(e.featureChildrenGetFunc(), e.e.DefaultObjectDeleteFunc())

	return deleter.Run(object, domain)
}

func (e *Executor) UpdateFeature(feature *Feature, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	e.setFeatureRoot(object)
	return e.e.UpdateObjectWithCustomizedData(feature, object)
}

func (e *Executor) GetFeature() ([]*caskin.CustomizedDataPair, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	objects, err := e.e.GetObjects(ObjectTypeFeature)
	if err != nil {
		return nil, err
	}
	return caskin.ObjectArray2CustomizedDataPair(objects, FeatureFactory)
}

func (e *Executor) featureChildrenGetFunc() caskin.TreeNodeEntryChildrenGetFunc {
	feature, _ := e.GetFeature()
	index := initTreeMapFromPair(feature)
	return func(p caskin.TreeNodeEntry, domain caskin.Domain) []caskin.TreeNodeEntry {
		var out []caskin.TreeNodeEntry
		in := e.e.Enforcer.GetChildrenForObjectInDomain(p.(caskin.Object), domain)
		for _, v := range in {
			if _, ok := index[v.GetID()]; ok {
				out = append(out, v)
			}
		}
		return out
	}
}
