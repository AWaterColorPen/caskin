package caskin

type IService interface {
}

type ICurrentService interface {
}

type IBaseService interface {
	SuperadminUserAdd(User) error
	SuperadminUserDelete(User) error
	SuperadminUserGet() ([]User, error)
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
	UserInDomainGet(Domain) ([]User, error)
	DomainByUserGet(User) ([]Domain, error)
	ObjectCreate(User, Domain, Object) error
	ObjectRecover(User, Domain, Object) error
	ObjectDelete(User, Domain, Object) error
	ObjectUpdate(User, Domain, Object) error
	ObjectGet(User, Domain, ...ObjectType) ([]Object, error)
}

type IFeatureService interface {
}

type ICreatorService interface {
}
