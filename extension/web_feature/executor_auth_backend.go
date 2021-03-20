package web_feature

import (
	"time"

	"github.com/awatercolorpen/caskin"
)

func (e *Executor) AuthBackendAPIEnforce(backend *Backend) error {
	object, err := e.getBackendAPIObject(backend)
	if err != nil || object == nil{
		object = &caskin.SampleNoPermissionObject{}
	}
	if e.check(object) != nil {
		return caskin.ErrNoBackendAPIPermission
	}
	return nil
}

func (e *Executor) getBackendAPIObject(backend *Backend) (caskin.Object, error) {
	if e.enableBackendAPIAuthCache {
		return e.getCacheBackendAPIObject(backend)
	}
	return e.getSyncBackendAPIObject(backend)
}

func (e *Executor) getCacheBackendAPIObject(backend *Backend) (caskin.Object, error) {
	key := backend.GetName()
	if u, ok := LocalCache.Get(key); ok {
		return u.(caskin.Object), nil
	} else {
		object, err := e.getSyncBackendAPIObject(backend)
		if err != nil {
			LocalCache.Set(key, nil, 30*time.Second)
		} else {
			LocalCache.SetDefault(key, object)
		}
		return object, err
	}
}

func (e *Executor) getSyncBackendAPIObject(backend *Backend) (caskin.Object, error) {
	return e.takeBackend(backend)
}
