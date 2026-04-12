package caskin

import (
	"encoding/json"
)

// DirectorySearchType specifies the scope of a directory search.
// Use [DirectorySearchAll] to include all descendants, or [DirectorySearchTop]
// to include only direct children.
type DirectorySearchType = string

// ObjectType is a string tag that categorises objects (e.g. "role", "menu",
// "api"). It is used to look up the correct GORM model when performing
// object-data operations.
type ObjectType = string

// Action is a permission verb. caskin defines three built-in actions:
// [Read], [Write], and [Manage]. Custom actions are also supported.
type Action = string

// User represents an actor in the permission system. Implementations must
// provide an integer ID ([idInterface]) and a string encoding
// ([codeInterface]) used as the casbin subject.
type User interface {
	idInterface
	codeInterface
}

// Domain represents an isolated permission scope such as a tenant or
// organization. Implementations must provide an integer ID and a string
// encoding used as the casbin domain token.
type Domain interface {
	idInterface
	codeInterface
}

// Role represents a named permission bundle within a domain. Roles are
// also [ObjectData], meaning they belong to an object type and can be
// organised in a hierarchy via [IBaseService.AddRoleG].
type Role interface {
	ObjectData
	codeInterface
}

// Object represents a resource in the permission system. Objects are
// arranged in a tree: each object may have a parent ([parentInterface]).
// A user's read/write/manage access to a child object is implicitly
// determined by the policy set on any ancestor.
type Object interface {
	idInterface
	codeInterface
	parentInterface
	domainInterface
	// GetObjectType returns the type tag for this object (e.g. "menu").
	GetObjectType() string
}

// ObjectData represents a domain-specific data record that is protected by
// an [Object]. For example, a "document" entity might be protected by a
// "documents" object. Permission checks are done via the associated object.
type ObjectData interface {
	idInterface
	domainInterface
	// GetObjectID returns the ID of the protecting Object.
	GetObjectID() uint64
	// SetObjectID sets the ID of the protecting Object.
	SetObjectID(uint64)
}

// IDirectory is the interface for directory search backends. It returns the
// list of [Directory] nodes reachable from the given root ID according to
// the specified [DirectorySearchType].
type IDirectory interface {
	Search(uint64, DirectorySearchType) []*Directory
}

// ID extracts the integer IDs from a slice of entities.
func ID[E idInterface](in []E) []uint64 {
	var m []uint64
	for _, v := range in {
		m = append(m, v.GetID())
	}
	return m
}

// IDMap converts a slice of entities into a map keyed by their integer ID.
func IDMap[E idInterface](in []E) map[uint64]E {
	m := map[uint64]E{}
	for _, v := range in {
		m[v.GetID()] = v
	}
	return m
}

// Policy is a 4-tuple of (Role, Object, Domain, Action) that grants the
// role permission to perform the action on the object within the domain.
type Policy struct {
	Role   Role   `json:"role"`
	Object Object `json:"object"`
	Domain Domain `json:"domain"`
	Action Action `json:"action"`
}

// Key returns a stable string that uniquely identifies this policy tuple.
// It is used internally for set-difference operations (e.g. [DiffPolicy]).
func (p *Policy) Key() string {
	s := []string{p.Role.Encode(), p.Object.Encode(), p.Domain.Encode(), p.Action}
	b, _ := json.Marshal(s)
	return string(b)
}

// UserRolePair binds a user to a role. It is the unit used in
// [IBaseService.AddUserRole] and related methods.
type UserRolePair struct {
	User User `json:"user"`
	Role Role `json:"role"`
}

// Directory decorates an [Object] with aggregate counts that describe
// the subtree rooted at that object.
type Directory struct {
	Object
	// AllDirectoryCount is the total number of directory nodes in the subtree.
	AllDirectoryCount uint64 `json:"all_directory_count"`
	// AllItemCount is the total number of leaf items in the subtree.
	AllItemCount uint64 `json:"all_item_count"`
	// TopDirectoryCount is the number of direct child directories.
	TopDirectoryCount uint64 `json:"top_directory_count"`
	// TopItemCount is the number of direct child items.
	TopItemCount uint64 `json:"top_item_count"`
}

// DirectoryRequest is the parameter bag for directory operations such as
// [IDirectoryService.GetDirectory], [IDirectoryService.MoveDirectory], and
// [IDirectoryService.DeleteDirectory].
type DirectoryRequest struct {
	// To is the target parent object ID for move operations.
	To uint64 `json:"to,omitempty"`
	// ID is the list of object IDs to operate on.
	ID []uint64 `json:"id,omitempty"`
	// Type is the ObjectType filter.
	Type string `json:"type,omitempty"`
	// Policy is an optional policy filter string.
	Policy string `json:"policy,omitempty"`
	// SearchType controls the scope of the directory traversal;
	// see [DirectorySearchAll] and [DirectorySearchTop].
	SearchType string `json:"search_type,omitempty"`
	// CountDirectory is an optional callback that returns per-directory item counts.
	CountDirectory func([]uint64) (map[uint64]uint64, error)
	// ActionDirectory is an optional callback that performs a side effect on a set of directories.
	ActionDirectory func([]uint64) error
}

// CountDirectoryItem is the function signature for per-directory item count
// callbacks used in [DirectoryRequest].
type CountDirectoryItem = func([]uint64) (map[uint64]uint64, error)

// DirectoryResponse summarises the outcome of a directory move or copy
// operation, reporting how many directories and items were already at the
// destination ("done") versus still pending ("to-do").
type DirectoryResponse struct {
	// DoneDirectoryCount is the number of directories already at the destination.
	DoneDirectoryCount uint64 `json:"done_directory_count,omitempty"`
	// DoneItemCount is the number of items already at the destination.
	DoneItemCount uint64 `json:"done_item_count,omitempty"`
	// ToDoDirectoryCount is the number of directories moved/copied.
	ToDoDirectoryCount uint64 `json:"to_do_directory_count,omitempty"`
	// ToDoItemCount is the number of items moved/copied.
	ToDoItemCount uint64 `json:"to_do_item_count,omitempty"`
}
