package caskin

import (
	"strings"
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
}

type Domain interface {
	entry
}

func ID[E idInterface](in []E) []uint64 {
	var m []uint64
	for _, v := range in {
		m = append(m, v.GetID())
	}
	return m
}

func SB[T map[K]V, K comparable, V any](in T) []uint64 {
	return nil
}

func SB2[K comparable, V map[K]V](in V) []uint64 {
	return nil
}

func IDMap[E idInterface](in []E) map[uint64]E {
	m := map[uint64]E{}
	for _, v := range in {
		m[v.GetID()] = v
	}
	return m
}

func Tree[E TreeNodeEntry](in []E) map[uint64]uint64 {
	m := map[uint64]uint64{}
	for _, v := range in {
		if v.GetParentID() != 0 {
			m[v.GetID()] = v.GetParentID()
		}
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
