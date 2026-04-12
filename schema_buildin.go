package caskin

// NamedObject is the built-in [Object] implementation used for special objects
// such as the superadmin sentinel. Its Encode/Decode methods simply use the
// Name string, so it can represent objects that exist only in casbin policy
// strings without corresponding database rows.
type NamedObject struct {
	Name string `json:"name"`
}

func (o *NamedObject) GetID() uint64 {
	return 0
}

func (o *NamedObject) SetID(uint64) {
}

func (o *NamedObject) Encode() string {
	return o.Name
}

func (o *NamedObject) Decode(code string) error {
	o.Name = code
	return nil
}

func (o *NamedObject) GetParentID() uint64 {
	return 0
}

func (o *NamedObject) SetParentID(uint64) {
}

func (o *NamedObject) GetDomainID() uint64 {
	return 0
}

func (o *NamedObject) SetDomainID(uint64) {
}

func (o *NamedObject) GetObjectType() string {
	return ""
}

// SampleSuperadminRole is the built-in [Role] that identifies the global
// superadmin. Its Encode method always returns the [SuperadminRole] constant
// and Decode accepts only that value, returning [ErrIsNotSuperadmin] otherwise.
type SampleSuperadminRole struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (s *SampleSuperadminRole) GetID() uint64 {
	return s.ID
}

func (s *SampleSuperadminRole) SetID(uint64) {
}

func (s *SampleSuperadminRole) Encode() string {
	return SuperadminRole
}

func (s *SampleSuperadminRole) Decode(code string) error {
	if code != SuperadminRole {
		return ErrIsNotSuperadmin
	}
	return nil
}

func (s *SampleSuperadminRole) GetObjectID() uint64 {
	return 0
}

func (s *SampleSuperadminRole) SetObjectID(uint64) {
}

func (s *SampleSuperadminRole) GetDomainID() uint64 {
	return 0
}

func (s *SampleSuperadminRole) SetDomainID(uint64) {
}

// SampleSuperadminDomain is the built-in [Domain] that represents the
// superadmin scope. Its Encode method always returns [SuperadminDomain] and
// Decode accepts only that value.
type SampleSuperadminDomain struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (s *SampleSuperadminDomain) GetID() uint64 {
	return s.ID
}

func (s *SampleSuperadminDomain) SetID(uint64) {
}

func (s *SampleSuperadminDomain) Encode() string {
	return SuperadminDomain
}

func (s *SampleSuperadminDomain) Decode(code string) error {
	if code != SuperadminDomain {
		return ErrIsNotSuperadmin
	}
	return nil
}

// GetSuperadminRole returns a [Role] that encodes to the [SuperadminRole]
// constant. It is used internally when granting or checking superadmin status.
func GetSuperadminRole() Role {
	return &SampleSuperadminRole{
		ID:   0,
		Name: DefaultSuperadminRoleName,
	}
}

// GetSuperadminDomain returns a [Domain] that encodes to the
// [SuperadminDomain] constant. It is used internally alongside
// [GetSuperadminRole] for superadmin enforcement.
func GetSuperadminDomain() Domain {
	return &SampleSuperadminDomain{
		ID:   0,
		Name: DefaultSuperadminDomainName,
	}
}
