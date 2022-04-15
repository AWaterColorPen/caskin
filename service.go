package caskin

type IService interface {
	IBaseService
	ICurrentService
}

type IBaseService interface {
	AddSuperadmin(User) error
	DeleteSuperadmin(User) error
	GetSuperadmin() ([]User, error)

	CreateUser(User) error
	RecoverUser(User) error
	DeleteUser(User) error
	UpdateUser(User) error

	CreateDomain(Domain) error
	RecoverDomain(Domain) error
	DeleteDomain(Domain) error
	UpdateDomain(Domain) error
	GetDomain() ([]Domain, error)
	ResetDomain(Domain) error

	CreateObject(User, Domain, Object) error
	RecoverObject(User, Domain, Object) error
	DeleteObject(User, Domain, Object) error
	UpdateObject(User, Domain, Object) error
	GetObject(User, Domain, Action, ...ObjectType) ([]Object, error)

	CreateRole(User, Domain, Role) error
	RecoverRole(User, Domain, Role) error
	DeleteRole(User, Domain, Role) error
	UpdateRole(User, Domain, Role) error
	GetRole(User, Domain) ([]Role, error)

	GetUserByDomain(Domain) ([]User, error)
	GetDomainByUser(User) ([]Domain, error)

	GetUserRole(User, Domain) ([]*UserRolePair, error)
	GetUserRoleByUser(User, Domain, User) ([]*UserRolePair, error)
	GetUserRoleByRole(User, Domain, Role) ([]*UserRolePair, error)
	ModifyUserRolePerUser(User, Domain, User, []*UserRolePair) error
	ModifyUserRolePerRole(User, Domain, Role, []*UserRolePair) error

	GetPolicy(User, Domain) ([]*Policy, error)
	GetPolicyByRole(User, Domain, Role) ([]*Policy, error)
	GetPolicyByObject(User, Domain, Object) ([]*Policy, error)
	ModifyPolicyPerRole(User, Domain, Role, []*Policy) error
	ModifyPolicyPerObject(User, Domain, Object, []*Policy) error

	CreateObjectData(User, Domain, ObjectData, ObjectType) error
	RecoverObjectData(User, Domain, ObjectData) error
	DeleteObjectData(User, Domain, ObjectData) error
	UpdateObjectData(User, Domain, ObjectData, ObjectType) error
	// GetObjectData(User, Domain, ObjectData) error // TODO

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
	GetBackend() ([]*Backend, error)
	GetFrontend() ([]*Frontend, error)
	GetFeature() ([]*Feature, error)
	AuthBackend(User, Domain, *Backend) error
	AuthFrontend(User, Domain) []*Frontend
	GetFeatureObject(User, Domain) ([]Object, error)
	GetFeaturePolicy(User, Domain) ([]*Policy, error)
	GetFeaturePolicyByRole(User, Domain, Role) ([]*Policy, error)
	ModifyFeaturePolicyPerRole(User, Domain, Role, []*Policy) error
	ResetFeature(Domain) error
}

type ICreatorService interface {
}
