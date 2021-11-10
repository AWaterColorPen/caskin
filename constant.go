package caskin

const (
	ObjectPType = "g2"

	SuperadminRole   = "superadmin"
	SuperadminDomain = "superdomain"
)

const (
	Read  Action = "read"
	Write Action = "write"
)

const (
	ObjectTypeDefault ObjectType = "default"
	ObjectTypeObject  ObjectType = "object"
	ObjectTypeRole    ObjectType = "role"
)

const (
	UpsertTypeCreate  UpsertType = "create"
	UpsertTypeRecover UpsertType = "recover"
	UpsertTypeUpdate  UpsertType = "update"
)
