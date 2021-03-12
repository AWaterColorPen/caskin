package caskin

import "github.com/ahmetb/go-linq/v3"

// FilterObjectData
// filter object_data with action
func (e *executor) FilterObjectData(source interface{}, action Action) ([]ObjectData, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	var result []ObjectData
	linq.From(source).Where(func(v interface{}) bool {
		return Check(e.e, u, d, v.(ObjectData), action)
	}).ToSlice(&result)
	return result, nil
}

// Enforce
// check permission of object_data with action
func (e *executor) Enforce(item ObjectData, action Action) error {
	return e.check(item, action)
}

// // CreateObjectDataCheck
// // if current user has object's write permission and there does not exist the object
// // then create a new one
// // 1. create a new object into metadata database
// // 2. update object to parent's g2 in the domain
// func (e *executor) CreateObjectDataCheck(item ObjectData, ty ObjectType) error {
// 	if err := e.check(item, Write); err != nil {
// 		return err
// 	}
// 	o := item.GetObject()
// 	if err := e.mdb.Take(o); err != nil {
// 		return E
// 	}
// 	if o.GetObjectType() != ty {
// 		return
// 	}
// 	return
// }

//
// // RecoverObject
// // if current user has object's write permission and there exist the object but soft deleted
// // then recover it
// // 1. recover the soft delete one object at metadata database
// // 2. update object to parent's g2 in the domain
// func (e *executor) RecoverObject(object Object) error {
// 	fn := func(domain Domain) error {
// 		if err := e.mdb.Recover(object); err != nil {
// 			return err
// 		}
// 		updater := e.objectParentUpdater()
// 		return updater.update(object, domain)
// 	}
//
// 	return e.parentEntryFlowHandler(object, e.recoverObjectDataEntryCheck, e.newObject, fn)
// }
//
// // DeleteObject
// // if current user has object's write permission
// // 1. delete object's g2 in the domain
// // 2. delete object's p in the domain
// // 3. soft delete one object in metadata database
// // 4. dfs to delete all son of the object in the domain
// func (e *executor) DeleteObject(object Object) error {
// 	fn := func(domain Domain) error {
// 		deleter := newParentEntryDeleter(e.objectChildrenFn(), e.objectDeleteFn())
// 		return deleter.dfs(object, domain)
// 	}
//
// 	return e.parentEntryFlowHandler(object, e.deleteObjectDataEntryCheck, e.newObject, fn)
// }
//
// // UpdateObject
// // if current user has object's write permission and there exist the object
// // 1. update object's properties
// // 2. update object to parent's g2 in the domain
// func (e *executor) UpdateObject(object Object) error {
// 	fn := func(domain Domain) error {
// 		if object.GetObjectType() == ObjectTypeObject &&
// 			object.GetObject().GetID() != object.GetID() {
// 			return ErrObjectTypeObjectIDMustBeItselfID
// 		}
// 		if err := e.mdb.Update(object); err != nil {
// 			return err
// 		}
// 		updater := e.objectParentUpdater()
// 		return updater.update(object, domain)
// 	}
//
// 	objectUpdateCheck := func(item objectDataEntry) error {
// 		tmp := e.newObject()
// 		if err := e.updateObjectDataEntryCheck(item, tmp); err != nil {
// 			return err
// 		}
// 		return e.parentEntryCheck(tmp, e.objectParentsFn())
// 	}
// 	return e.parentEntryFlowHandler(object, objectUpdateCheck, e.newObject, fn)
// }
//
// // GetObject
// // if current user has object's read permission
// // 1. get objects by ty
// func (e *executor) GetObjects(ty ...ObjectType) ([]Object, error) {
// 	currentUser, currentDomain, err := e.provider.Get()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	ds := e.e.GetObjectsInDomain(currentDomain)
// 	tree := getTree(ds)
// 	objects, err := e.mdb.GetObjectInDomain(currentDomain, ty...)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	os := e.filterWithNoError(currentUser, currentDomain, Read, objects)
// 	objects = []Object{}
// 	for _, v := range os {
// 		objects = append(objects, v.(Object))
// 	}
//
// 	for _, v := range objects {
// 		if p, ok := tree[v.GetID()]; ok {
// 			v.SetParentID(p)
// 		}
// 	}
//
// 	return objects, nil
// }
