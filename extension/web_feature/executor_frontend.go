package web_feature

import "github.com/awatercolorpen/caskin"

func (e *Executor) CreateFrontend(frontend *Frontend, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	if err := isEmptyObject(object); err != nil {
		return err
	}
	e.setFrontendRoot(object)
	return e.e.CreateObjectWithCustomizedData(frontend, object)
}

func (e *Executor) RecoverFrontend(frontend *Frontend, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	if err := isEmptyObject(object); err != nil {
		return err
	}
	return e.e.RecoverObjectWithCustomizedData(frontend, object)
}

func (e *Executor) DeleteFrontend(object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.DeleteObjectWithCustomizedData(&Frontend{}, object)
}

func (e *Executor) UpdateFrontend(frontend *Frontend, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	e.setFrontendRoot(object)
	caskin.CustomizedData2Object(frontend, object)
	return e.frontendAndBackendUpdate(object)
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

func (e *Executor) getFrontend() ([]*caskin.CustomizedDataPair, error) {
	objects, err := e.e.DB.GetObjectInDomain(e.operationDomain, ObjectTypeFrontend)
	if err != nil {
		return nil, err
	}
	return caskin.ObjectArray2CustomizedDataPair(objects, FrontendFactory)
}
