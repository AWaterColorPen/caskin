package caskin

// CreateObject
// if current user has object's write permission and there does not exist the object
// then create a new one
// 1. create a new object into metadata database
// 2. update object to parent's g2 in the domain
func (e *executor) CreateObject(object Object) error {
	fn := func() error {
		_, domain, err := e.provider.Get()
		if err != nil {
			return err
		}

		object.SetDomainID(domain.GetID())
		if err := e.mdb.Create(object); err != nil {
			return err
		}

		if object.GetParentID() != 0 {
			parent := e.factory.NewObject()
			parent.SetID(object.GetParentID())
			return e.e.AddParentForObjectInDomain(object, parent, domain)
		}
		return nil
	}

	return e.flowHandler(object, e.createCheck, e.newObject, fn)
}

// RecoverObject
// if current user has object's write permission and there exist the object but soft deleted
// then recover it
// 1. recover the soft delete one object at metadata database
// 2. update object to parent's g2 in the domain
func (e *executor) RecoverObject(object Object) error {
	fn := func() error {
		_, domain, err := e.provider.Get()
		if err != nil {
			return err
		}
		object.SetDomainID(domain.GetID())
		if err := e.mdb.Recover(object); err != nil {
			return err
		}
		if object.GetParentID() != 0 {
			parent := e.factory.NewObject()
			parent.SetID(object.GetParentID())
			return e.e.AddParentForObjectInDomain(object, parent, domain)
		}
		return nil
	}

	return e.flowHandler(object, e.recoverEntryCheck, e.newObject, fn)
}

// DeleteObject
// if current user has object's write permission
// 1. delete object's g2 in the domain
// 2. delete object's p in the domain
// 3. soft delete one object in metadata database
func (e *executor) DeleteObject(object Object) error {
	fn := func() error {
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
	fn := func() error {
		_, domain, err := e.provider.Get()
		if err != nil {
			return err
		}
		object.SetDomainID(domain.GetID())
		if err := e.mdb.Update(object); err != nil {
			return err
		}
		// TODO 这里的parent关系可能会有错误
		if object.GetParentID() != 0 {
			parent := e.factory.NewObject()
			parent.SetID(object.GetParentID())
			return e.e.AddParentForObjectInDomain(object, parent, domain)
		}
		return nil
	}
	return e.flowHandler(object, e.updateEntryCheck, e.newObject, fn)
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

func (e *executor) createEntryCheck(entry objectDataEntry) error {
	if err := e.mdb.Take(entry); err == nil {
		return ErrAlreadyExists
	}
	return e.check(Write, entry)
}

func (e *executor) recoverEntryCheck(entry objectDataEntry) error {
	if err := e.mdb.Take(entry); err == nil {
		return ErrAlreadyExists
	}
	if err := e.mdb.TakeUnscoped(entry); err != nil {
		return err
	}
	return e.check(Write, entry)
}

func (e *executor) updateEntryCheck(entry objectDataEntry) error {
	_, domain, err := e.provider.Get()
	if err != nil {
		return err
	}

	if err := isValid(entry); err != nil {
		return err
	}

	tmpObject := e.factory.NewObject()
	tmpObject.SetID(entry.GetID())
	if err := e.mdb.Take(tmpObject); err != nil {
		return ErrNotExists
	}

	tmpParents := e.e.GetParentsForObjectInDomain(tmpObject, domain)
	for _, v := range tmpParents {
		if err := e.check(Write, v); err != nil {
			return ErrNoWritePermission
		}
	}

	return e.check(Write, tmpObject)
}

func (e *executor) deleteEntryCheck(entry objectDataEntry) error {
	if err := isValid(entry); err != nil {
		return err
	}
	if err := e.mdb.Take(entry); err != nil {
		return ErrNotExists
	}
	return e.check(Write, entry)
}

func (e *executor) flowHandler(entry parentEntry,
	check func(objectDataEntry) error,
	newEntry func() parentEntry,
	handle func() error) error {
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

	return handle()
}
