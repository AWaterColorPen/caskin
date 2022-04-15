package caskin

type UpsertType string

type MetaDB interface {
	Create(any) error
	Recover(any) error
	Update(any) error
	UpsertType(any) UpsertType
	Take(any) error
	TakeUnscoped(any) error
	Find(any, ...any) error
	DeleteByID(any, uint64) error

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
