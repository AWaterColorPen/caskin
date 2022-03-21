package web_feature

import "github.com/awatercolorpen/caskin"

func (e *Executor) GetBackend() ([]*Backend, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Executor) takeBackend(*Backend) (caskin.Object, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Executor) GetFrontend() ([]*Frontend, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *Executor) GetFeature() ([]*Feature, error) {
	if err := e.operationPermissionCheck(); err != nil {
		return nil, err
	}
	return nil, nil
}
