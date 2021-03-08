package caskin

type ObjectType string

type Action string

type entry interface {
	// get id method
	GetID() uint64
	// get id method
	SetID(uint64)
	// encode entry to string method
	Encode() string
	// decode string to entry method
	Decode(string) error
}

type parent interface {
	// get parent id method
	GetParentID() uint64
	// set parent id method
	SetParentID(uint64)
}

type ObjectData interface {
	// get object interface method
	GetObject() Object
	// set object
	SetObjectId(uint64)
	// set domain id method
	SetDomainID(uint64)
}

type parentEntry interface {
	entry
	parent
	ObjectData
}

type User interface {
	entry
}

type Role interface {
	parentEntry
}

type Object interface {
	parentEntry
	GetObjectType() ObjectType
}

type Domain interface {
	entry
}

type Users = []User

type Roles = []Role

type Objects = []Object

type Domains = []Domain

// EntryFactory
type EntryFactory interface {
	NewUser() User
	NewRole() Role
	NewObject() Object
	NewDomain() Domain
}

// DomainCreator create new domain's function
type DomainCreator = func(Domain) Creator

// Creator interface to create a domain
type Creator interface {
	BuildCreator() (Roles, Objects)
	SetRelation()
	GetPolicy() []*Policy
	GetRoles() Roles
	GetObjects() Objects
}

// CurrentProvider 包含发起当前请求的user及其domain
type CurrentProvider interface {
	Get() (User, Domain, error)
}

type Policy struct {
	Role   Role   `json:"role"`
	Object Object `json:"object"`
	Domain Domain `json:"domain"`
	Action Action `json:"action"`
}

type UserRolePair struct {
	User User `json:"user"`
	Role Role `json:"role"`
}

type RolesForUser struct {
	User  User  `json:"user"`
	Roles Roles `json:"roles"`
}

type UsersForRole struct {
	Role  Role  `json:"role"`
	Users Users `json:"users"`
}

type PoliciesForRole struct {
	Role     Role      `json:"role"`
	Policies []*Policy `json:"policies"`
}
