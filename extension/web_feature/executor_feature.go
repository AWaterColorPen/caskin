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
	setFeatureRoot(object)
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
	return e.e.DeleteObjectWithCustomizedData(&Feature{}, object)
}

func (e *Executor) UpdateFeature(feature *Feature, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	setFeatureRoot(object)
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
