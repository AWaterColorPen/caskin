package caskin

const (
	// ObjectPType is the casbin named-policy type used for object hierarchy
	// edges (parent→child relationships between Objects).
	ObjectPType = "g2"

	// SuperadminRole is the casbin role name reserved for the superadmin role.
	// It is used in the dedicated superadmin domain and bypasses all normal
	// permission checks.
	SuperadminRole = "superadmin"

	// SuperadminDomain is the casbin domain name reserved for the superadmin
	// scope. Assignments in this domain confer superadmin privileges.
	SuperadminDomain = "superdomain"
)

const (
	// Read is the read permission action.
	Read Action = "read"
	// Write is the write/create/update permission action.
	Write Action = "write"
	// Manage is the administrative permission action, which typically implies
	// Read and Write as well.
	Manage Action = "manage"
)

const (
	// ObjectTypeRole is the built-in ObjectType for Role objects.
	// Roles are stored as ObjectData with this type.
	ObjectTypeRole ObjectType = "role"
)

const (
	// DirectorySearchAll traverses the entire subtree rooted at the given node.
	DirectorySearchAll DirectorySearchType = "all"
	// DirectorySearchTop returns only the direct children of the given node.
	DirectorySearchTop DirectorySearchType = "top"
)

const (
	// UpsertTypeCreate indicates the upsert operation created a new record.
	UpsertTypeCreate UpsertType = "create"
	// UpsertTypeRecover indicates the upsert operation recovered a soft-deleted record.
	UpsertTypeRecover UpsertType = "recover"
	// UpsertTypeUpdate indicates the upsert operation updated an existing record.
	UpsertTypeUpdate UpsertType = "update"
)
