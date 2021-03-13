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
		return Check(e.e, u, d, v.(ObjectData), action)
	}).ToSlice(&result)
	return result, nil
}

// Enforce
// check permission of object_data with action
func (e *Executor) Enforce(item ObjectData, action Action) error {
	return e.check(item, action)
}

// CreateObjectDataCheck
// check permission of creating object_data
func (e *Executor) CreateObjectDataCheck(item ObjectData, ty ObjectType) error {
	return e.writeObjectDataCheck(item, ty)
}

// RecoverObjectDataCheck
// check permission of recover object_data
func (e *Executor) RecoverObjectDataCheck(item ObjectData, ty ObjectType) error {
	return e.writeObjectDataCheck(item, ty)
}

// UpdateObjectDataCheck
// check permission of updating object_data's
func (e *Executor) UpdateObjectDataCheck(item ObjectData, old ObjectData, ty ObjectType) error {
	list := []ObjectData{item}
	if item.GetObject().GetID() != old.GetObject().GetID() {
		list = append(list, old)
	}
	for _, v := range list {
		if err := e.writeObjectDataCheck(v, ty); err != nil {
			return err
		}
	}
	return nil
}

// DeleteObjectDataCheck
// check permission of deleting object_data
func (e *Executor) DeleteObjectDataCheck(item ObjectData, ty ObjectType) error {
	return e.writeObjectDataCheck(item, ty)
}

func (e *Executor) writeObjectDataCheck(item ObjectData, ty ObjectType) error {
	if err := e.check(item, Write); err != nil {
		return err
	}
	o := item.GetObject()
	if err := e.db.Take(o); err != nil {
		return ErrInValidObject
	}
	if o.GetObjectType() != ty {
		return ErrInValidObjectType
	}
	return nil
}
