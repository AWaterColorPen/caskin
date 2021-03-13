package caskin

type MetaDB interface {
	Create(interface{}) error
	Recover(interface{}) error
	Update(interface{}) error
	Upsert(interface{}) error
	Take(interface{}) error
	TakeUnscoped(interface{}) error
	DeleteByID(interface{}, uint64) error

	// User API
	GetUserByID([]uint64) ([]User, error)

	// Role API
	GetRoleInDomain(Domain) ([]Role, error)
	GetRoleByID([]uint64) ([]Role, error)

	// Object API
	GetObjectInDomain(Domain, ...ObjectType) ([]Object, error)
	GetObjectByID([]uint64) ([]Object, error)

	// Domain API
	GetAllDomain() ([]Domain, error)
}

type MetaDBBindObjectAPI interface {
	Create(ObjectData, Object) error
	Recover(ObjectData, Object) error
	Update(ObjectData, Object) error
	Delete(ObjectData, Object) error
}
