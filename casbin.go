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

func CasbinModel() (model.Model, error) {
	return model.NewModelFromString(CasbinModelText)
}

type WatcherOption struct {
	Type     string `json:"type"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Channel  string `json:"channel"`
	AutoLoad int64  `json:"auto_load"`
}

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
		startAutoLoadPolicy.StartAutoLoadPolicy(time.Duration(option.AutoLoad))
		return nil
	}
}

type startAutoLoadPolicyInterface interface {
	StartAutoLoadPolicy(time.Duration)
}

type IEnforcer interface {
	Enforce(User, Object, Domain, Action) (bool, error)
	EnforceRole(son Role, parent Role, domain Domain) (bool, error)
	EnforceObject(son Object, parent Object, domain Domain) (bool, error)
	IsSuperadmin(User) (bool, error)

	GetDomainsIncludeUser(User) []Domain

	GetRolesForUserInDomain(User, Domain) []Role
	GetUsersForRoleInDomain(Role, Domain) []User
	GetParentsForRoleInDomain(Role, Domain) []Role
	GetChildrenForRoleInDomain(Role, Domain) []Role
	GetParentsForObjectInDomain(Object, Domain) []Object
	GetChildrenForObjectInDomain(Object, Domain) []Object
	GetPoliciesForRoleInDomain(Role, Domain) []*Policy
	GetPoliciesForObjectInDomain(Object, Domain) []*Policy

	RemoveUserInDomain(User, Domain) error
	RemoveRoleInDomain(Role, Domain) error
	RemoveObjectInDomain(Object, Domain) error

	AddPolicyInDomain(Role, Object, Domain, Action) error
	RemovePolicyInDomain(Role, Object, Domain, Action) error

	AddRoleForUserInDomain(User, Role, Domain) error
	RemoveRoleForUserInDomain(User, Role, Domain) error

	AddParentForRoleInDomain(Role, Role, Domain) error
	RemoveParentForRoleInDomain(Role, Role, Domain) error

	AddParentForObjectInDomain(Object, Object, Domain) error
	RemoveParentForObjectInDomain(Object, Object, Domain) error

	GetUsersInDomain(Domain) []User
	GetRolesInDomain(Domain) []Role
	GetObjectsInDomain(Domain) []Object
	GetPoliciesInDomain(Domain) []*Policy

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
	rules := e.e.GetFilteredGroupingPolicy(0, user.Encode())
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
	ps := e.e.GetPermissionsForUser(role.Encode(), domain.Encode())
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
	ps := e.e.GetFilteredPolicy(1, domain.Encode(), object.Encode())
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
	rules := e.e.GetFilteredGroupingPolicy(2, domain.Encode())
	for _, rule := range rules {
		if user, err := e.factory.User(rule[0]); err == nil {
			users = append(users, user)
		}
	}

	return users
}

func (e *enforcer) GetRolesInDomain(domain Domain) []Role {
	var roles []Role
	rules := e.e.GetFilteredGroupingPolicy(2, domain.Encode())
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
	rules := e.e.GetFilteredNamedGroupingPolicy(ObjectPType, 2, domain.Encode())
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
	ps := e.e.GetFilteredPolicy(1, domain.Encode())
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
	gp := e.e.GetFilteredGroupingPolicy(2, domain.Encode())
	var rules [][]string
	for _, rule := range gp {
		if _, err := e.factory.User(rule[0]); err == nil {
			rules = append(rules, rule)
		}
	}

	_, err := e.e.RemoveGroupingPolicies(rules)
	return err
}

func NewEnforcer(e casbin.IEnforcer, factory Factory) IEnforcer {
	return &enforcer{
		e:       e,
		factory: factory,
	}
}
