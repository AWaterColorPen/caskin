package caskin

type BindExecutor struct {
	e  executor
	db MetaDBBindObjectAPI
}

// CreateBindObjectData
// if current user has object's write permission and there does not exist the pair of object_data and bind object
// then create new
// 1. create a new pair into metadata database
// 2. update bind object to parent's g2 in the domain
func (b *BindExecutor) CreateBindObjectData(item ObjectData, bind Object, ty ObjectType) error {
	check := func(objectDataEntry) error {
		if bind.GetObjectType() != ty {
			return ErrInValidObjectType
		}
		if err := b.e.createObjectDataEntryCheck(bind); err != nil {
			return err
		}
		return b.e.createCheck(item)
	}

	fn := func(domain Domain) error {
		if err := b.db.Create(item, bind); err != nil {
			return err
		}
		updater := b.e.objectParentUpdater()
		return updater.update(bind, domain)
	}

	return b.e.parentEntryFlowHandler(bind, check, b.e.newObject, fn)
}

// RecoverBindObjectData
// if current user has object's write permission and there exist the pair of object_data and bind object
// then recover them
// 1. recover the soft delete pair at metadata database
// 2. update bind object to parent's g2 in the domain
func (b *BindExecutor) RecoverBindObjectData(item ObjectData, bind Object, ty ObjectType) error {
	check := func(objectDataEntry) error {
		if bind.GetObjectType() != ty {
			return ErrInValidObjectType
		}
		if err := b.e.recoverObjectDataEntryCheck(bind); err != nil {
			return err
		}
		return b.e.recoverCheck(item)
	}

	fn := func(domain Domain) error {
		if err := b.db.Recover(item, bind); err != nil {
			return err
		}
		updater := b.e.objectParentUpdater()
		return updater.update(bind, domain)
	}

	return b.e.parentEntryFlowHandler(bind, check, b.e.newObject, fn)
}

// DeleteBindObjectData
// if current user has object's write permission
// 1. delete object's g2 in the domain
// 2. delete object's p in the domain
// 3. soft delete pair in metadata database
// 4. dfs to delete all son of the pairs in the domain
func (b *BindExecutor) DeleteBindObjectData(item ObjectData, bind Object, ty ObjectType) error {
	check := func(objectDataEntry) error {
		if bind.GetObjectType() != ty {
			return ErrInValidObjectType
		}
		if err := b.e.deleteObjectDataEntryCheck(bind); err != nil {
			return err
		}
		return b.e.deleteEntryCheck(item)
	}

	delFn := func(p parentEntry, d Domain) error {
		if err := b.e.e.RemoveObjectInDomain(p.(Object), d); err != nil {
			return err
		}
		return b.db.DeleteByID(item.GetID(), p.GetID())
	}

	fn := func(domain Domain) error {
		deleter := newParentEntryDeleter(b.e.objectChildrenFn(), delFn)
		return deleter.dfs(bind, domain)
	}

	return b.e.parentEntryFlowHandler(bind, check, b.e.newObject, fn)
}

// UpdateBindObjectData
// if current user has object's write permission and there exist the pair of object_data and bind object
// 1. update pair's properties
// 2. update object to parent's g2 in the domain
func (b *BindExecutor) UpdateBindObjectData(item ObjectData, bind Object, ty ObjectType) error {
	check := func(objectDataEntry) error {
		if bind.GetObjectType() != ty {
			return ErrInValidObjectType
		}
		tmp := b.e.newObject()
		if err := b.e.updateObjectDataEntryCheck(bind, tmp); err != nil {
			return err
		}
		if err := b.e.updateObjectDataEntryCheck(bind, tmp); err != nil {
			return err
		}
		if err := b.e.parentEntryCheck(bind, b.e.objectParentsFn()); err != nil {
			return err
		}

		// TODO
		// b.e.updateEntryCheck(item)
		return nil
	}

	fn := func(domain Domain) error {
		if err := b.db.Update(item, bind); err != nil {
			return err
		}
		updater := b.e.objectParentUpdater()
		return updater.update(bind, domain)
	}

	return b.e.parentEntryFlowHandler(bind, check, b.e.newObject, fn)
}
