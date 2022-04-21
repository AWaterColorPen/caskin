package caskin

type IService interface {
	IBaseService
	IDirectoryService
	IFeatureService
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

	AddUserRole(User, Domain, []*UserRolePair) error
	RemoveUserRole(User, Domain, []*UserRolePair) error
	AddRoleG(User, Domain, Role, Role) error
	RemoveRoleG(User, Domain, Role, Role) error

	GetUserByDomain(Domain) ([]User, error)
	GetDomainByUser(User) ([]Domain, error)

	GetUserRole(User, Domain) ([]*UserRolePair, error)
	GetUserRoleByUser(User, Domain, User) ([]*UserRolePair, error)
	GetUserRoleByRole(User, Domain, Role) ([]*UserRolePair, error)
	ModifyUserRolePerUser(User, Domain, User, []*UserRolePair) error
	ModifyUserRolePerRole(User, Domain, Role, []*UserRolePair) error

	GetPolicy(User, Domain) ([]*Policy, error)
	GetPolicyByRole(User, Domain, Role) ([]*Policy, error)
	ModifyPolicyPerRole(User, Domain, Role, []*Policy) error

	CreateObjectData(User, Domain, ObjectData, ObjectType) error
	RecoverObjectData(User, Domain, ObjectData) error
	DeleteObjectData(User, Domain, ObjectData) error
	UpdateObjectData(User, Domain, ObjectData, ObjectType) error
	// GetObjectData(User, Domain, ObjectData) ([]ObjectData, error) // TODO

	ObjectDataCreateCheck(User, Domain, ObjectData, ObjectType) error
	ObjectDataRecoverCheck(User, Domain, ObjectData) error
	ObjectDataDeleteCheck(User, Domain, ObjectData) error
	ObjectDataUpdateCheck(User, Domain, ObjectData, ObjectType) error
	ObjectDataModifyCheck(User, Domain, ObjectData) error
	ObjectDataGetCheck(User, Domain, ObjectData) error
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

type IDirectoryService interface {
	CreateDirectory(User, Domain, Object) error
	UpdateDirectory(User, Domain, Object) error
	DeleteDirectory(User, Domain, *DirectoryRequest) error
	GetDirectory(User, Domain, *DirectoryRequest) ([]*Directory, error)
	MoveDirectory(User, Domain, *DirectoryRequest) (*DirectoryResponse, error)
	MoveItem(User, Domain, ObjectData, *DirectoryRequest) (*DirectoryResponse, error)
	CopyItem(User, Domain, ObjectData, *DirectoryRequest) (*DirectoryResponse, error)
}

type ICurrentService interface {
	SetCurrent(User, Domain) IService

	CreateObjectWithCurrent(Object) error
	RecoverObjectWithCurrent(Object) error
	DeleteObjectWithCurrent(Object) error
	UpdateObjectWithCurrent(Object) error
	GetObjectWithCurrent(Action, ...ObjectType) ([]Object, error)

	CreateRoleWithCurrent(Role) error
	RecoverRoleWithCurrent(Role) error
	DeleteRoleWithCurrent(Role) error
	UpdateRoleWithCurrent(Role) error
	GetRole() ([]Role, error)

	AddUserRoleWithCurrent([]*UserRolePair) error
	RemoveUserRoleWithCurrent([]*UserRolePair) error
	AddRoleGWithCurrent(Role, Role) error
	RemoveRoleGWithCurrent(Role, Role) error

	GetUserByDomain(Domain) ([]User, error)
	GetDomainByUser(User) ([]Domain, error)

	GetUserRole() ([]*UserRolePair, error)
	GetUserRoleByUserWithCurrent(User) ([]*UserRolePair, error)
	GetUserRoleByRoleWithCurrent(Role) ([]*UserRolePair, error)
	ModifyUserRolePerUserWithCurrent(User, []*UserRolePair) error
	ModifyUserRolePerRoleWithCurrent(Role, []*UserRolePair) error

	GetPolicy() ([]*Policy, error)
	GetPolicyByRoleWithCurrent(Role) ([]*Policy, error)
	ModifyPolicyPerRoleWithCurrent(Role, []*Policy) error

	CreateObjectDataWithCurrent(ObjectData, ObjectType) error
	RecoverObjectDataWithCurrent(ObjectData) error
	DeleteObjectDataWithCurrent(ObjectData) error
	UpdateObjectDataWithCurrent(ObjectData, ObjectType) error
	// GetObjectDataWithCurrent(ObjectData) ([]ObjectData, error) // TODO

	CreateObjectDataCheckWithCurrent(ObjectData, ObjectType) error
	RecoverObjectDataCheckWithCurrent(ObjectData) error
	DeleteObjectDataCheckWithCurrent(ObjectData) error
	UpdateObjectDataCheckWithCurrent(ObjectData, ObjectType) error
	ModifyObjectDataCheckWithCurrent(ObjectData) error
	GetObjectDataCheckWithCurrent(ObjectData) error
}
