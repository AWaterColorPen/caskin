package caskin

import "github.com/ahmetb/go-linq/v3"

// CreateRole
// if current user has role's write permission and there does not exist the role
// then create a new one
// 1. create a new role into metadata database
// 2. Run role to parent's g in the domain
func (e *Executor) CreateRole(role Role) error {
	if err := e.ObjectDataCreateCheck(role, ObjectTypeRole); err != nil {
		return err
	}
	if err := e.TreeNodeEntryParentCheck(role, e.newRole()); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	role.SetDomainID(domain.GetID())
	if err := e.DB.Create(role); err != nil {
		return err
	}
	updater := e.DefaultRoleUpdater()
	return updater.Run(role, domain)
}

// RecoverRole
// if current user has role's write permission and there exist the role but soft deleted
// then recover it
// 1. recover the soft delete one role at metadata database
// 2. Run role to parent's g in the domain
func (e *Executor) RecoverRole(role Role) error {
	if err := e.ObjectDataRecoverCheck(role); err != nil {
		return err
	}
	if err := e.TreeNodeEntryParentCheck(role, e.newRole()); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	role.SetDomainID(domain.GetID())
	if err := e.DB.Recover(role); err != nil {
		return err
	}
	updater := e.DefaultRoleUpdater()
	return updater.Run(role, domain)
}

// DeleteRole
// if current user has role's write permission
// 1. delete role's g in the domain
// 2. delete role's p in the domain
// 3. soft delete one role in metadata database
// 4. Run to delete all son of the role in the domain
func (e *Executor) DeleteRole(role Role) error {
	if err := e.ObjectDataDeleteCheck(role); err != nil {
		return err
	}
	if err := e.TreeNodeEntryParentCheck(role, e.newRole()); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	role.SetDomainID(domain.GetID())
	deleter := NewTreeNodeEntryDeleter(e.DefaultRoleChildrenGetFunc(), e.DefaultRoleDeleteFunc())
	return deleter.Run(role, domain)
}

// UpdateRole
// if current user has role's write permission and there exist the role
// 1. Run role's properties
// 2. Run role to parent's g in the domain
func (e *Executor) UpdateRole(role Role) error {
	tmp1, tmp2 := e.newRole(), e.newRole()
	if err := e.TreeNodeEntryUpdateCheck(role, tmp1, tmp2, ObjectTypeRole); err != nil {
		return err
	}
	if err := e.TreeNodeEntryParentCheck(role, e.newRole()); err != nil {
		return err
	}
	_, domain, _ := e.provider.Get()
	role.SetDomainID(domain.GetID())
	if err := e.DB.Update(role); err != nil {
		return err
	}
	updater := e.DefaultRoleUpdater()
	return updater.Run(role, domain)
}

// GetRoles
// if current user has role's read permission
// 1. get roles in the domain
func (e *Executor) GetRoles() ([]Role, error) {
	currentUser, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	roles, err := e.DB.GetRoleInDomain(currentDomain)
	if err != nil {
		return nil, err
	}
	out := e.filterWithNoError(currentUser, currentDomain, Read, roles)
	linq.From(out).ToSlice(&roles)
	return roles, nil
}
