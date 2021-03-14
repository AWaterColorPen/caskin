package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

func (e *Executor) CreateBackend(backend *Backend) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.CreateObjectWithCustomizedData(backend, e.objectFactory())
}

func (e *Executor) RecoverBackend(backend *Backend) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.RecoverObjectWithCustomizedData(backend)
}

func (e *Executor) DeleteBackend(backend *Backend, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.DeleteObjectWithCustomizedData(backend, object)
}

func (e *Executor) UpdateBackend(backend *Backend, object caskin.Object) error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
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
	return caskin.ObjectArray2CustomizedDataPair(objects, backendFactory)
}
