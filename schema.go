package caskin

import (
	"strings"

	"gorm.io/datatypes"
)

type ObjectType string

type Action string

type ObjectData interface {
	idInterface
	// GetObject get object interface method
	GetObject() Object
	// SetObjectID set object
	SetObjectID(uint64)
	// GetDomainID get domain id method
	GetDomainID() uint64
	// SetDomainID set domain id method
	SetDomainID(uint64)
}

type TreeNodeEntry interface {
	entry
	treeNode
	ObjectData
}

type User interface {
	entry
}

type Role interface {
	TreeNodeEntry
	nameInterface
}

type Object interface {
	TreeNodeEntry
	nameInterface
	GetObjectType() ObjectType
	SetObjectType(ObjectType)
	GetCustomizedData() datatypes.JSON
	SetCustomizedData(datatypes.JSON)
}

type Domain interface {
	entry
}

type Users []User

func (u Users) ID() []uint64 {
	var m []uint64
	for _, v := range u {
		m = append(m, v.GetID())
	}
	return m
}

func (u Users) IDMap() map[uint64]User {
	m := map[uint64]User{}
	for _, v := range u {
		m[v.GetID()] = v
	}
	return m
}

type Roles []Role

func (r Roles) ID() []uint64 {
	var m []uint64
	for _, v := range r {
		m = append(m, v.GetID())
	}
	return m
}

func (r Roles) IDMap() map[uint64]Role {
	m := map[uint64]Role{}
	for _, v := range r {
		m[v.GetID()] = v
	}
	return m
}

func (r Roles) Tree() map[uint64]uint64 {
	m := map[uint64]uint64{}
	for _, v := range r {
		if v.GetParentID() != 0 {
			m[v.GetID()] = v.GetParentID()
		}
	}
	return m
}

type Objects []Object

func (o Objects) ID() []uint64 {
	var m []uint64
	for _, v := range o {
		m = append(m, v.GetID())
	}
	return m
}

func (o Objects) IDMap() map[uint64]Object {
	m := map[uint64]Object{}
	for _, v := range o {
		m[v.GetID()] = v
	}
	return m
}

func (o Objects) Tree() map[uint64]uint64 {
	m := map[uint64]uint64{}
	for _, v := range o {
		if v.GetParentID() != 0 {
			m[v.GetID()] = v.GetParentID()
		}
	}
	return m
}

type Domains []Domain

func (d Domains) ID() []uint64 {
	var m []uint64
	for _, v := range d {
		m = append(m, v.GetID())
	}
	return m
}

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
