package example

import "github.com/awatercolorpen/caskin"

// EntryFactory sample for caskin.EntryFactory interface
type EntryFactory struct {
}

func (e *EntryFactory) NewUser() caskin.User {
	return &User{}
}

func (e *EntryFactory) NewRole() caskin.Role {
	return &Role{}
}

func (e *EntryFactory) NewObject() caskin.Object {
	return &Object{}
}

func (e *EntryFactory) NewDomain() caskin.Domain {
	return &Domain{}
}
