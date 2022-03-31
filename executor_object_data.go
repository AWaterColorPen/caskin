package caskin

import "github.com/ahmetb/go-linq/v3"

// FilterObjectData
// filter object_data with action
func (e *Executor) FilterObjectData(source any, action Action) ([]ObjectData, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	var result []ObjectData
	linq.From(source).Where(func(v any) bool {
		return ObjectDataCheck(e.Enforcer, u, d, v.(ObjectData), action)
	}).ToSlice(&result)
	return result, nil
}

// FilterObject
// filter object with action
func (e *Executor) FilterObject(source any, action Action) ([]Object, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	var result []Object
	linq.From(source).Where(func(v any) bool {
		return ObjectCheck(e.Enforcer, u, d, v.(Object), action)
	}).ToSlice(&result)
	return result, nil
}

// EnforceObjectData
// check permission of object_data with action
func (e *Executor) EnforceObjectData(item ObjectData, action Action) error {
	return e.checkObjectData(item, action)
}

// EnforceObject
// check permission of object with action
func (e *Executor) EnforceObject(item Object, action Action) error {
	return e.checkObject(item, action)
}
