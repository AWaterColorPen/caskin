package caskin

// CreateObject
// if current user has object's write permission and there does not exist the object
// then create a new one
// 1. create a new object into metadata database
func (e *executor) CreateObject(object Object) error {
	return e.createOrRecoverObject(object, e.mdb.CreateObject)
}

// RecoverObject
// if current user has object's write permission and there exist the object but soft deleted
// then recover it
// 1. recover the soft delete one object at metadata database
func (e *executor) RecoverObject(object Object) error {
	return e.createOrRecoverObject(object, e.mdb.RecoverObject)
}

// DeleteObject
// if current user has object's write permission
// 1. delete object's g in the domain
// 2. delete object's p in the domain
// 3. soft delete one object in metadata database
func (e *executor) DeleteObject(object Object) error {
	fn := func(o Object) error {
		_, domain, err := e.provider.Get()
		if err != nil {
			return err
		}
		if err := e.e.RemoveObjectInDomain(o, domain); err != nil {
			return err
		}
		return e.mdb.DeleteObjectByID(o.GetID())
	}

	return e.writeObject(object, fn)
}

// UpdateObject
// if current user has object's write permission and there exist the object
// 1. update object's properties
func (e *executor) UpdateObject(object Object) error {
	return e.writeObject(object, e.mdb.UpdateObject)
}

// GetObject
// if current user has object's read permission
// 1. get objects by ty
func (e *executor) GetObject(ty ...ObjectType) ([]Object, error) {
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

func (e *executor) createOrRecoverObject(object Object, fn func(Object) error) error {
	if err := e.mdb.TakeObject(object); err == nil {
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
		err := e.mdb.TakeObject(o)
		return o, err
	}

	if err := e.checkParentEntryWrite(object, take); err != nil {
		return err
	}

	object.SetDomainID(domain.GetID())
	return fn(object)
}

func (e *executor) writeObject(role Object, fn func(Object) error) error {
	if err := isValid(role); err != nil {
		return err
	}

	if err := e.mdb.TakeObject(role); err != nil {
		return ErrNotExists
	}

	_, domain, err := e.provider.Get()
	if err != nil {
		return err
	}

	take := func(id uint64) (parentEntry, error) {
		r := e.factory.NewObject()
		r.SetDomainID(domain.GetID())
		err := e.mdb.TakeObject(r)
		return r, err
	}

	if err := e.checkParentEntryWrite(role, take); err != nil {
		return err
	}

	role.SetDomainID(domain.GetID())
	return fn(role)
}
