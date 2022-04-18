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

	GetUserByID([]uint64) ([]User, error)
	GetRoleInDomain(Domain) ([]Role, error)
	GetRoleByID([]uint64) ([]Role, error)
	GetObjectInDomain(Domain, ...ObjectType) ([]Object, error)
	GetObjectByID([]uint64) ([]Object, error)
	GetDomainByID([]uint64) ([]Domain, error)
	GetAllDomain() ([]Domain, error)
}
