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
		return CheckObjectData(e.Enforcer, u, d, v.(ObjectData), action)
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
		return CheckObject(e.Enforcer, u, d, v.(Object), action)
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

// CreateObjectDataCheckPermission
// check permission of creating object_data
func (e *Executor) CreateObjectDataCheckPermission(item ObjectData, ty ObjectType) error {
	return e.ObjectDataCreateCheck(item, ty)
}

// RecoverObjectDataCheckPermission
// check permission of recover object_data
func (e *Executor) RecoverObjectDataCheckPermission(item ObjectData) error {
	return e.ObjectDataRecoverCheck(item)
}

// UpdateObjectDataCheckPermission
// check permission of updating object_data's
func (e *Executor) UpdateObjectDataCheckPermission(item ObjectData, old ObjectData, ty ObjectType) error {
	return e.ObjectDataUpdateCheck(item, old, ty)
}

// DeleteObjectDataCheckPermission
// check permission of deleting object_data
func (e *Executor) DeleteObjectDataCheckPermission(item ObjectData) error {
	return e.ObjectDataDeleteCheck(item)
}
