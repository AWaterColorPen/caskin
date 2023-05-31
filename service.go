package caskin

// IService is the interface that defines all the methods for caskin service
type IService interface {
	IBaseService      // basic CRUD operations for users, domains, objects and roles
	IDirectoryService // directory-related operations for objects and object data
	IFeatureService   // feature-related operations for backends, frontends and policies
	ICurrentService   // current user-related operations
}

// IBaseService is the interface that defines the basic CRUD operations for users, domains, objects and roles
type IBaseService interface {
	// AddSuperadmin adds a superadmin user
	AddSuperadmin(User) error
	// DeleteSuperadmin deletes a superadmin user
	DeleteSuperadmin(User) error
	// GetSuperadmin gets all superadmin users
	GetSuperadmin() ([]User, error)

	// CreateUser creates a new user
	CreateUser(User) error
	// RecoverUser recovers a deleted user
	RecoverUser(User) error
	// DeleteUser deletes a user
	DeleteUser(User) error
	// UpdateUser updates a user
	UpdateUser(User) error

	// CreateDomain creates a new domain
	CreateDomain(Domain) error
	// RecoverDomain recovers a deleted domain
	RecoverDomain(Domain) error
	// DeleteDomain deletes a domain
	DeleteDomain(Domain) error
	// UpdateDomain updates a domain
	UpdateDomain(Domain) error
	// GetDomain gets all domains
	GetDomain() ([]Domain, error)
	// ResetDomain resets a domain to its initial state
	ResetDomain(Domain) error

	// CreateObject creates a new object in a domain
	CreateObject(User, Domain, Object) error
	// RecoverObject recovers a deleted object in a domain
	RecoverObject(User, Domain, Object) error
	// DeleteObject deletes an object in a domain
	DeleteObject(User, Domain, Object) error
	// UpdateObject updates an object in a domain
	UpdateObject(User, Domain, Object) error
	// GetObject gets all objects in a domain that the user can perform an action on
	GetObject(User, Domain, Action, ...ObjectType) ([]Object, error)
	// GetObjectHierarchyLevel gets the hierarchy level of an object in a domain
	GetObjectHierarchyLevel(user User, domain Domain, object Object) (int, error)

	// CreateRole creates a new role in a domain
	CreateRole(User, Domain, Role) error
	// RecoverRole recovers a deleted role in a domain
	RecoverRole(User, Domain, Role) error
	// DeleteRole deletes a role in a domain
	DeleteRole(User, Domain, Role) error
	// UpdateRole updates a role in a domain
	UpdateRole(User, Domain, Role) error
	// GetRole gets all roles in a domain
	GetRole(User, Domain) ([]Role, error)

	// AddUserRole adds user-role pairs in a domain
	AddUserRole(User, Domain, []*UserRolePair) error
	// RemoveUserRole removes user-role pairs in a domain
	RemoveUserRole(User, Domain, []*UserRolePair) error
	// AddRoleG adds a role inheritance relation in a domain
	AddRoleG(User, Domain, Role, Role) error
	// RemoveRoleG removes a role inheritance relation in a domain
	RemoveRoleG(User, Domain, Role, Role) error

	// GetUserByDomain gets all users in a domain
	GetUserByDomain(Domain) ([]User, error)
	// GetDomainByUser gets all domains that a user belongs to
	GetDomainByUser(User) ([]Domain, error)

	// GetUserRole gets all user-role pairs in a domain
	GetUserRole(User, Domain) ([]*UserRolePair, error)
	// GetUserRoleByUser gets all user-role pairs in a domain for a specific user
	GetUserRoleByUser(User, Domain, User) ([]*UserRolePair, error)
	// GetUserRoleByRole gets all user-role pairs in a domain for a specific role
	GetUserRoleByRole(User, Domain, Role) ([]*UserRolePair, error)
	// ModifyUserRolePerUser modifies the user-role pairs in a domain for a specific user
	ModifyUserRolePerUser(User, Domain, User, []*UserRolePair) error
	// ModifyUserRolePerRole modifies the user-role pairs in a domain for a specific role
	ModifyUserRolePerRole(User, Domain, Role, []*UserRolePair) error

	// GetPolicy gets all policies in a domain
	GetPolicy(User, Domain) ([]*Policy, error)
	// GetPolicyByRole gets all policies in a domain for a specific role
	GetPolicyByRole(User, Domain, Role) ([]*Policy, error)
	// ModifyPolicyPerRole modifies the policies in a domain for a specific role
	ModifyPolicyPerRole(User, Domain, Role, []*Policy) error

	// CreateObjectData creates a new object data in a domain with an object type
	CreateObjectData(User, Domain, ObjectData, ObjectType) error
	// RecoverObjectData recovers a deleted object data in a domain
	RecoverObjectData(User, Domain, ObjectData) error
	// DeleteObjectData deletes an object data in a domain
	DeleteObjectData(User, Domain, ObjectData) error
	// UpdateObjectData updates an object data in a domain with an object type
	UpdateObjectData(User, Domain, ObjectData, ObjectType) error
	// GetObjectData gets all object data in a domain
	// GetObjectData(User, Domain, ObjectData) ([]ObjectData, error) // TODO

	// CheckCreateObjectData checks if the user can create an object data in a domain with an object type
	CheckCreateObjectData(User, Domain, ObjectData, ObjectType) error
	// CheckRecoverObjectData checks if the user can recover an object data in a domain
	CheckRecoverObjectData(User, Domain, ObjectData) error
	// CheckDeleteObjectData checks if the user can delete an object data in a domain
	CheckDeleteObjectData(User, Domain, ObjectData) error
	// CheckWriteObjectData checks if the user can write an object data in a domain with an object type
	CheckWriteObjectData(User, Domain, ObjectData, ObjectType) error
	// CheckUpdateObjectData checks if the user can update an object data in a domain with an object type
	CheckUpdateObjectData(User, Domain, ObjectData, ObjectType) error
	// CheckModifyObjectData checks if the user can modify an object data in a domain
	CheckModifyObjectData(User, Domain, ObjectData) error
	// CheckGetObjectData checks if the user can get an object data in a domain
	CheckGetObjectData(User, Domain, ObjectData) error
}

// IFeatureService is the interface that defines the feature-related operations for backends, frontends and policies
type IFeatureService interface {
	// AuthBackend authenticates a user for a backend in a domain
	AuthBackend(User, Domain, *Backend) error
	// AuthFrontend authenticates a user for frontends in a domain
	AuthFrontend(User, Domain) []*Frontend
	// GetFeature gets all features in a domain
	GetFeature(User, Domain) ([]*Feature, error)
	// GetFeaturePolicy gets all feature policies in a domain
	GetFeaturePolicy(User, Domain) ([]*Policy, error)
	// GetFeaturePolicyByRole gets all feature policies in a domain for a specific role
	GetFeaturePolicyByRole(User, Domain, Role) ([]*Policy, error)
	// ModifyFeaturePolicyPerRole modifies the feature policies in a domain for a specific role
	ModifyFeaturePolicyPerRole(User, Domain, Role, []*Policy) error
	// ResetFeature resets the features in a domain to their initial state
	ResetFeature(Domain) error
}

// IDirectoryService is the interface that defines the directory-related operations for objects and object data
type IDirectoryService interface {
	// CreateDirectory creates a new directory for an object in a domain
	CreateDirectory(User, Domain, Object) error
	// UpdateDirectory updates an existing directory for an object in a domain
	UpdateDirectory(User, Domain, Object) error
	// DeleteDirectory deletes a directory and its subdirectories in a domain based on a request
	DeleteDirectory(User, Domain, *DirectoryRequest) error
	// GetDirectory gets all directories and their subdirectories in a domain based on a request
	GetDirectory(User, Domain, *DirectoryRequest) ([]*Directory, error)
	// MoveDirectory moves a directory and its subdirectories to another directory in a domain based on a request and returns the updated directory structure
	MoveDirectory(User, Domain, *DirectoryRequest) (*DirectoryResponse, error)
	// MoveItem moves an object data to another directory in a domain based on a request and returns the updated directory structure
	MoveItem(User, Domain, ObjectData, *DirectoryRequest) (*DirectoryResponse, error)
	// CopyItem copies an object data to another directory in a domain based on a request and returns the updated directory structure
	CopyItem(User, Domain, ObjectData, *DirectoryRequest) (*DirectoryResponse, error)
}

// ICurrentService is the interface that defines the current user-related operations
type ICurrentService interface {
	// SetCurrent sets the current user and domain for the service and returns a new service instance
	SetCurrent(User, Domain) IService

	// CreateObjectWithCurrent creates a new object in the current domain
	// CreateObjectWithCurrent(Object) error
	// RecoverObjectWithCurrent recovers a deleted object in the current domain
	// RecoverObjectWithCurrent(Object) error
	// DeleteObjectWithCurrent deletes an object in the current domain
	// DeleteObjectWithCurrent(Object) error
	// UpdateObjectWithCurrent updates an object in the current domain
	// UpdateObjectWithCurrent(Object) error
	// GetObjectWithCurrent gets all objects in the current domain that the current user can perform an action on
	// GetObjectWithCurrent(Action, ...ObjectType) ([]Object, error)
	//
	// CreateRoleWithCurrent creates a new role in the current domain
	// CreateRoleWithCurrent(Role) error
	// RecoverRoleWithCurrent recovers a deleted role in the current domain
	// RecoverRoleWithCurrent(Role) error
	// DeleteRoleWithCurrent deletes a role in the current domain
	// DeleteRoleWithCurrent(Role) error
	// UpdateRoleWithCurrent updates a role in the current domain
	// UpdateRoleWithCurrent(Role) error
	// GetRoleWithCurrent gets all roles in the current domain
	// GetRoleWithCurrent() ([]Role, error)
	//
	// AddUserRoleWithCurrent adds user-role pairs in the current domain
	// AddUserRoleWithCurrent([]*UserRolePair) error
	// RemoveUserRoleWithCurrent removes user-role pairs in the current domain
	// RemoveUserRoleWithCurrent([]*UserRolePair) error
	// AddRoleGWithCurrent adds a role inheritance relation in the current domain
	// AddRoleGWithCurrent(Role, Role) error
	// RemoveRoleGWithCurrent removes a role inheritance relation in the current domain
	// RemoveRoleGWithCurrent(Role, Role) error
	//
	// GetUserRoleWithCurrent gets all user-role pairs in the current domain
	// GetUserRoleWithCurrent() ([]*UserRolePair, error)
	// GetUserRoleByUserWithCurrent gets all user-role pairs in the current domain for a specific user
	// GetUserRoleByUserWithCurrent(User) ([]*UserRolePair, error)
	// GetUserRoleByRoleWithCurrent gets all user-role pairs in the current domain for a specific role
	// GetUserRoleByRoleWithCurrent(Role) ([]*UserRolePair, error)
	// ModifyUserRolePerUserWithCurrent modifies the user-role pairs in the current domain for a specific user
	// ModifyUserRolePerUserWithCurrent(User, []*UserRolePair) error
	// ModifyUserRolePerRoleWithCurrent modifies the user-role pairs in the current domain for a specific role
	// ModifyUserRolePerRoleWithCurrent(Role, []*UserRolePair) error
	//
	// GetPolicyWithCurrent gets all policies in the current domain
	// GetPolicyWithCurrent() ([]*Policy, error)
	// GetPolicyByRoleWithCurrent gets all policies in the current domain for a specific role
	// GetPolicyByRoleWithCurrent(Role) ([]*Policy, error)
	// ModifyPolicyPerRoleWithCurrent modifies the policies in the current domain for a specific role
	// ModifyPolicyPerRoleWithCurrent(Role, []*Policy) error

	// CreateObjectDataWithCurrent creates a new object data in the current domain with an object type
	CreateObjectDataWithCurrent(ObjectData, ObjectType) error
	// RecoverObjectDataWithCurrent recovers a deleted object data in the current domain
	RecoverObjectDataWithCurrent(ObjectData) error
	// DeleteObjectDataWithCurrent deletes an object data in the current domain
	DeleteObjectDataWithCurrent(ObjectData) error
	// UpdateObjectDataWithCurrent updates an object data in the current domain with an object type
	UpdateObjectDataWithCurrent(ObjectData, ObjectType) error
	// GetObjectDataWithCurrent gets all object data in the current domain
	// GetObjectDataWithCurrent(ObjectData) ([]ObjectData, error) // TODO

	// CheckCreateObjectDataWithCurrent checks if the current user can create an object data in the current domain with an object type
	CheckCreateObjectDataWithCurrent(ObjectData, ObjectType) error
	// CheckRecoverObjectDataWithCurrent checks if the current user can recover an object data in the current domain
	CheckRecoverObjectDataWithCurrent(ObjectData) error
	// CheckDeleteObjectDataWithCurrent checks if the current user can delete an object data in the current domain
	CheckDeleteObjectDataWithCurrent(ObjectData) error
	// CheckWriteObjectDataWithCurrent checks if the current user can write an object data in the current domain with an object type
	CheckWriteObjectDataWithCurrent(ObjectData, ObjectType) error
	// CheckUpdateObjectDataWithCurrent checks if the current user can update an object data in the current domain with an object type
	CheckUpdateObjectDataWithCurrent(ObjectData, ObjectType) error
	// CheckModifyObjectDataWithCurrent checks if the current user can modify an object data in the current domain
	CheckModifyObjectDataWithCurrent(ObjectData) error
	// CheckGetObjectDataWithCurrent checks if the current user can get an object data in the current domain
	CheckGetObjectDataWithCurrent(ObjectData) error
}
