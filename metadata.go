package caskin

type UpsertType string

type MetaDB interface {
	AutoMigrate(...any) error

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

func GetByID[T any](db MetaDB, id []uint64) ([]T, error) {
	var out []T
	if err := db.Find(&out, "id IN ?", id); err != nil {
		return nil, err
	}
	return out, nil
}

func GetInDomain[T any](db MetaDB, domain Domain) ([]T, error) {
	var out []T
	if err := db.Find(&out, "domain_id = ?", domain.GetID()); err != nil {
		return nil, err
	}
	return out, nil
}
