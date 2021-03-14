package caskin

// CreateRole
// if current user has role's write permission and there does not exist the role
// then create a new one
// 1. create a new role into metadata database
// 2. update role to parent's g in the domain
func (e *Executor) CreateRole(role Role) error {
	if err := e.treeNodeEntryCheckFlow(role, e.ObjectDataCreateCheck, e.newRole); err != nil {
		return err
	}
	if err := e.DB.Create(role); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	updater := e.roleParentUpdater()
	return updater.update(role, domain)
}

// RecoverRole
// if current user has role's write permission and there exist the role but soft deleted
// then recover it
// 1. recover the soft delete one role at metadata database
// 2. update role to parent's g in the domain
func (e *Executor) RecoverRole(role Role) error {
	if err := e.treeNodeEntryCheckFlow(role, e.ObjectDataRecoverCheck, e.newRole); err != nil {
		return err
	}
	if err := e.DB.Recover(role); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	updater := e.roleParentUpdater()
	return updater.update(role, domain)
}

// DeleteRole
// if current user has role's write permission
// 1. delete role's g in the domain
// 2. delete role's p in the domain
// 3. soft delete one role in metadata database
// 4. dfs to delete all son of the role in the domain
func (e *Executor) DeleteRole(role Role) error {
	if err := e.treeNodeEntryCheckFlow(role, e.ObjectDataDeleteCheck, e.newRole); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	deleter := newParentEntryDeleter(e.roleChildrenFn(), e.roleDeleteFn())
	return deleter.dfs(role, domain)
}

// UpdateRole
// if current user has role's write permission and there exist the role
// 1. update role's properties
// 2. update role to parent's g in the domain
func (e *Executor) UpdateRole(role Role) error {
	tmp := e.newRole()
	if err := e.ObjectDataUpdateCheck(role, tmp); err != nil {
		return err
	}
	if err := e.treeNodeParentCheck(tmp, e.newRole); err != nil {
		return err
	}
	if err := e.treeNodeEntryCheckFlow(role, func(ObjectData) error { return nil }, e.newRole); err != nil {
		return err
	}
	if err := e.DB.Update(role); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	updater := e.roleParentUpdater()
	return updater.update(role, domain)
}

// GetRoles
// if current user has role's read permission
// 1. get roles in the domain
func (e *Executor) GetRoles() (Roles, error) {
	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	rs := e.e.GetRolesInDomain(currentDomain)
	tree := getTree(rs)
	roles, err := e.DB.GetRoleInDomain(currentDomain)
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
