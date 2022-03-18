package caskin

// idInterface it is only for get/set id method
type idInterface interface {
	// GetID get id method
	GetID() uint64
	// SetID set id method
	SetID(uint64)
}

// nameInterface it is only for get/set name method
type nameInterface interface {
	// GetName get name method
	GetName() string
	// SetName set name method
	SetName(string)
}

// entry it is casbin entry of User Role Object Domain
type entry interface {
	idInterface
	// Encode entry to string method
	Encode() string
	// Decode decode string to entry method
	Decode(string) error
}

type treeNode interface {
	// GetParentID get parent id method
	GetParentID() uint64
	// SetParentID set parent id method
	SetParentID(uint64)
}
