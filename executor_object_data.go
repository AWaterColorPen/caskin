package caskin

import "github.com/ahmetb/go-linq/v3"

// FilterObjectData
// filter object_data with action
func (e *Executor) FilterObjectData(source interface{}, action Action) ([]ObjectData, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	var result []ObjectData
	linq.From(source).Where(func(v interface{}) bool {
		return Check(e.Enforcer, u, d, v.(ObjectData), action)
	}).ToSlice(&result)
	return result, nil
}

// Enforce
// check permission of object_data with action
func (e *Executor) Enforce(item ObjectData, action Action) error {
	return e.check(item, action)
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
