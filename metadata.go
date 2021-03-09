package caskin

type MetaDB interface {
	Create(interface{}) error
	Recover(interface{}) error
	Update(interface{}) error
	Take(interface{}) error

	// User API
	GetUserInDomain(Domain) ([]User, error)
	GetUserByID([]uint64) ([]User, error)
	UpsertUser(User) error
	DeleteUserByID(uint64) error

	// Role API
	GetRoleInDomain(Domain) ([]Role, error)
	GetRoleByID([]uint64) ([]Role, error)
	UpsertRole(Role) error
	DeleteRoleByID(uint64) error

	// Object API
	GetObjectInDomain(Domain, ...ObjectType) ([]Object, error)
	GetObjectByID([]uint64) ([]Object, error)
	UpsertObject(Object) error
	DeleteObjectByID(uint64) error

	// Domain API
	GetAllDomain() ([]Domain, error)
	DeleteDomainByID(uint64) error
}
