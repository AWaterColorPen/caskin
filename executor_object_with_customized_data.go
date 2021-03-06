package caskin

func (e *Executor) CreateObjectWithCustomizedData(customized CustomizedData, object Object) error {
	CustomizedData2Object(customized, object)
	return e.CreateObject(object)
}

func (e *Executor) RecoverObjectWithCustomizedData(customized CustomizedData, object Object) error {
	CustomizedData2Object(customized, object)
	return e.RecoverObject(object)
}

func (e *Executor) DeleteObjectWithCustomizedData(customized CustomizedData, object Object) error {
	object.SetObjectType(customized.GetObjectType())
	return e.DeleteObject(object)
}

func (e *Executor) UpdateObjectWithCustomizedData(customized CustomizedData, object Object) error {
	CustomizedData2Object(customized, object)
	return e.UpdateObject(object)
}
