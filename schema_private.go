package caskin

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

type treeNode interface {
	// get parent id method
	GetParentID() uint64
	// set parent id method
	SetParentID(uint64)
}

type parentEntry interface {
	entry
	treeNode
	ObjectData
}

type treeNodeEntry = parentEntry


type objectDataEntry interface {
	entry
	ObjectData
}

type idInterface interface {
	// get id method
	GetID() uint64
}

