package caskin

import (
	"encoding/json"
)

type ObjectType string

type Action string

type User interface {
	idInterface
	codeInterface
}

type Domain interface {
	idInterface
	codeInterface
}

type Role interface {
	ObjectData
	codeInterface
	// TODO no parent
	parentInterface
}

type Object interface {
	idInterface
	// nameInterface // TODO
	codeInterface
	parentInterface
	domainInterface
	GetObjectType() ObjectType
}

type ObjectData interface {
	idInterface
	domainInterface
	// GetObjectID get object
	GetObjectID() uint64
	// SetObjectID set object
	SetObjectID(uint64)
}

func ID[E idInterface](in []E) []uint64 {
	var m []uint64
	for _, v := range in {
		m = append(m, v.GetID())
	}
	return m
}

func IDMap[E idInterface](in []E) map[uint64]E {
	m := map[uint64]E{}
	for _, v := range in {
		m[v.GetID()] = v
	}
	return m
}

func Tree[E treeNode](in []E) map[uint64]uint64 {
	m := map[uint64]uint64{}
	for _, v := range in {
		if v.GetParentID() != 0 {
			m[v.GetID()] = v.GetParentID()
		}
	}
	return m
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
	b, _ := json.Marshal(s)
	return string(b)
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
