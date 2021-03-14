package caskin

func (e *Executor) CreateObjectWithCustomizedData(customized CustomizedData, object Object) error {
	CustomizedData2Object(customized, object)
	return e.CreateObject(object)
}

func (e *Executor) RecoverObjectWithCustomizedData(customized CustomizedData) error {
	object := e.newObject().(Object)
	CustomizedData2Object(customized, object)
	return e.RecoverObject(object)
}

func (e *Executor) DeleteObjectWithCustomizedData(customized CustomizedData, object Object) error {
	if !CustomizedDataEqualObject(customized, object) {
		return ErrCustomizedDataIsNotBelongToObject
	}
	return e.DeleteObject(object)
}

func (e *Executor) UpdateObjectWithCustomizedData(customized CustomizedData, object Object) error {
	CustomizedData2Object(customized, object)
	return e.UpdateObject(object)
}
