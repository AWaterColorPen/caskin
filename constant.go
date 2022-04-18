package caskin

const (
	ObjectPType = "g2"

	SuperadminRole   = "superadmin"
	SuperadminDomain = "superdomain"
)

const (
	Read   Action = "read"
	Write  Action = "write"
	Manage Action = "manage"
)

const (
	ObjectTypeRole ObjectType = "role"
)

const (
	UpsertTypeCreate  UpsertType = "create"
	UpsertTypeRecover UpsertType = "recover"
	UpsertTypeUpdate  UpsertType = "update"
)
