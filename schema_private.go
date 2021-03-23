package caskin

// idInterface it is only for get/set id method
type idInterface interface {
	// get id method
	GetID() uint64
	// set id method
	SetID(uint64)
}

// nameInterface it is only for get/set name method
type nameInterface interface {
	// get name method
	GetName() string
	// set name method
	SetName(string)
}

// entry it is casbin entry of User Role Object Domain
type entry interface {
	idInterface
	// encode entry to string method
	Encode() string
	// decode string to entry method
	Decode(string) error
}

type treeNode interface {
	// get parent id method
	GetParentID() uint64
	// set parent id method
	SetParentID(uint64)
}
