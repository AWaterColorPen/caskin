package caskin

// RoleCreate
// if there does not exist the role then create a new one
// 1. current user has role's write permission
// 2. create a new role into metadata database
// 3. set role to parent's g in the domain
func (e *server) RoleCreate(user User, domain Domain, role Role) error {
	if err := e.ObjectDataCreateCheck(user, domain, role, ObjectTypeRole); err != nil {
		return err
	}
	if err := e.RoleParentCheck(role); err != nil {
		return err
	}
	role.SetDomainID(domain.GetID())
	if err := e.DB.Create(role); err != nil {
		return err
	}
	updater := e.DefaultRoleUpdater()
	return updater.Run(role, domain)
}

// RoleRecover
// if there exist the role but soft deleted then recover it
// 1. current user has role's write permission
// 2. recover the soft delete one role at metadata database
// 3. set role to parent's g in the domain
func (e *server) RoleRecover(user User, domain Domain, role Role) error {
	if err := e.ObjectDataRecoverCheck(user, domain, role); err != nil {
		return err
	}
	if err := e.RoleParentCheck(role); err != nil {
		return err
	}
	role.SetDomainID(domain.GetID())
	if err := e.DB.Recover(role); err != nil {
		return err
	}
	updater := e.DefaultRoleUpdater()
	return updater.Run(role, domain)
}

// RoleDelete
// if there exist the object
// 1. current user has role's write permission
// 1. delete role's g in the domain
// 2. delete role's p in the domain
// 3. soft delete one role in metadata database
// 4. delete all son of the role in the domain
func (e *server) RoleDelete(user User, domain Domain, role Role) error {
	if err := e.ObjectDataDeleteCheck(user, domain, role); err != nil {
		return err
	}
	if err := e.RoleParentCheck(role); err != nil {
		return err
	}
	role.SetDomainID(domain.GetID())
	deleter := NewTreeNodeDeleter(e.DefaultRoleChildrenGetFunc(), e.DefaultRoleDeleteFunc())
	return deleter.Run(role, domain)
}

// RoleUpdate
// if there exist the role
// 1. current user has role's write permission and
// 2. update role's properties
// 3. update role to parent's g in the domain
func (e *server) RoleUpdate(user User, domain Domain, role Role) error {
	if err := e.RoleUpdateCheck(user, domain, role); err != nil {
		return err
	}
	if err := e.RoleParentCheck(role); err != nil {
		return err
	}
	role.SetDomainID(domain.GetID())
	if err := e.DB.Update(role); err != nil {
		return err
	}
	updater := e.DefaultRoleUpdater()
	return updater.Run(role, domain)
}

// RoleGet
// get role
// 1. current user has role's read permission
func (e *server) RoleGet(user User, domain Domain) ([]Role, error) {
	roles, err := e.DB.GetRoleInDomain(domain)
	if err != nil {
		return nil, err
	}
	return Filter(e.Enforcer, user, domain, Read, roles), nil
}
