package caskin

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
)

//go:embed configs/casbin_model.conf
var CasbinModelText string

// CasbinModel loads the embedded casbin RBAC model configuration and returns
// a parsed model.Model ready to be passed to a casbin Enforcer.
func CasbinModel() (model.Model, error) {
	return model.NewModelFromString(CasbinModelText)
}

// WatcherOption configures an optional casbin policy watcher that keeps
// multiple enforcer instances in sync when the policy changes.
type WatcherOption struct {
	// Type selects the watcher backend. Currently "redis" is supported.
	// An empty or unrecognised type falls back to auto-load polling when
	// AutoLoad > 0.
	Type string `json:"type"`
	// Address is the host:port of the watcher backend (e.g. "localhost:6379").
	Address string `json:"address"`
	// Password is the authentication password for the watcher backend.
	Password string `json:"password"`
	// Channel is the pub/sub channel name used by the Redis watcher.
	Channel string `json:"channel"`
	// AutoLoad, when > 0 and Type is not "redis", sets the interval in
	// seconds for the enforcer's built-in StartAutoLoadPolicy poller.
	AutoLoad int64 `json:"auto_load"`
}

// SetWatcher attaches a policy watcher to the given casbin enforcer based on
// the provided [WatcherOption]. If option is nil, SetWatcher is a no-op.
//
// Supported watcher types:
//   - "redis": sets up a Redis pub/sub watcher via github.com/casbin/redis-watcher.
//   - "": or unknown type with AutoLoad > 0 — enables the enforcer's built-in
//     periodic policy reload via StartAutoLoadPolicy.
func SetWatcher(e casbin.IEnforcer, option *WatcherOption) error {
	if option == nil {
		return nil
	}
	switch option.Type {
	case "redis":
		wOption := rediswatcher.WatcherOptions{
			Options: redis.Options{
				Network:  "tcp",
				Password: option.Password,
			},
			Channel: option.Channel,
		}
		w, err := rediswatcher.NewWatcher(option.Address, wOption)
		if err != nil {
			return err
		}
		err = e.SetWatcher(w)
		if err != nil {
			return err
		}
		return w.SetUpdateCallback(func(v string) { _ = e.LoadPolicy() })
	default:
		if option.AutoLoad == 0 {
			return nil
		}

		startAutoLoadPolicy, ok := e.(startAutoLoadPolicyInterface)
		if !ok {
			return fmt.Errorf("watcher type %v, enforcer has no StartAutoLoadPolicy method", option.Type)
		}
		startAutoLoadPolicy.StartAutoLoadPolicy(time.Duration(option.AutoLoad) * time.Second)
		return nil
	}
}

type startAutoLoadPolicyInterface interface {
	StartAutoLoadPolicy(time.Duration)
}

// IEnforcer is the caskin-internal interface that wraps casbin's enforcer with
// domain-aware, strongly-typed methods. All parameters and return values use
// the caskin domain model types ([User], [Role], [Object], [Domain], [Action],
// [Policy]) rather than raw strings.
//
// Obtain an IEnforcer via [NewEnforcer].
type IEnforcer interface {
	// Enforce checks whether user u can perform action on object o within domain d.
	Enforce(User, Object, Domain, Action) (bool, error)
	// EnforceRole checks whether son inherits from parent within domain.
	EnforceRole(son Role, parent Role, domain Domain) (bool, error)
	// EnforceObject checks whether son is a descendant of parent within domain.
	EnforceObject(son Object, parent Object, domain Domain) (bool, error)
	// IsSuperadmin returns true if user is a global superadmin.
	IsSuperadmin(User) (bool, error)

	// GetDomainsIncludeUser returns every domain the user belongs to.
	GetDomainsIncludeUser(User) []Domain

	// GetRolesForUserInDomain returns the roles directly assigned to user in domain.
	GetRolesForUserInDomain(User, Domain) []Role
	// GetUsersForRoleInDomain returns the users that have role in domain.
	GetUsersForRoleInDomain(Role, Domain) []User
	// GetParentsForRoleInDomain returns the parent roles of role in domain.
	GetParentsForRoleInDomain(Role, Domain) []Role
	// GetChildrenForRoleInDomain returns the child roles of role in domain.
	GetChildrenForRoleInDomain(Role, Domain) []Role
	// GetParentsForObjectInDomain returns the parent objects of object in domain.
	GetParentsForObjectInDomain(Object, Domain) []Object
	// GetChildrenForObjectInDomain returns the child objects of object in domain.
	GetChildrenForObjectInDomain(Object, Domain) []Object
	// GetPoliciesForRoleInDomain returns all policies granted to role in domain.
	GetPoliciesForRoleInDomain(Role, Domain) []*Policy
	// GetPoliciesForObjectInDomain returns all policies that reference object in domain.
	GetPoliciesForObjectInDomain(Object, Domain) []*Policy

	// RemoveUserInDomain removes all role assignments for user within domain.
	RemoveUserInDomain(User, Domain) error
	// RemoveRoleInDomain removes all policies, inheritance edges, and user
	// assignments associated with role within domain.
	RemoveRoleInDomain(Role, Domain) error
	// RemoveObjectInDomain removes all policies and hierarchy edges associated
	// with object within domain.
	RemoveObjectInDomain(Object, Domain) error

	// AddPolicyInDomain grants role the given action on object within domain.
	AddPolicyInDomain(Role, Object, Domain, Action) error
	// RemovePolicyInDomain revokes role's action on object within domain.
	RemovePolicyInDomain(Role, Object, Domain, Action) error

	// AddRoleForUserInDomain assigns role to user within domain.
	AddRoleForUserInDomain(User, Role, Domain) error
	// RemoveRoleForUserInDomain unassigns role from user within domain.
	RemoveRoleForUserInDomain(User, Role, Domain) error

	// AddParentForRoleInDomain makes son inherit from parent within domain.
	AddParentForRoleInDomain(Role, Role, Domain) error
	// RemoveParentForRoleInDomain removes the son→parent inheritance in domain.
	RemoveParentForRoleInDomain(Role, Role, Domain) error

	// AddParentForObjectInDomain makes son a child of parent within domain.
	AddParentForObjectInDomain(Object, Object, Domain) error
	// RemoveParentForObjectInDomain removes the son→parent object edge in domain.
	RemoveParentForObjectInDomain(Object, Object, Domain) error

	// GetUsersInDomain returns all users that have any role in domain.
	GetUsersInDomain(Domain) []User
	// GetRolesInDomain returns all roles defined in domain.
	GetRolesInDomain(Domain) []Role
	// GetObjectsInDomain returns all objects defined in domain.
	GetObjectsInDomain(Domain) []Object
	// GetPoliciesInDomain returns all policies defined in domain.
	GetPoliciesInDomain(Domain) []*Policy

	// RemoveUsersInDomain removes all user→role assignments within domain,
	// effectively clearing the domain's user membership.
	RemoveUsersInDomain(Domain) error
}

type enforcer struct {
	e       casbin.IEnforcer
	factory Factory
}

func (e *enforcer) Enforce(user User, object Object, domain Domain, action Action) (bool, error) {
	return e.e.Enforce(user.Encode(), domain.Encode(), object.Encode(), action)
}

func (e *enforcer) EnforceRole(son Role, parent Role, domain Domain) (bool, error) {
	rs, err := e.e.GetImplicitRolesForUser(son.Encode(), domain.Encode())
	if err != nil {
		return false, err
	}
	for _, r := range rs {
		if r == parent.Encode() {
			return true, nil
		}
	}
	return false, nil
}

func (e *enforcer) EnforceObject(son Object, parent Object, domain Domain) (bool, error) {
	os, err := e.e.GetImplicitRolesForUser(parent.Encode(), domain.Encode())
	if err != nil {
		return false, err
	}
	for _, o := range os {
		if o == son.Encode() {
			return true, nil
		}
	}
	return false, nil
}

func (e *enforcer) IsSuperadmin(user User) (bool, error) {
	return e.e.HasRoleForUser(user.Encode(), SuperadminRole, SuperadminDomain)
}

func (e *enforcer) GetDomainsIncludeUser(user User) []Domain {
	var domains []Domain
	rules, _ := e.e.GetFilteredGroupingPolicy(0, user.Encode())
	for _, rule := range rules {
		if domain, err := e.factory.Domain(rule[2]); err == nil {
			domains = append(domains, domain)
		}
	}
	return domains
}

func (e *enforcer) GetRolesForUserInDomain(user User, domain Domain) []Role {
	var roles []Role
	rs := e.e.GetRolesForUserInDomain(user.Encode(), domain.Encode())
	for _, r := range rs {
		if role, err := e.factory.Role(r); err == nil {
			roles = append(roles, role)
		}
	}
	return roles
}

func (e *enforcer) GetUsersForRoleInDomain(role Role, domain Domain) []User {
	var users []User
	us := e.e.GetUsersForRoleInDomain(role.Encode(), domain.Encode())
	for _, u := range us {
		if user, err := e.factory.User(u); err == nil {
			users = append(users, user)
		}
	}
	return users
}

func (e *enforcer) GetParentsForRoleInDomain(role Role, domain Domain) []Role {
	var roles []Role
	rs := e.e.GetUsersForRoleInDomain(role.Encode(), domain.Encode())
	for _, v := range rs {
		if _, err := e.factory.Role(v); err == nil {
			roles = append(roles, role)
		}
	}

	return roles
}

func (e *enforcer) GetChildrenForRoleInDomain(role Role, domain Domain) []Role {
	var roles []Role
	rs := e.e.GetRolesForUserInDomain(role.Encode(), domain.Encode())
	for _, v := range rs {
		if r, err := e.factory.Role(v); err == nil {
			roles = append(roles, r)
		}
	}

	return roles
}

func (e *enforcer) GetParentsForObjectInDomain(object Object, domain Domain) []Object {
	var objects []Object
	os, _ := e.e.GetModel()["g"][ObjectPType].RM.GetRoles(object.Encode(), domain.Encode())
	for _, v := range os {
		if o, err := e.factory.Object(v); err == nil {
			objects = append(objects, o)
		}
	}

	return objects
}

func (e *enforcer) GetChildrenForObjectInDomain(object Object, domain Domain) []Object {
	var objects []Object
	os, _ := e.e.GetModel()["g"][ObjectPType].RM.GetUsers(object.Encode(), domain.Encode())
	for _, v := range os {
		if o, err := e.factory.Object(v); err == nil {
			objects = append(objects, o)
		}
	}

	return objects
}

func (e *enforcer) GetPoliciesForRoleInDomain(role Role, domain Domain) []*Policy {
	var policies []*Policy
	ps, _ := e.e.GetPermissionsForUser(role.Encode(), domain.Encode())
	for _, p := range ps {
		r, err1 := e.factory.Role(p[0])
		o, err2 := e.factory.Object(p[2])
		if err1 != nil || err2 != nil {
			continue
		}
		action := p[3]
		pp := &Policy{
			Role:   r,
			Object: o,
			Domain: domain,
			Action: action,
		}
		policies = append(policies, pp)
	}

	return policies
}

func (e *enforcer) GetPoliciesForObjectInDomain(object Object, domain Domain) []*Policy {
	var policies []*Policy
	ps, _ := e.e.GetFilteredPolicy(1, domain.Encode(), object.Encode())
	for _, p := range ps {
		r, err1 := e.factory.Role(p[0])
		o, err2 := e.factory.Object(p[2])
		if err1 != nil || err2 != nil {
			continue
		}
		action := p[3]
		pp := &Policy{
			Role:   r,
			Object: o,
			Domain: domain,
			Action: action,
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
	if _, err := e.e.RemoveFilteredPolicy(0, role.Encode(), domain.Encode()); err != nil {
		return err
	}

	cs := e.GetChildrenForRoleInDomain(role, domain)
	for _, v := range cs {
		if err := e.RemoveParentForRoleInDomain(v, role, domain); err != nil {
			return err
		}
	}

	ps := e.GetParentsForRoleInDomain(role, domain)
	for _, v := range ps {
		if err := e.RemoveParentForRoleInDomain(role, v, domain); err != nil {
			return err
		}
	}

	us := e.GetUsersForRoleInDomain(role, domain)
	for _, v := range us {
		if err := e.RemoveRoleForUserInDomain(v, role, domain); err != nil {
			return err
		}
	}

	return nil
}

func (e *enforcer) RemoveObjectInDomain(object Object, domain Domain) error {
	if _, err := e.e.RemoveFilteredPolicy(1, domain.Encode(), object.Encode()); err != nil {
		return err
	}

	cs := e.GetChildrenForObjectInDomain(object, domain)
	for _, v := range cs {
		if err := e.RemoveParentForObjectInDomain(v, object, domain); err != nil {
			return err
		}
	}

	ps := e.GetParentsForObjectInDomain(object, domain)
	for _, v := range ps {
		if err := e.RemoveParentForObjectInDomain(object, v, domain); err != nil {
			return err
		}
	}

	return nil
}

func (e *enforcer) AddPolicyInDomain(role Role, object Object, domain Domain, action Action) error {
	_, err := e.e.AddPolicy(role.Encode(), domain.Encode(), object.Encode(), action)
	return err
}

func (e *enforcer) RemovePolicyInDomain(role Role, object Object, domain Domain, action Action) error {
	_, err := e.e.RemovePolicy(role.Encode(), domain.Encode(), object.Encode(), action)
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
	_, err := e.e.AddRoleForUserInDomain(son.Encode(), parent.Encode(), domain.Encode())
	return err
}

func (e *enforcer) RemoveParentForRoleInDomain(son Role, parent Role, domain Domain) error {
	_, err := e.e.DeleteRoleForUserInDomain(son.Encode(), parent.Encode(), domain.Encode())
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
	rules, _ := e.e.GetFilteredGroupingPolicy(2, domain.Encode())
	for _, rule := range rules {
		if user, err := e.factory.User(rule[0]); err == nil {
			users = append(users, user)
		}
	}

	return users
}

func (e *enforcer) GetRolesInDomain(domain Domain) []Role {
	var roles []Role
	rules, _ := e.e.GetFilteredGroupingPolicy(2, domain.Encode())
	for _, rule := range rules {
		r1, err1 := e.factory.Role(rule[0])
		if err1 != nil {
			continue
		}
		roles = append(roles, r1)
		r2, err2 := e.factory.Role(rule[1])
		if err2 != nil {
			continue
		}
		roles = append(roles, r2)
	}

	return roles
}

func (e *enforcer) GetObjectsInDomain(domain Domain) []Object {
	var objects []Object
	rules, _ := e.e.GetFilteredNamedGroupingPolicy(ObjectPType, 2, domain.Encode())
	for _, rule := range rules {
		o1, err1 := e.factory.Object(rule[0])
		if err1 != nil {
			continue
		}
		objects = append(objects, o1)
		o2, err2 := e.factory.Object(rule[1])
		if err2 != nil {
			continue
		}
		o1.SetParentID(o2.GetID())
		objects = append(objects, o2)
	}

	return objects
}

func (e *enforcer) GetPoliciesInDomain(domain Domain) []*Policy {
	var policies []*Policy
	ps, _ := e.e.GetFilteredPolicy(1, domain.Encode())
	for _, p := range ps {
		r, err1 := e.factory.Role(p[0])
		o, err2 := e.factory.Object(p[2])
		if err1 != nil || err2 != nil {
			continue
		}
		action := p[3]

		pp := &Policy{
			Role:   r,
			Object: o,
			Domain: domain,
			Action: action,
		}
		policies = append(policies, pp)
	}

	return policies
}

func (e *enforcer) RemoveUsersInDomain(domain Domain) error {
	gp, _ := e.e.GetFilteredGroupingPolicy(2, domain.Encode())
	var rules [][]string
	for _, rule := range gp {
		if _, err := e.factory.User(rule[0]); err == nil {
			rules = append(rules, rule)
		}
	}

	_, err := e.e.RemoveGroupingPolicies(rules)
	return err
}

// NewEnforcer wraps a casbin enforcer and a [Factory] into an [IEnforcer].
// The factory is used to decode casbin string tokens back into typed caskin values.
func NewEnforcer(e casbin.IEnforcer, factory Factory) IEnforcer {
	return &enforcer{
		e:       e,
		factory: factory,
	}
}
