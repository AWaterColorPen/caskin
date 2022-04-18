package caskin

// idInterface it is only for get/set id method
type idInterface interface {
	// GetID get id method
	GetID() uint64
	// SetID set id method
	SetID(uint64)
}

func isValid(item idInterface) error {
	if item == nil {
		return ErrNil
	}
	if item.GetID() == 0 {
		return ErrEmptyID
	}
	return nil
}

// nameInterface it is only for get/set name method
type nameInterface interface {
	// GetName get name method
	GetName() string
	// SetName set name method
	SetName(string)
}

// domainInterface it is only for get/set domain_id method
type domainInterface interface {
	// GetDomainID get domain_id method
	GetDomainID() uint64
	// SetDomainID set domain_id method
	SetDomainID(uint64)
}

// codeInterface it is only for encode/decode method
type codeInterface interface {
	// Encode to string method
	Encode() string
	// Decode string to instance method
	Decode(string) error
}

// parentInterface it is only for get/set parent_id method
type parentInterface interface {
	// GetParentID get parent id method
	GetParentID() uint64
	// SetParentID set parent id method
	SetParentID(uint64)
}
