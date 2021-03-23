package caskin

const (
	ObjectPType = "g2"

	SuperadminRole   = "superadmin"
	SuperadminDomain = "superdomain"

	ModelConfPathSuperadmin   = "http://raw.githubusercontent.com/awatercolorpen/caskin/main/configs/casbin_model.conf"
	ModelConfPathNoSuperadmin = "http://raw.githubusercontent.com/awatercolorpen/caskin/main/configs/casbin_model.no_superadmin.conf"
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
