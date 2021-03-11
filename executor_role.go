package caskin

// CreateRole
// if current user has role's write permission and there does not exist the role
// then create a new one
// 1. create a new role into metadata database
// 2. update role to parent's g in the domain
func (e *executor) CreateRole(role Role) error {
	fn := func(interface{}) error {
		if err := e.mdb.Create(role); err != nil {
			return err
		}

		_, domain, err := e.provider.Get()
		if err != nil {
			return err
		}

		if role.GetParentID() != 0 {
			parent := e.factory.NewRole()
			parent.SetID(role.GetParentID())
			if err := e.e.AddParentForRoleInDomain(role, parent, domain); err != nil {
				return err
			}
		}
		return nil
	}
	return e.createOrRecoverRole(role, fn)
}

// RecoverRole
// if current user has role's write permission and there exist the role but soft deleted
// then recover it
// 1. recover the soft delete one role at metadata database
// 2. update role to parent's g in the domain
func (e *executor) RecoverRole(role Role) error {
	return e.createOrRecoverRole(role, e.mdb.Recover)
}

// DeleteRole
// if current user has role's write permission
// 1. delete role's g in the domain
// 2. delete role's p in the domain
// 3. soft delete one role in metadata database
func (e *executor) DeleteRole(role Role) error {
	fn := func(interface{}) error {
		_, domain, err := e.provider.Get()
		if err != nil {
			return err
		}
		if err := e.e.RemoveRoleInDomain(role, domain); err != nil {
			return err
		}
		return e.mdb.DeleteRoleByID(role.GetID())
	}

	return e.writeRole(role, fn)
}

// UpdateRole
// if current user has role's write permission and there exist the role
// 1. update role's properties
// 2. update role to parent's g in the domain
func (e *executor) UpdateRole(role Role) error {
	return e.writeRole(role, e.mdb.Update)
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
	roles, err := e.mdb.GetRoleInDomain(currentDomain)
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

func (e *executor) createOrRecoverRole(role Role, fn func(interface{}) error) error {
	if err := e.mdb.Take(role); err == nil {
		return ErrAlreadyExists
	}

	_, domain, err := e.provider.Get()
	if err != nil {
		return err
	}

	if err := e.check(role, Write); err != nil {
		return err
	}

	var parents []Role
	if role.GetParentID() != 0 {
		nextParent := e.factory.NewRole()
		nextParent.SetID(role.GetParentID())
		parents = append(parents, nextParent)
	}
	for _, v := range parents {
		v.SetDomainID(domain.GetID())
		if err := e.mdb.Take(v); err != nil {
			return err
		}
		if err := e.check(v, Write); err != nil {
			return err
		}
	}

	role.SetDomainID(domain.GetID())
	return fn(role)
}

func (e *executor) writeRole(role Role, fn func(interface{}) error) error {
	if err := isValid(role); err != nil {
		return err
	}

	tmp := e.factory.NewRole()
	tmp.SetID(role.GetID())
	if err := e.mdb.Take(tmp); err != nil {
		return ErrNotExists
	}

	if err := e.check(tmp, Write); err != nil {
		return err
	}

	_, domain, err := e.provider.Get()
	if err != nil {
		return err
	}

	parents := e.e.GetParentsForRoleInDomain(tmp, domain)

	nextParent := e.factory.NewRole()
	nextParent.SetID(role.GetParentID())
	parents = append(parents, nextParent)
	for _, v := range parents {
		v.SetDomainID(domain.GetID())
		if err := e.mdb.Take(v); err != nil {
			return err
		}
		if err := e.check(v, Write); err != nil {
			return err
		}
	}

	role.SetDomainID(domain.GetID())
	return fn(role)
}
