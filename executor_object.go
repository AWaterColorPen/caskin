package caskin

import "github.com/ahmetb/go-linq/v3"

// ObjectCreate
// if there does not exist the object then create a new one
// 1. current user has manage permission of object's parent
// 2. create a new object into metadata database
// 3. set object to parent's g2 in the domain
func (e *Executor) ObjectCreate(user User, domain Domain, object Object) error {
	if err := e.DBCreateCheck(object); err != nil {
		return err
	}
	if isRoot(object) {
		return ErrEmptyParentIdOrNotSuperadmin
	}
	if err := e.ObjectParentCheck(user, domain, object); err != nil {
		return err
	}
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
	if err := e.ObjectParentCheck(object); err != nil {
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
	if err := e.ObjectParentCheck(object); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	object.SetDomainID(domain.GetID())
	deleter := NewTreeNodeDeleter(e.DefaultObjectChildrenGetFunc(), e.DefaultObjectDeleteFunc())
	return deleter.Run(object, domain)
}

// UpdateObject
// if current user has object's write permission and there exist the object
// 1. Run object's properties
// 2. Run object to parent's g2 in the domain
func (e *Executor) UpdateObject(object Object) error {
	if err := e.ObjectTreeNodeUpdateCheck(object, e.factory.NewObject()); err != nil {
		return err
	}
	if err := e.ObjectParentCheck(object); err != nil {
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

// GetObjects
// if current user has object's read permission
// 1. get objects by type
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

// GetExplicitObjects
// if current user has explicit object's action permission
// 1. get objects by type and action
func (e *Executor) GetExplicitObjects(action Action, ty ...ObjectType) ([]Object, error) {
	_, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}
	objects, err := e.DB.GetObjectInDomain(currentDomain, ty...)
	if err != nil {
		return nil, err
	}
	return e.FilterObject(objects, action)
}
