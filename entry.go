package caskin

type ObjectType string

type Action string

const (
	Read  Action = "read"
	Write Action = "write"
)

type Policy struct {
	Role   Role
	Object Object
	Domain Domain
	Action Action
}

type entry interface {
	// get id method
	GetID() uint64

	// encode and decode method
	Encode() string
	Decode(string) error

	// object method
	IsObject() bool
	GetObject() string
}

type parent interface {
	GetParentID() uint64
	SetParentID(uint64)
}

type User interface {
	entry
}

type Role interface {
	entry
	parent
}

type Object interface {
	entry
	parent
	GetObjectType() ObjectType
}

type Domain interface {
	entry
}

type EntryFactory interface {
	NewUser() User
	NewRole() Role
	NewObject() Object
	NewDomain() Domain
}
