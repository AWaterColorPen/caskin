package caskin

import "github.com/ahmetb/go-linq/v3"

// CreateObject
// if current user has object's write permission and there does not exist the object
// then create a new one
// 1. create a new object into metadata database
// 2. update object to parent's g2 in the domain
func (e *Executor) CreateObject(object Object) error {
	if err := e.objectCheckFlow(object, e.ObjectDataCreateCheck); err != nil {
		return err
	}
	if err := e.db.Create(object); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	updater := e.objectParentUpdater()
	return updater.update(object, domain)
}

// RecoverObject
// if current user has object's write permission and there exist the object but soft deleted
// then recover it
// 1. recover the soft delete one object at metadata database
// 2. update object to parent's g2 in the domain
func (e *Executor) RecoverObject(object Object) error {
	if err := e.objectCheckFlow(object, e.ObjectDataRecoverCheck); err != nil {
		return err
	}
	if err := e.db.Recover(object); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	updater := e.objectParentUpdater()
	return updater.update(object, domain)
}

// DeleteObject
// if current user has object's write permission
// 1. delete object's g2 in the domain
// 2. delete object's p in the domain
// 3. soft delete one object in metadata database
// 4. dfs to delete all son of the object in the domain
func (e *Executor) DeleteObject(object Object) error {
	if err := e.objectCheckFlow(object, e.ObjectDataDeleteCheck); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	deleter := newParentEntryDeleter(e.objectChildrenFn(), e.objectDeleteFn())
	return deleter.dfs(object, domain)
}

// UpdateObject
// if current user has object's write permission and there exist the object
// 1. update object's properties
// 2. update object to parent's g2 in the domain
func (e *Executor) UpdateObject(object Object) error {
	tmp := e.factory.NewObject()
	if err := e.ObjectDataUpdateCheck(object, tmp); err != nil {
		return err
	}
	if err := e.objectTreeNodeParentCheck(tmp); err != nil {
		return err
	}

	if err := e.objectCheckFlow(object, func(ObjectData) error { return nil }); err != nil {
		return err
	}
	if object.GetObjectType() == ObjectTypeObject &&
		object.GetObject().GetID() != object.GetID() {
		return ErrObjectTypeObjectIDMustBeItselfID
	}
	if err := e.db.Update(object); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	updater := e.objectParentUpdater()
	return updater.update(object, domain)
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
