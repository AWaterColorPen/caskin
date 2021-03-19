package caskin

import (
	"strings"

	"gorm.io/datatypes"
)

type ObjectType string

type Action string

type ObjectData interface {
	idInterface
	// get object interface method
	GetObject() Object
	// set object
	SetObjectID(uint64)
	// set domain id method
	SetDomainID(uint64)
}

type User interface {
	entry
}

type Role interface {
	treeNodeEntry
}

type Object interface {
	treeNodeEntry
	GetName() string
	SetName(string)
	GetObjectType() ObjectType
	SetObjectType(ObjectType)
	GetCustomizedData() datatypes.JSON
	SetCustomizedData(datatypes.JSON)
}

type Domain interface {
	entry
}

type Users = []User

type Roles []Role

type Objects []Object

func (o Objects) IDMap() map[uint64]Object {
	m := map[uint64]Object{}
	for _, v := range o {
		m[v.GetID()] = v
	}
	return m
}

type Domains = []Domain

// DomainCreator create new domain's function
type DomainCreator = func(Domain) Creator

// Creator interface to create a domain
type Creator interface {
	BuildCreator() ([]Role, []Object)
	SetRelation()
	GetPolicy() []*Policy
	GetRoles() []Role
	GetObjects() []Object
}

// Policy tuple of role-object-domain-action
type Policy struct {
	Role   Role   `json:"role"`
	Object Object `json:"object"`
	Domain Domain `json:"domain"`
	Action Action `json:"action"`
}

// Key get the unique identify of the policy
func (p *Policy) Key() string {
	s := []string{p.Role.Encode(), p.Object.Encode(), p.Domain.Encode(), string(p.Action)}
	return strings.Join(s, DefaultSeparator)
}

// PolicyList list of policy
type PolicyList []*Policy

func (p PolicyList) IsValidWithObject(object Object) error {
	encode := object.Encode()
	for _, v := range p {
		if v.Object.Encode() != encode {
			return ErrInputPolicyListNotBelongSameObject
		}
	}
	return nil
}

func (p PolicyList) IsValidWithRole(role Role) error {
	encode := role.Encode()
	for _, v := range p {
		if v.Role.Encode() != encode {
			return ErrInputPolicyListNotBelongSameRole
		}
	}
	return nil
}

// UserRolePair pair of user and role
type UserRolePair struct {
	User User `json:"user"`
	Role Role `json:"role"`
}

// UserRolePairs list of user and role's pair
type UserRolePairs []*UserRolePair

func (u UserRolePairs) IsValidWithRole(role Role) error {
	encode := role.Encode()
	for _, v := range u {
		if v.Role.Encode() != encode {
			return ErrInputPairArrayNotBelongSameRole
		}
	}
	return nil
}

func (u UserRolePairs) IsValidWithUser(user User) error {
	encode := user.Encode()
	for _, v := range u {
		if v.User.Encode() != encode {
			return ErrInputPairArrayNotBelongSameUser
		}
	}
	return nil
}

func (u UserRolePairs) RoleID() []uint64 {
	var id []uint64
	for _, v := range u {
		id = append(id, v.Role.GetID())
	}
	return id
}

func (u UserRolePairs) UserID() []uint64 {
	var id []uint64
	for _, v := range u {
		id = append(id, v.User.GetID())
	}
	return id
}

type CustomizedDataPair struct {
	Object               Object         `json:"object"`
	ObjectCustomizedData CustomizedData `json:"customized_data"`
}

// InheritanceRelation value is sons' id list
type InheritanceRelation = []uint64

// InheritanceRelations key is parent id, value is sons' id list
type InheritanceRelations = map[uint64]InheritanceRelation
