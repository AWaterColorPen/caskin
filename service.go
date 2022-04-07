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
	DomainReset(Domain) error

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

	UserByDomainGet(Domain) ([]User, error)
	DomainByUserGet(User) ([]Domain, error)

	UserRoleGet(User, Domain) ([]*UserRolePair, error)
	UserRoleByUserGet(User, Domain, User) ([]*UserRolePair, error)
	UserRoleByRoleGet(User, Domain, Role) ([]*UserRolePair, error)
	UserRolePerUserModify(User, Domain, User, []*UserRolePair) error
	UserRolePerRoleModify(User, Domain, Role, []*UserRolePair) error

	PolicyGet(User, Domain) ([]*Policy, error)
	PolicyByRoleGet(User, Domain, Role) ([]*Policy, error)
	PolicyByObjectGet(User, Domain, Object) ([]*Policy, error)
	PolicyPerRoleModify(User, Domain, Role, []*Policy) error
	PolicyPerObjectModify(User, Domain, Object, []*Policy) error
}

type ICurrentService interface {
}

type IFeatureService interface {
	GetBackend() ([]*Backend, error)
	GetFrontend() ([]*Frontend, error)
	GetFeature() ([]*Feature, error)
	BackendAuth(User, Domain, *Backend) error
	FrontendAuth(User, Domain) []*Frontend
	FeatureObjectGet(User, Domain) ([]Object, error)
	FeaturePolicyGet(User, Domain) ([]*Policy, error)
	FeaturePolicyByRoleGet(User, Domain, Role) ([]*Policy, error)
	FeaturePolicyPerRoleModify(User, Domain, Role, []*Policy) error
	FeatureReset(Domain) error
}

type ICreatorService interface {
}

type server struct {
	Enforcer   IEnforcer
	DB         MetaDB
	Dictionary IDictionary
}

type currentServer struct {
	server
	CurrentUser   User
	CurrentDomain Domain
}
