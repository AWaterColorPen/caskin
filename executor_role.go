package caskin

// CreateRole
// if current user has role's write permission and there does not exist the role
// then create a new one
// 1. create a new role into metadata database
// 2. update role to parent's g in the domain
func (e *executor) CreateRole(role Role) error {
	fn := func(domain Domain) error {
		if err := e.db.Create(role); err != nil {
			return err
		}
		updater := e.roleParentUpdater()
		return updater.update(role, domain)
	}

	return e.parentEntryFlowHandler(role, e.createObjectDataEntryCheck, e.newRole, fn)
}

// RecoverRole
// if current user has role's write permission and there exist the role but soft deleted
// then recover it
// 1. recover the soft delete one role at metadata database
// 2. update role to parent's g in the domain
func (e *executor) RecoverRole(role Role) error {
	fn := func(domain Domain) error {
		if err := e.db.Recover(role); err != nil {
			return err
		}
		updater := e.roleParentUpdater()
		return updater.update(role, domain)
	}

	return e.parentEntryFlowHandler(role, e.recoverObjectDataEntryCheck, e.newRole, fn)
}

// DeleteRole
// if current user has role's write permission
// 1. delete role's g in the domain
// 2. delete role's p in the domain
// 3. soft delete one role in metadata database
// 4. dfs to delete all son of the role in the domain
func (e *executor) DeleteRole(role Role) error {
	fn := func(domain Domain) error {
		deleter := newParentEntryDeleter(e.roleChildrenFn(), e.roleDeleteFn())
		return deleter.dfs(role, domain)
	}

	return e.parentEntryFlowHandler(role, e.deleteObjectDataEntryCheck, e.newRole, fn)
}

// UpdateRole
// if current user has role's write permission and there exist the role
// 1. update role's properties
// 2. update role to parent's g in the domain
func (e *executor) UpdateRole(role Role) error {
	fn := func(domain Domain) error {
		if err := e.db.Update(role); err != nil {
			return err
		}
		updater := e.roleParentUpdater()
		return updater.update(role, domain)
	}

	roleUpdateCheck := func(item objectDataEntry) error {
		tmp := e.newRole()
		if err := e.updateObjectDataEntryCheck(item, tmp); err != nil {
			return err
		}
		return e.treeNodeParentCheck(tmp, e.newRole)
	}
	return e.parentEntryFlowHandler(role, roleUpdateCheck, e.newRole, fn)
}

// GetRoles
// if current user has role's read permission
// 1. get roles in the domain
func (e *executor) GetRoles() (Roles, error) {
	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	rs := e.e.GetRolesInDomain(currentDomain)
	tree := getTree(rs)
	roles, err := e.db.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	r := e.filterWithNoError(currentUser, currentDomain, Read, roles)
	roles = []Role{}
	for _, v := range r {
		roles = append(roles, v.(Role))
	}

	for _, v := range roles {
		if p, ok := tree[v.GetID()]; ok {
			v.SetParentID(p)
		}
	}

	return roles, nil
}
