package caskin

type UpsertType string

type MetaDB interface {
	AutoMigrate(...interface{}) error

	Create(interface{}) error
	Recover(interface{}) error
	Update(interface{}) error
	UpsertType(interface{}) UpsertType
	Take(interface{}) error
	TakeUnscoped(interface{}) error
	First(interface{}, ...interface{}) error
	Find(interface{}, ...interface{}) error
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
	GetDomainByID([]uint64) ([]Domain, error)
	GetAllDomain() ([]Domain, error)
}
