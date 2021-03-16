package web_feature

import "github.com/awatercolorpen/caskin"

func (e *Executor) CreateFrontend(frontend *Frontend) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.CreateObjectWithCustomizedData(frontend, e.objectFactory())
}

func (e *Executor) RecoverFrontend(frontend *Frontend) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.RecoverObjectWithCustomizedData(frontend)
}

func (e *Executor) DeleteFrontend(frontend *Frontend, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.DeleteObjectWithCustomizedData(frontend, object)
}

func (e *Executor) UpdateFrontend(frontend *Frontend, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.UpdateObjectWithCustomizedData(frontend, object)
}

func (e *Executor) GetFrontend() ([]*caskin.CustomizedDataPair, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	objects, err := e.e.GetObjects(ObjectTypeFrontend)
	if err != nil {
		return nil, err
	}
	return caskin.ObjectArray2CustomizedDataPair(objects, FrontendFactory)
}
