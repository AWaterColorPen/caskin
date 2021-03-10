package caskin

func (e *executor) newObject() parentEntry {
	return e.factory.NewObject()
}

// CreateObject
// if current user has object's write permission and there does not exist the object
// then create a new one
// 1. create a new object into metadata database
// 2. update object to parent's g2 in the domain
func (e *executor) CreateObject(object Object) error {
	return e.flowHandler(object, e.createEntryCheck, e.newObject, e.mdb.Create)
}

// RecoverObject
// if current user has object's write permission and there exist the object but soft deleted
// then recover it
// 1. recover the soft delete one object at metadata database
// 2. update object to parent's g2 in the domain
func (e *executor) RecoverObject(object Object) error {
	return e.flowHandler(object, e.recoverEntryCheck, e.newObject, e.mdb.Recover)
}

// DeleteObject
// if current user has object's write permission
// 1. delete object's g2 in the domain
// 2. delete object's p in the domain
// 3. soft delete one object in metadata database
func (e *executor) DeleteObject(object Object) error {
	fn := func(interface{}) error {
		_, domain, err := e.provider.Get()
		if err != nil {
			return err
		}

		parent := e.e.GetParentsForObjectInDomain(object, domain)
		for _, v := range parent {
			if err := e.e.RemoveParentForObjectInDomain(object, v, domain); err != nil {
				return err
			}
		}
		// 5. 删除object在domain中的关系
		object.SetDomainID(domain.GetID())
		if err := e.e.RemoveObjectInDomain(object, domain); err != nil {
			return err
		}
		// 6. 删除object的数据
		return e.mdb.DeleteObjectByID(object.GetID())
	}

	return e.flowHandler(object, e.deleteEntryCheck, e.newObject, fn)
}

// UpdateObject
// if current user has object's write permission and there exist the object
// 1. update object's properties
// 2. update object to parent's g2 in the domain
func (e *executor) UpdateObject(object Object) error {
	return e.flowHandler(object, e.updateEntryCheck, e.newObject, e.mdb.Update)
}

// GetObject
// if current user has object's read permission
// 1. get objects by ty
func (e *executor) GetObjects(ty ...ObjectType) ([]Object, error) {
	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	ds := e.e.GetObjectsInDomain(currentDomain)
	tree := getTree(ds)
	objects, err := e.mdb.GetObjectInDomain(currentDomain, ty...)
	if err != nil {
		return nil, err
	}

	os := e.filterWithNoError(currentUser, currentDomain, Read, objects)
	objects = []Object{}
	for _, v := range os {
		objects = append(objects, v.(Object))
	}

	for _, v := range objects {
		if p, ok := tree[v.GetID()]; ok {
			v.SetParentID(p)
		}
	}

	return objects, nil
}

func (e *executor) createEntryCheck(entry parentEntry) error {
	if err := e.mdb.Take(entry); err == nil {
		return ErrAlreadyExists
	}
	return e.check(Write, entry)
}

func (e *executor) recoverEntryCheck(entry parentEntry) error {
	// 1. 首先查看当前的object是否存在，如果存在，那么不需要恢复，直接报错
	if err := e.mdb.Take(entry); err == nil {
		return ErrAlreadyExists
	}
	// 2. 否则，查看记录是否在数据库里，使用unscoped来查看，如果该记录存在，那么可以recover
	if err := e.mdb.TakeUnscoped(entry); err != nil {
		return err
	}
	// 3. 查看对当前object是否有权限
	return e.check(Write, entry)
}

func (e *executor) updateEntryCheck(entry parentEntry) error {
	// 1. 首先查看id是否合法
	if err := isValid(entry); err != nil {
		return err
	}
	// 2. 首先查看当前的object是否存在，需要构造一个tmp的去查，需要存在才能够update
	tmpObject := e.factory.NewObject()
	tmpObject.SetID(entry.GetID())
	if err := e.mdb.Take(tmpObject); err != nil {
		return ErrNotExists
	}
	// 3. 然后查看tmp的相关权限
	return e.check(Write, tmpObject)
}

func (e *executor) deleteEntryCheck(entry parentEntry) error {
	// 1. 首先查看id是否合法
	if err := isValid(entry); err != nil {
		return err
	}
	// 2. 首先查看当前的object是否存在,需要存在才能够delete
	if err := e.mdb.Take(entry); err != nil {
		return ErrNotExists
	}
	// 3. 然后查看object的相关权限
	return e.check(Write, entry)
}

func (e *executor) flowHandler(entry parentEntry,
	check func(parentEntry) error,
	newEntry func() parentEntry,
	handle func(interface{}) error) error {
	if err := check(entry); err != nil {
		return err
	}

	_, domain, err := e.provider.Get()
	if err != nil {
		return err
	}
	parent := newEntry()
	if entry.GetParentID() != 0 {
		parent.SetID(entry.GetParentID())
		parent.SetDomainID(domain.GetID())
		if err := e.mdb.Take(parent); err != nil {
			return err
		}
		if err := e.check(Write, parent); err != nil {
			return ErrNoWritePermission
		}
	}
	entry.SetDomainID(domain.GetID())
	return handle(entry)
}
