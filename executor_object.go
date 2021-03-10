package caskin

// CreateObject
// if current user has object's write permission and there does not exist the object
// then create a new one
// 1. create a new object into metadata database
// 2. update object to parent's g2 in the domain
func (e *executor) CreateObject(object Object) error {
	return e.createOrRecoverObject(object, e.mdb.Create, e.mdb.Take)
}

// RecoverObject
// if current user has object's write permission and there exist the object but soft deleted
// then recover it
// 1. recover the soft delete one object at metadata database
// 2. update object to parent's g2 in the domain
func (e *executor) RecoverObject(object Object) error {
	return e.createOrRecoverObject(object, e.mdb.Recover, e.mdb.TakeUnscoped)
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
		if err := e.e.RemoveObjectInDomain(object, domain); err != nil {
			return err
		}
		return e.mdb.DeleteObjectByID(object.GetID())
	}

	return e.writeObject(object, fn)
}

// UpdateObject
// if current user has object's write permission and there exist the object
// 1. update object's properties
// 2. update object to parent's g2 in the domain
func (e *executor) UpdateObject(object Object) error {
	return e.writeObject(object, e.mdb.Update)
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

func (e *executor) createOrRecoverObject(object Object, fn func(interface{}) error, takeObject func(interface{}) error) error {
	if err := e.mdb.Take(object); err == nil {
		return ErrAlreadyExists
	}

	_, domain, err := e.provider.Get()
	if err != nil {
		return err
	}

	take := func(id uint64) (parentEntry, error) {
		o := e.factory.NewObject()
		o.SetID(id)
		o.SetDomainID(domain.GetID())
		err := takeObject(o)
		return o, err
	}

	if err := e.checkParentEntryWrite(object, take); err != nil {
		return err
	}

	object.SetDomainID(domain.GetID())
	return fn(object)
}

func (e *executor) writeObject(object Object, fn func(interface{}) error) error {
	if err := isValid(object); err != nil {
		return err
	}

	tmpObject := e.factory.NewObject()
	tmpObject.SetID(object.GetID())
	if err := e.mdb.Take(tmpObject); err != nil {
		return ErrNotExists
	}

	if tmpObject.GetObjectType() != object.GetObjectType() {
		return ErrNotValidObjectType
	}

	_, domain, err := e.provider.Get()
	if err != nil {
		return err
	}

	take := func(id uint64) (parentEntry, error) {
		o := e.factory.NewObject()
		o.SetID(id)
		o.SetDomainID(domain.GetID())
		err := e.mdb.Take(o)
		return o, err
	}

	if err := e.checkParentEntryWrite(tmpObject, take); err != nil {
		return err
	}
	if err := e.checkParentEntryWrite(object, take); err != nil {
		return err
	}

	object.SetDomainID(domain.GetID())
	return fn(object)
}

// 创建object的逻辑
func (e *executor) createObject(object Object) error {
	_, domain, err := e.provider.Get()
	if err != nil {
		return nil
	}
	// 1. 首先查看当前object是否存在，如果存在，那么直接已经存在的错误
	if err := e.mdb.Take(object); err == nil {
		return ErrAlreadyExists
	}
	// 2. 否则，查看对当前的object是否有权限
	if err := e.check(Write, object); err != nil {
		return ErrNoWritePermission

	}
	// 3. 如果当前object有parent，那么查看是否对parent有权限
	parent := e.factory.NewObject()
	if object.GetParentID() != 0 {
		parent.SetID(object.GetParentID())
		parent.SetDomainID(domain.GetID())
		// 3.1. 去数据库里面拿出数据，包含有objectId
		if err := e.mdb.Take(parent); err != nil {
			return err
		}
		if err := e.check(Write, parent); err != nil {
			return ErrNoWritePermission
		}
	}
	// 4. 创建object
	object.SetDomainID(domain.GetID())
	if err := e.mdb.Create(object); err != nil {
		return err
	}
	// 3.2. 处理parent的关系
	if parent.GetID() != 0 {
		return e.e.AddParentForObjectInDomain(object, parent, domain)
	}
	return nil
}

// 恢复object的逻辑
func (e *executor) recoverObject(object Object) error {
	_, domain, err := e.provider.Get()
	if err != nil {
		return nil
	}
	// 1. 首先查看当前的object是否存在，如果存在，那么不需要恢复，直接报错
	if err := e.mdb.Take(object); err == nil {
		return ErrAlreadyExists
	}
	// 2. 否则，查看记录是否在数据库里，使用unscoped来查看，如果该记录存在，那么可以recover
	if err := e.mdb.TakeUnscoped(object); err != nil {
		return err
	}
	// 3. 查看对当前object是否有权限
	if err := e.check(Write, object); err != nil {
		return ErrNoWritePermission
	}
	// 4. 如果当前的object有parent，那么查看对parent是否有权限
	parent := e.factory.NewObject()
	if object.GetParentID() != 0 {
		parent.SetID(object.GetParentID())
		parent.SetDomainID(domain.GetID())
		if err := e.mdb.Take(parent); err != nil {
			return err
		}
		if err := e.check(Write, parent); err != nil {
			return ErrNoWritePermission
		}
	}
	object.SetDomainID(domain.GetID())
	if err := e.mdb.Recover(object); err != nil {
		return err
	}
	if parent.GetID() != 0 {
		return e.e.AddParentForObjectInDomain(object, parent, domain)
	}
	return nil
}

// 修改object的逻辑
func (e *executor) update(object Object) error {
	_, domain, err := e.provider.Get()
	if err != nil {
		return nil
	}
	// 1. 首先查看id是否合法
	if err := isValid(object); err != nil {
		return err
	}
	// 2. 首先查看当前的object是否存在，需要构造一个tmp的去查，需要存在才能够update
	tmpObject := e.factory.NewObject()
	tmpObject.SetID(object.GetID())
	if err := e.mdb.Take(tmpObject); err != nil {
		return ErrNotExists
	}
	// 3. 然后查看tmp的相关权限
	if err := e.check(Write, tmpObject); err != nil {
		return ErrNoWritePermission
	}
	tmpParents := e.e.GetParentsForObjectInDomain(tmpObject, domain)
	for _, v := range tmpParents {
		if err := e.check(Write, v); err != nil {
			return ErrNoWritePermission
		}
	}
	// 4. 查看当前object的相关权限
	if err := e.check(Write, object); err != nil {
		return ErrNoWritePermission
	}
	parent := e.factory.NewObject()
	if object.GetParentID() != 0 {
		parent.SetID(object.GetParentID())
		parent.SetDomainID(domain.GetID())
		if err := e.mdb.Take(parent); err != nil {
			return err
		}
		if err := e.check(Write, parent); err != nil {
			return ErrNoWritePermission
		}
	}
	object.SetDomainID(domain.GetID())
	if err := e.mdb.Update(object); err != nil {
		return err
	}
	if parent.GetID() != 0 {
		return e.e.AddParentForObjectInDomain(object, parent, domain)
	}
	return nil
}

// 删除object的逻辑
func (e *executor) delete(object Object) error {
	_, domain, err := e.provider.Get()
	if err != nil {
		return nil
	}
	// 1. 首先查看id是否合法
	if err := isValid(object); err != nil {
		return err
	}
	// 2. 首先查看当前的object是否存在,需要存在才能够delete
	if err := e.mdb.Take(object); err != nil {
		return ErrNotExists
	}
	// 3. 然后查看object的相关权限
	if err := e.check(Write, object); err != nil {
		return err
	}
	if object.GetParentID() != 0 {
		parent := e.factory.NewObject()
		parent.SetID(object.GetParentID())
		parent.SetDomainID(domain.GetID())
		if err := e.mdb.Take(parent); err != nil {
			return err
		}
		if err := e.check(Write, parent); err != nil {
			return err
		}
	}
	// 4. 删除其parent的关系
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
