package caskin

type MetaDB interface {
	// User API
	TakeUser(User) error
	GetUserByID([]uint64) ([]User, error)

	// Role API
	TakeRole(Role) error
	GetRoleInDomain(Domain) ([]Role, error)
	GetRoleByID([]uint64) ([]Role, error)
	UpsertRole(Role) error
	DeleteRoleByID(uint64) error

	// Object API
	TakeObject(Object) error
	GetObjectInDomain(Domain, ...ObjectType) ([]Object, error)
	GetObjectByID([]uint64) ([]Object, error)
	UpsertObject(Object) error
	DeleteObjectByID(uint64) error

	// Domain API
	TakeDomain(Domain) error
	GetAllDomain() ([]Domain, error)
	DeleteDomainByID(uint64) error
}
