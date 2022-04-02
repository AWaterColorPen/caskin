package caskin

type IService interface {
	IBaseService
	ICurrentService
}

type IBaseService interface {
	SuperadminAdd(User) error
	SuperadminDelete(User) error
	SuperadminGet() ([]User, error)

	UserCreate(User) error
	UserRecover(User) error
	UserDelete(User) error
	UserUpdate(User) error

	DomainCreate(Domain) error
	DomainRecover(Domain) error
	DomainDelete(Domain) error
	DomainUpdate(Domain) error
	DomainGet() ([]Domain, error)
	DomainInitialize(Domain) error

	UserByDomainGet(Domain) ([]User, error)
	DomainByUserGet(User) ([]Domain, error)

	ObjectCreate(User, Domain, Object) error
	ObjectRecover(User, Domain, Object) error
	ObjectDelete(User, Domain, Object) error
	ObjectUpdate(User, Domain, Object) error
	ObjectGet(User, Domain, Action, ...ObjectType) ([]Object, error)

	RoleCreate(User, Domain, Role) error
	RoleRecover(User, Domain, Role) error
	RoleDelete(User, Domain, Role) error
	RoleUpdate(User, Domain, Role) error
	RoleGet(User, Domain) ([]Role, error)

	ObjectDataWriteCheck(User, Domain, ObjectData, ObjectType) error
	ObjectDataCreateCheck(User, Domain, ObjectData, ObjectType) error
	ObjectDataRecoverCheck(User, Domain, ObjectData) error
	ObjectDataDeleteCheck(User, Domain, ObjectData) error
	ObjectDataUpdateCheck(User, Domain, ObjectData, ObjectType) error
	ObjectDataModifyCheck(User, Domain, ObjectData) error
	ObjectDataGetCheck(User, Domain, ObjectData) error
}

type ICurrentService interface {
}

type IFeatureService interface {
}

type ICreatorService interface {
}
