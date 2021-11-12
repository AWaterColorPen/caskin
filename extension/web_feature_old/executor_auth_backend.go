package web_feature_old

import (
	"github.com/awatercolorpen/caskin"
)

func (e *Executor) AuthBackendAPIEnforce(backend *Backend) error {
	object, err := e.getCacheBackendAPIObject(backend)
	if err != nil || object == nil {
		object = &caskin.SampleNoPermissionObject{}
	}
	if e.e.EnforceObject(object, caskin.Read) != nil {
		return caskin.ErrNoBackendAPIPermission
	}
	return nil
}

func (e *Executor) getCacheBackendAPIObject(backend *Backend) (caskin.Object, error) {
	key := backend.GetName()
	if u, ok := LocalCache.Get(key); ok {
		if u == nil {
			return nil, nil
		}
		return u.(caskin.Object), nil
	}

	object, err := e.takeBackend(backend)
	if err != nil {
		LocalCache.SetDefault(key, nil)
	} else {
		LocalCache.SetDefault(key, object)
	}
	return object, err
}
