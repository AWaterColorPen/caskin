package caskin

// EntryFactory
type EntryFactory interface {
	NewUser() User
	NewRole() Role
	NewObject() Object
	NewDomain() Domain
}

type UserFactory = func() User
type RoleFactory = func() Role
type ObjectFactory = func() Object
type DomainFactory = func() Domain
