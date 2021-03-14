package caskin

import "github.com/ahmetb/go-linq/v3"

// CreateObject
// if current user has object's write permission and there does not exist the object
// then create a new one
// 1. create a new object into metadata database
// 2. update object to parent's g2 in the domain
func (e *Executor) CreateObject(object Object) error {
	fn := func(domain Domain) error {
		if err := e.db.Create(object); err != nil {
			return err
		}
		updater := e.objectParentUpdater()
		return updater.update(object, domain)
	}

	return e.parentEntryFlowHandler(object, e.createObjectDataEntryCheck, e.newObject, fn)
}

// RecoverObject
// if current user has object's write permission and there exist the object but soft deleted
// then recover it
// 1. recover the soft delete one object at metadata database
// 2. update object to parent's g2 in the domain
func (e *Executor) RecoverObject(object Object) error {
	fn := func(domain Domain) error {
		if err := e.db.Recover(object); err != nil {
			return err
		}
		updater := e.objectParentUpdater()
		return updater.update(object, domain)
	}

	return e.parentEntryFlowHandler(object, e.recoverObjectDataEntryCheck, e.newObject, fn)
}

// DeleteObject
// if current user has object's write permission
// 1. delete object's g2 in the domain
// 2. delete object's p in the domain
// 3. soft delete one object in metadata database
// 4. dfs to delete all son of the object in the domain
func (e *Executor) DeleteObject(object Object) error {
	fn := func(domain Domain) error {
		deleter := newParentEntryDeleter(e.objectChildrenFn(), e.objectDeleteFn())
		return deleter.dfs(object, domain)
	}
	return e.parentEntryFlowHandler(object, e.deleteObjectDataEntryCheck, e.newObject, fn)
}

// UpdateObject
// if current user has object's write permission and there exist the object
// 1. update object's properties
// 2. update object to parent's g2 in the domain
func (e *Executor) UpdateObject(object Object) error {
	fn := func(domain Domain) error {
		if object.GetObjectType() == ObjectTypeObject &&
			object.GetObject().GetID() != object.GetID() {
			return ErrObjectTypeObjectIDMustBeItselfID
		}
		if err := e.db.Update(object); err != nil {
			return err
		}
		updater := e.objectParentUpdater()
		return updater.update(object, domain)
	}

	objectUpdateCheck := func(item ObjectData) error {
		tmp := e.newObject()
		if err := e.updateObjectDataEntryCheck(item, tmp); err != nil {
			return err
		}
		if item.(Object).GetObjectType() != tmp.(Object).GetObjectType() {
			return ErrInValidObjectType
		}
		return e.treeNodeParentCheck(tmp, e.newObject)
	}
	return e.parentEntryFlowHandler(object, objectUpdateCheck, e.newObject, fn)
}

// GetObject
// if current user has object's read permission
// 1. get objects by ty
func (e *Executor) GetObjects(ty ...ObjectType) ([]Object, error) {
	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	ds := e.e.GetObjectsInDomain(currentDomain)
	tree := getTree(ds)
	objects, err := e.db.GetObjectInDomain(currentDomain, ty...)
	if err != nil {
		return nil, err
	}

	os := e.filterWithNoError(currentUser, currentDomain, Read, objects)
	linq.From(os).ToSlice(&objects)

	for _, v := range objects {
		if p, ok := tree[v.GetID()]; ok {
			v.SetParentID(p)
		}
	}

	return objects, nil
}
