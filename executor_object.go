package caskin

import "github.com/ahmetb/go-linq/v3"

// CreateObject
// if current user has object's write permission and there does not exist the object
// then create a new one
// 1. create a new object into metadata database
// 2. Run object to parent's g2 in the domain
func (e *Executor) CreateObject(object Object) error {
	if err := e.ObjectDataCreateCheck(object, ObjectTypeObject); err != nil {
		return err
	}
	if err := e.ObjectTreeNodeParentCheck(object); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	object.SetDomainID(domain.GetID())
	if err := e.DB.Create(object); err != nil {
		return err
	}
	updater := e.DefaultObjectUpdater()
	return updater.Run(object, domain)
}

// RecoverObject
// if current user has object's write permission and there exist the object but soft deleted
// then recover it
// 1. recover the soft delete one object at metadata database
// 2. Run object to parent's g2 in the domain
func (e *Executor) RecoverObject(object Object) error {
	if err := e.ObjectDataRecoverCheck(object); err != nil {
		return err
	}
	if err := e.ObjectTreeNodeParentCheck(object); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	object.SetDomainID(domain.GetID())
	if err := e.DB.Recover(object); err != nil {
		return err
	}
	updater := e.DefaultObjectUpdater()
	return updater.Run(object, domain)
}

// DeleteObject
// if current user has object's write permission
// 1. delete object's g2 in the domain
// 2. delete object's p in the domain
// 3. soft delete one object in metadata database
// 4. Run to delete all son of the object in the domain
func (e *Executor) DeleteObject(object Object) error {
	if err := e.ObjectDataDeleteCheck(object); err != nil {
		return err
	}
	if err := e.ObjectTreeNodeParentCheck(object); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	object.SetDomainID(domain.GetID())
	deleter := NewTreeNodeEntryDeleter(e.DefaultObjectChildrenGetFunc(), e.DefaultObjectDeleteFunc())
	return deleter.Run(object, domain)
}

// UpdateObject
// if current user has object's write permission and there exist the object
// 1. Run object's properties
// 2. Run object to parent's g2 in the domain
func (e *Executor) UpdateObject(object Object) error {
	if err := e.objectTreeNodeUpdateCheck(object, e.factory.NewObject()); err != nil {
		return err
	}
	if err := e.ObjectTreeNodeParentCheck(object); err != nil {
		return err
	}
	if err := isObjectTypeObjectIDBeSelfIDCheck(object); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	object.SetDomainID(domain.GetID())
	if err := e.DB.Update(object); err != nil {
		return err
	}
	updater := e.DefaultObjectUpdater()
	return updater.Run(object, domain)
}

// GetObject
// if current user has object's read permission
// 1. get objects by ty
func (e *Executor) GetObjects(ty ...ObjectType) ([]Object, error) {
	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}
	objects, err := e.DB.GetObjectInDomain(currentDomain, ty...)
	if err != nil {
		return nil, err
	}
	os := e.filterWithNoError(currentUser, currentDomain, Read, objects)
	linq.From(os).ToSlice(&objects)
	return objects, nil
}
