package caskin

// ObjectCreate
// if there does not exist the object then create a new one
// 1. current user has manage permission of object's parent
// 2. create a new object into metadata database
// 3. set object to parent's g2 in the domain
func (e *Executor) ObjectCreate(user User, domain Domain, object Object) error {
	if err := e.DBCreateCheck(object); err != nil {
		return err
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

// ObjectRecover
// if there exist the object but soft deleted then recover it
// 1. current user has manage permission of object's parent
// 2. recover the soft delete one object at metadata database
// 3. set object to parent's g2 in the domain
func (e *Executor) ObjectRecover(user User, domain Domain, object Object) error {
	if err := e.DBRecoverCheck(object); err != nil {
		return err
	}
	if err := e.ObjectParentCheck(user, domain, object); err != nil {
		return err
	}
	object.SetDomainID(domain.GetID())
	if err := e.DB.Recover(object); err != nil {
		return err
	}
	updater := e.DefaultObjectUpdater()
	return updater.Run(object, domain)
}

// ObjectDelete
// if there exist the object
// 1. current user has manage permission of object's parent
// 2. delete object's g2 in the domain
// 3. delete object's p in the domain
// 4. soft delete one object in metadata database
// 5. delete all son of the object in the domain
func (e *Executor) ObjectDelete(user User, domain Domain, object Object) error {
	if err := e.IDInterfaceDeleteCheck(object); err != nil {
		return err
	}
	if err := e.ObjectParentCheck(user, domain, object); err != nil {
		return err
	}
	object.SetDomainID(domain.GetID())
	deleter := NewTreeNodeDeleter(e.DefaultObjectChildrenGetFunc(), e.DefaultObjectDeleteFunc())
	return deleter.Run(object, domain)
}

// ObjectUpdate
// if there exist the object
// 1. current user has manage permission of object's parent to change parent_id
//    current user has manage permission of object to change properties
// 2. update object's properties
// 3. update object to parent's g2 in the domain
func (e *Executor) ObjectUpdate(user User, domain Domain, object Object) error {
	if err := e.ObjectUpdateCheck(user, domain, object); err != nil {
		return err
	}
	if err := e.ObjectParentCheck(user, domain, object); err != nil {
		return err
	}
	object.SetDomainID(domain.GetID())
	if err := e.DB.Update(object); err != nil {
		return err
	}
	updater := e.DefaultObjectUpdater()
	return updater.Run(object, domain)
}

// ObjectGet
// get choose object
// 1. current user has manage permission of object
//    manage permission for admin to manage
//    read/write permission for customer to get directory
// 2. get object by type
func (e *Executor) ObjectGet(user User, domain Domain, action Action, ty ...ObjectType) ([]Object, error) {
	objects, err := e.DB.GetObjectInDomain(domain, ty...)
	if err != nil {
		return nil, err
	}
	return Filter(e.Enforcer, user, domain, action, objects), nil
}
