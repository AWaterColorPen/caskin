package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

func (e *Executor) CreateBackend(backend *Backend, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	if err := isEmptyObject(object); err != nil {
		return err
	}
	setBackendRoot(object)
	return e.e.CreateObjectWithCustomizedData(backend, object)
}

func (e *Executor) RecoverBackend(backend *Backend) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.RecoverObjectWithCustomizedData(backend)
}

func (e *Executor) DeleteBackend(object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.DeleteObjectWithCustomizedData(&Backend{}, object)
}

func (e *Executor) UpdateBackend(backend *Backend, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	setBackendRoot(object)
	return e.e.UpdateObjectWithCustomizedData(backend, object)
}

func (e *Executor) GetBackend() ([]*caskin.CustomizedDataPair, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	objects, err := e.e.GetObjects(ObjectTypeBackend)
	if err != nil {
		return nil, err
	}
	return caskin.ObjectArray2CustomizedDataPair(objects, BackendFactory)
}

func (e *Executor) takeBackend(backend *Backend) (caskin.Object, error) {
	object := e.objectFactory()
	caskin.CustomizedData2Object(backend, object)
	return object, e.e.DB.Take(object)
}
