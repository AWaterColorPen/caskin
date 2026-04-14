package caskin

// UpsertType describes the outcome of a database upsert operation performed
// by [MetaDB]. See [UpsertTypeCreate], [UpsertTypeRecover], and
// [UpsertTypeUpdate] for the possible values.
type UpsertType string

// MetaDB is the storage abstraction used by caskin to persist Users, Roles,
// Objects, Domains, and ObjectData. Provide a concrete implementation (or use
// the built-in GORM-backed one via [Register]) to connect caskin to your
// database.
//
// The interface deliberately remains generic (operating on any) for Create,
// Recover, Update, and similar operations so that it can handle all entity
// types without requiring separate methods per type.
type MetaDB interface {
	// Create inserts a new record. The entity must not already exist (not
	// soft-deleted). Returns [ErrAlreadyExists] if a live record is found.
	Create(any) error
	// Recover undeletes a soft-deleted record. Returns [ErrNotExists] if no
	// deleted record is found.
	Recover(any) error
	// Update saves changes to an existing record.
	Update(any) error
	// UpsertType determines whether the next Upsert will create, recover, or
	// update; it does not perform any write itself.
	UpsertType(any) UpsertType
	// Take loads a single live (not soft-deleted) record by primary key.
	Take(any) error
	// TakeUnscoped loads a single record regardless of soft-delete status.
	TakeUnscoped(any) error
	// Find loads all matching records into the destination slice, with optional
	// filter conditions forwarded to the underlying ORM.
	Find(any, ...any) error
	// DeleteByID soft-deletes the record of the given type with the specified ID.
	DeleteByID(any, uint64) error

	// GetUserByID fetches [User] records for the given IDs.
	GetUserByID([]uint64) ([]User, error)
	// GetRoleInDomain fetches all [Role] records that belong to the given domain.
	GetRoleInDomain(Domain) ([]Role, error)
	// GetRoleByID fetches [Role] records for the given IDs.
	GetRoleByID([]uint64) ([]Role, error)
	// GetObjectInDomain fetches all [Object] records in the given domain,
	// optionally filtered by one or more [ObjectType] values.
	GetObjectInDomain(Domain, ...ObjectType) ([]Object, error)
	// GetObjectByID fetches [Object] records for the given IDs.
	GetObjectByID([]uint64) ([]Object, error)
	// GetDomainByID fetches [Domain] records for the given IDs.
	GetDomainByID([]uint64) ([]Domain, error)
	// GetAllDomain fetches every [Domain] record in the database.
	GetAllDomain() ([]Domain, error)
}
