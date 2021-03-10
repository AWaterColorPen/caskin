package caskin

import (
	"github.com/casbin/casbin/v2"
)

type ienforcer interface {
	// check permission
	Enforce(User, Object, Domain, Action) (bool, error)
	IsSuperAdmin(User) (bool, error)

	// get grouping entry in domain
	GetRolesForUserInDomain(User, Domain) []Role
	GetUsersForRoleInDomain(Role, Domain) []User
	GetParentsForRoleInDomain(Role, Domain) []Role
	GetChildrenForRoleInDomain(Role, Domain) []Role
	GetParentsForObjectInDomain(Object, Domain) []Object
	GetChildrenForObjectInDomain(Object, Domain) []Object
	GetPoliciesForRoleInDomain(Role, Domain) []*Policy

	// remove entry in domain
	RemoveUserInDomain(User, Domain) error
	RemoveRoleInDomain(Role, Domain) error
	RemoveObjectInDomain(Object, Domain) error

	// add or remove policy information
	AddPolicyInDomain(Role, Object, Domain, Action) error
	RemovePolicyInDomain(Role, Object, Domain, Action) error

	// add or remove user-role grouping information
	AddRoleForUserInDomain(User, Role, Domain) error
	RemoveRoleForUserInDomain(User, Role, Domain) error

	// add or remove role-parent grouping information
	AddParentForRoleInDomain(Role, Role, Domain) error
	RemoveParentForRoleInDomain(Role, Role, Domain) error

	// add or remove object-parent grouping information
	AddParentForObjectInDomain(Object, Object, Domain) error
	RemoveParentForObjectInDomain(Object, Object, Domain) error

	// get all entry in domain
	GetUsersInDomain(Domain) []User
	GetRolesInDomain(Domain) []Role
	GetObjectsInDomain(Domain) []Object
	GetPoliciesInDomain(Domain) []*Policy

	// remove entry in domain
	RemoveUsersInDomain(Domain) error
}

type enforcer struct {
	e       casbin.IEnforcer
	factory EntryFactory
}

func (e *enforcer) Enforce(user User, object Object, domain Domain, action Action) (bool, error) {
	return e.e.Enforce(user.Encode(), domain.Encode(), object.Encode(), string(action))
}

func (e *enforcer) IsSuperAdmin(user User) (bool, error) {
	return e.e.HasRoleForUser(user.Encode(), SuperadminRole, SuperadminDomain)
}

func (e *enforcer) GetRolesForUserInDomain(user User, domain Domain) []Role {
	var roles []Role
	rs := e.e.GetRolesForUserInDomain(user.Encode(), domain.Encode())
	for _, r := range rs {
		role := e.factory.NewRole()
		if err := role.Decode(r); err == nil {
			roles = append(roles, role)
		}
	}

	return roles
}

func (e *enforcer) GetUsersForRoleInDomain(role Role, domain Domain) []User {
	var users []User
	us := e.e.GetUsersForRoleInDomain(role.Encode(), domain.Encode())
	for _, u := range us {
		user := e.factory.NewUser()
		if err := user.Decode(u); err == nil {
			users = append(users, user)
		}
	}

	return users
}

func (e *enforcer) GetParentsForRoleInDomain(role Role, domain Domain) []Role {
	var roles []Role
	rs := e.e.GetRolesForUserInDomain(role.Encode(), domain.Encode())
	for _, v := range rs {
		r := e.factory.NewRole()
		if err := r.Decode(v); err == nil {
			roles = append(roles, role)
		}
	}

	return roles
}

func (e *enforcer) GetChildrenForRoleInDomain(role Role, domain Domain) []Role {
	var roles []Role
	rs := e.e.GetUsersForRoleInDomain(role.Encode(), domain.Encode())
	for _, v := range rs {
		r := e.factory.NewRole()
		if err := r.Decode(v); err == nil {
			roles = append(roles, role)
		}
	}

	return roles
}

func (e *enforcer) GetParentsForObjectInDomain(object Object, domain Domain) []Object {
	var objects []Object
	os, _ := e.e.GetModel()["g"][ObjectPType].RM.GetRoles(object.Encode(), domain.Encode())
	for _, v := range os {
		o := e.factory.NewObject()
		if err := o.Decode(v); err == nil {
			objects = append(objects, o)
		}
	}

	return objects
}

func (e *enforcer) GetChildrenForObjectInDomain(object Object, domain Domain) []Object {
	var objects []Object
	os, _ := e.e.GetModel()["g"][ObjectPType].RM.GetUsers(object.Encode(), domain.Encode())
	for _, v := range os {
		o := e.factory.NewObject()
		if err := o.Decode(v); err == nil {
			objects = append(objects, o)
		}
	}

	return objects
}

func (e *enforcer) GetPoliciesForRoleInDomain(role Role, domain Domain) []*Policy {
	var policies []*Policy
	ps := e.e.GetPermissionsForUser(role.Encode(), domain.Encode())
	for _, p := range ps {
		r := e.factory.NewRole()
		o := e.factory.NewObject()
		if err := r.Decode(p[0]); err != nil {
			continue
		}
		if err := o.Decode(p[2]); err != nil {
			continue
		}
		action := p[3]

		pp := &Policy{
			Role:   r,
			Object: o,
			Domain: domain,
			Action: Action(action),
		}
		policies = append(policies, pp)
	}

	return policies
}

func (e *enforcer) RemoveUserInDomain(user User, domain Domain) error {
	roles := e.GetRolesForUserInDomain(user, domain)
	for _, role := range roles {
		if err := e.RemoveRoleForUserInDomain(user, role, domain); err == nil {
			return err
		}
	}

	return nil
}

func (e *enforcer) RemoveRoleInDomain(role Role, domain Domain) error {
	us := e.GetUsersForRoleInDomain(role, domain)
	for _, v := range us {
		if err := e.RemoveRoleForUserInDomain(v, role, domain); err != nil {
			return err
		}
	}

	ps := e.GetParentsForRoleInDomain(role, domain)
	for _, v := range ps {
		if err := e.RemoveRoleForUserInDomain(v, role, domain); err != nil {
			return err
		}
	}

	cs := e.GetChildrenForRoleInDomain(role, domain)
	for _, v := range cs {
		if err := e.RemoveRoleForUserInDomain(role, v, domain); err != nil {
			return err
		}
	}

	_, err := e.e.RemoveFilteredPolicy(0, role.Encode(), domain.Encode())
	return err
}

func (e *enforcer) RemoveObjectInDomain(object Object, domain Domain) error {
	ps := e.GetParentsForObjectInDomain(object, domain)
	for _, v := range ps {
		if err := e.RemoveParentForObjectInDomain(object, v, domain); err != nil {
			return err
		}
	}

	cs := e.GetChildrenForObjectInDomain(object, domain)
	for _, v := range cs {
		if err := e.RemoveParentForObjectInDomain(v, object, domain); err != nil {
			return err
		}
	}

	_, err := e.e.RemoveFilteredPolicy(1, domain.Encode(), object.Encode())
	return err
}

func (e *enforcer) AddPolicyInDomain(role Role, object Object, domain Domain, action Action) error {
	_, err := e.e.AddPolicy(role.Encode(), domain.Encode(), object.Encode(), string(action))
	return err
}

func (e *enforcer) RemovePolicyInDomain(role Role, object Object, domain Domain, action Action) error {
	_, err := e.e.RemovePolicy(role.Encode(), domain.Encode(), object.Encode(), string(action))
	return err
}

func (e *enforcer) AddRoleForUserInDomain(user User, role Role, domain Domain) error {
	_, err := e.e.AddRoleForUserInDomain(user.Encode(), role.Encode(), domain.Encode())
	return err
}

func (e *enforcer) RemoveRoleForUserInDomain(user User, role Role, domain Domain) error {
	_, err := e.e.DeleteRoleForUserInDomain(user.Encode(), role.Encode(), domain.Encode())
	return err
}

func (e *enforcer) AddParentForRoleInDomain(son Role, parent Role, domain Domain) error {
	_, err := e.e.AddRoleForUserInDomain(parent.Encode(), son.Encode(), domain.Encode())
	return err
}

func (e *enforcer) RemoveParentForRoleInDomain(son Role, parent Role, domain Domain) error {
	_, err := e.e.DeleteRoleForUserInDomain(parent.Encode(), son.Encode(), domain.Encode())
	return err
}

func (e *enforcer) AddParentForObjectInDomain(son Object, parent Object, domain Domain) error {
	_, err := e.e.AddNamedGroupingPolicy(ObjectPType, son.Encode(), parent.Encode(), domain.Encode())
	return err
}

func (e *enforcer) RemoveParentForObjectInDomain(son Object, parent Object, domain Domain) error {
	_, err := e.e.RemoveNamedGroupingPolicy(ObjectPType, son.Encode(), parent.Encode(), domain.Encode())
	return err
}

func (e *enforcer) GetUsersInDomain(domain Domain) []User {
	var users []User
	rules := e.e.GetFilteredGroupingPolicy(2, domain.Encode())
	for _, rule := range rules {
		user := e.factory.NewUser()
		if err := user.Decode(rule[0]); err == nil {
			users = append(users, user)
		}
	}

	return users
}

func (e *enforcer) GetRolesInDomain(domain Domain) []Role {
	var roles []Role
	rules := e.e.GetFilteredGroupingPolicy(2, domain.Encode())
	for _, rule := range rules {
		r1 := e.factory.NewRole()
		if err := r1.Decode(rule[0]); err != nil {
			continue
		}
		roles = append(roles, r1)
		r2 := e.factory.NewRole()
		if err := r2.Decode(rule[1]); err != nil {
			continue
		}
		r2.SetParentID(r1.GetID())
		roles = append(roles, r2)
	}

	return roles
}

func (e *enforcer) GetObjectsInDomain(domain Domain) []Object {
	var objects []Object
	rules := e.e.GetFilteredNamedGroupingPolicy(ObjectPType, 2, domain.Encode())
	for _, rule := range rules {
		o1 := e.factory.NewObject()
		if err := o1.Decode(rule[0]); err != nil {
			continue
		}
		objects = append(objects, o1)
		o2 := e.factory.NewObject()
		if err := o2.Decode(rule[1]); err != nil {
			continue
		}
		o1.SetParentID(o2.GetID())
		objects = append(objects, o2)
	}

	return objects
}

func (e *enforcer) GetPoliciesInDomain(domain Domain) []*Policy {
	var policies []*Policy
	ps := e.e.GetFilteredPolicy(1, domain.Encode())
	for _, p := range ps {
		r := e.factory.NewRole()
		o := e.factory.NewObject()
		if err := r.Decode(p[0]); err != nil {
			continue
		}
		if err := o.Decode(p[2]); err != nil {
			continue
		}
		action := p[3]

		pp := &Policy{
			Role:   r,
			Object: o,
			Domain: domain,
			Action: Action(action),
		}
		policies = append(policies, pp)
	}

	return policies
}

func (e *enforcer) RemoveUsersInDomain(domain Domain) error {
	user := e.factory.NewUser()
	gp := e.e.GetFilteredGroupingPolicy(2, domain.Encode())
	var rules [][]string
	for _, rule := range gp {
		if err := user.Decode(rule[0]); err == nil {
			rules = append(rules, rule)
		}
	}

	_, err := e.e.RemoveGroupingPolicies(rules)
	return err
}

func NewEnforcer(e casbin.IEnforcer, factory EntryFactory) ienforcer {
	return &enforcer{
		e:       e,
		factory: factory,
	}
}
