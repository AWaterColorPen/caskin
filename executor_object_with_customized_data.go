package caskin

func (e *Executor) CreateObjectWithCustomizedData(customized ObjectCustomizedData, object Object) error {
	ObjectCustomizedData2Object(customized, object)
	return e.CreateObject(object)
}

func (e *Executor) RecoverObjectWithCustomizedData(customized ObjectCustomizedData) error {
	object := e.newObject().(Object)
	ObjectCustomizedData2Object(customized, object)
	return e.RecoverObject(object)
}

func (e *Executor) DeleteObjectWithCustomizedData(customized ObjectCustomizedData, object Object) error {
	if !ObjectCustomizedDataEqualObject(customized, object) {
		return ErrCustomizedDataIsNotBelongToObject
	}
	return e.DeleteObject(object)
}

func (e *Executor) UpdateObjectWithCustomizedData(customized ObjectCustomizedData, object Object) error {
	ObjectCustomizedData2Object(customized, object)
	return e.UpdateObject(object)
}
