package caskin

type MetaDB interface {
	// User API
	TakeUser(User) error
	GetUserByID([]uint64) ([]User, error)

	// Role API
	CreateRole(Role) error
	RecoverRole(Role) error
	UpdateRole(Role) error
	TakeRole(Role) error
	GetRoleInDomain(Domain) ([]Role, error)
	GetRoleByID([]uint64) ([]Role, error)
	UpsertRole(Role) error
	DeleteRoleByID(uint64) error

	// Object API
	CreateObject(Object) error
	RecoverObject(Object) error
	UpdateObject(Object) error
	TakeObject(Object) error
	GetObjectInDomain(Domain, ...ObjectType) ([]Object, error)
	GetObjectByID([]uint64) ([]Object, error)
	UpsertObject(Object) error
	DeleteObjectByID(uint64) error

	// Domain API
	CreateDomain(Domain) error
	RecoverDomain(Domain) error
	UpdateDomain(Domain) error
	TakeDomain(Domain) error
	GetAllDomain() ([]Domain, error)
	DeleteDomainByID(uint64) error
}
