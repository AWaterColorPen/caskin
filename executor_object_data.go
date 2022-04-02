package caskin

func (e *Executor) ObjectDataWriteCheck(user User, domain Domain, item ObjectData, ty ObjectType) error {
	if err := e.CheckObjectData(user, domain, item, Write); err != nil {
		return err
	}
	o := item.GetObject()
	if err := e.DB.Take(o); err != nil {
		return ErrInValidObject
	}
	if o.GetObjectType() != ty {
		return ErrInValidObjectType
	}
	return nil
}

func (e *Executor) ObjectDataCreateCheck(user User, domain Domain, item ObjectData, ty ObjectType) error {
	if err := e.DBCreateCheck(item); err != nil {
		return err
	}
	return e.ObjectDataWriteCheck(user, domain, item, ty)
}

func (e *Executor) ObjectDataRecoverCheck(user User, domain Domain, item ObjectData) error {
	if err := e.DBRecoverCheck(item); err != nil {
		return err
	}
	return e.CheckObjectData(user, domain, item, Write)
}

func (e *Executor) ObjectDataDeleteCheck(user User, domain Domain, item ObjectData) error {
	if err := e.IDInterfaceDeleteCheck(item); err != nil {
		return err
	}
	return e.CheckObjectData(user, domain, item, Write)
}

func (e *Executor) ObjectDataUpdateCheck(user User, domain Domain, item ObjectData, ty ObjectType) error {
	old := createByE(item)
	if err := e.IDInterfaceUpdateCheck(item, old); err != nil {
		return err
	}
	if err := e.ObjectDataWriteCheck(user, domain, old, ty); err != nil {
		return err
	}
	if item.GetObject().GetID() != old.GetObject().GetID() {
		return e.ObjectDataWriteCheck(user, domain, item, ty)
	}
	return nil
}

func (e *Executor) ObjectDataModifyCheck(user User, domain Domain, item ObjectData) error {
	if err := e.IDInterfaceModifyCheck(item); err != nil {
		return err
	}
	return e.CheckObjectData(user, domain, item, Write)
}

func (e *Executor) ObjectDataGetCheck(user User, domain Domain, item ObjectData) error {
	if err := e.IDInterfaceGetCheck(item); err != nil {
		return err
	}
	return e.CheckObjectData(user, domain, item, Read)
}
