package caskin

import (
	"encoding/json"
)

type DirectorySearchType = string

type ObjectType = string

type Action = string

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
}

type Object interface {
	idInterface
	codeInterface
	parentInterface
	domainInterface
	GetObjectType() string
}

type ObjectData interface {
	idInterface
	domainInterface
	// GetObjectID get object
	GetObjectID() uint64
	// SetObjectID set object
	SetObjectID(uint64)
}

type IDirectory interface {
	Search(uint64, DirectorySearchType) []*Directory
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

// Policy tuple of role-object-domain-action
type Policy struct {
	Role   Role   `json:"role"`
	Object Object `json:"object"`
	Domain Domain `json:"domain"`
	Action Action `json:"action"`
}

// Key get the unique identify of the policy
func (p *Policy) Key() string {
	s := []string{p.Role.Encode(), p.Object.Encode(), p.Domain.Encode(), p.Action}
	b, _ := json.Marshal(s)
	return string(b)
}

// UserRolePair pair of user and role
type UserRolePair struct {
	User User `json:"user"`
	Role Role `json:"role"`
}

type Directory struct {
	Object
	AllDirectoryCount uint64 `json:"all_directory_count"`
	AllItemCount      uint64 `json:"all_item_count"`
	TopDirectoryCount uint64 `json:"top_directory_count"`
	TopItemCount      uint64 `json:"top_item_count"`
}

type DirectoryRequest struct {
	To              uint64   `json:"to,omitempty"`
	ID              []uint64 `json:"id,omitempty"`
	Type            string   `json:"type,omitempty"`
	Policy          string   `json:"policy,omitempty"`
	SearchType      string   `json:"search_type,omitempty"`
	CountDirectory  func([]uint64) (map[uint64]uint64, error)
	ActionDirectory func([]uint64) error
}

type CountDirectoryItem = func([]uint64) (map[uint64]uint64, error)

type DirectoryResponse struct {
	DoneDirectoryCount uint64 `json:"done_directory_count,omitempty"`
	DoneItemCount      uint64 `json:"done_item_count,omitempty"`
	ToDoDirectoryCount uint64 `json:"to_do_directory_count,omitempty"`
	ToDoItemCount      uint64 `json:"to_do_item_count,omitempty"`
}
