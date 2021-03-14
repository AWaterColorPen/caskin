package web_feature

import "github.com/awatercolorpen/caskin"

func (e *Executor) CreateFeature(feature *Feature) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	object := e.objectFactory()
	setFeatureParentID(object)
	return e.e.CreateObjectWithCustomizedData(feature, object)
}

func (e *Executor) RecoverFeature(feature *Feature) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.RecoverObjectWithCustomizedData(feature)
}

func (e *Executor) DeleteFeature(feature *Feature, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.DeleteObjectWithCustomizedData(feature, object)
}

func (e *Executor) UpdateFeature(feature *Feature, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	setFeatureParentID(object)
	return e.e.UpdateObjectWithCustomizedData(feature, object)
}

func (e *Executor) GetFeature() ([]*caskin.CustomizedDataPair, error) {
	objects, err := e.e.GetObjects(ObjectTypeFeature)
	if err != nil {
		return nil, err
	}
	return caskin.ObjectArray2CustomizedDataPair(objects, backendFactory)
}
