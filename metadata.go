package caskin

type MetaDB interface {
	Create(interface{}) error
	Recover(interface{}) error
	Update(interface{}) error
	Upsert(interface{}) error
	Take(interface{}) error
	TakeUnscoped(interface{}) error

	// User API
	GetUserByID([]uint64) ([]User, error)
	DeleteUserByID(uint64) error

	// Role API
	GetRoleInDomain(Domain) ([]Role, error)
	GetRoleByID([]uint64) ([]Role, error)
	DeleteRoleByID(uint64) error

	// Object API
	GetObjectInDomain(Domain, ...ObjectType) ([]Object, error)
	GetObjectByID([]uint64) ([]Object, error)
	DeleteObjectByID(uint64) error

	// Domain API
	GetAllDomain() ([]Domain, error)
	DeleteDomainByID(uint64) error
}

type MetaDBBindObjectAPI interface {
	Create(ObjectData, Object) error
	Recover(ObjectData, Object) error
	Update(ObjectData, Object) error
	DeleteByID(objectDataID, bindObjectID uint64) error
}
