package caskin

type SampleSuperadminRole struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (s *SampleSuperadminRole) GetID() uint64 {
	return DefaultSuperadminRoleID
}

func (s *SampleSuperadminRole) SetID(uint64) {
}

func (s *SampleSuperadminRole) Encode() string {
	return SuperadminRole
}

func (s *SampleSuperadminRole) Decode(code string) error {
	if code != SuperadminRole {
		return ErrIsNotSuperAdmin
	}
	return nil
}

func (s *SampleSuperadminRole) GetObject() Object {
	return nil
}

func (s *SampleSuperadminRole) SetObjectID(uint64) {
}

func (s *SampleSuperadminRole) GetParentID() uint64 {
	return 0
}

func (s *SampleSuperadminRole) SetParentID(uint64) {
}

func (s *SampleSuperadminRole) GetDomainID() uint64 {
	return DefaultSuperadminDomainID
}

func (s *SampleSuperadminRole) SetDomainID(uint64) {
}

func (s *SampleSuperadminRole) GetName() string {
	return DefaultSuperadminRoleName
}

func (s *SampleSuperadminRole) SetName(string) {
}

type SampleSuperadminDomain struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (s *SampleSuperadminDomain) GetID() uint64 {
	return DefaultSuperadminDomainID
}

func (s *SampleSuperadminDomain) SetID(uint64) {
}

func (s *SampleSuperadminDomain) Encode() string {
	return SuperadminDomain
}

func (s *SampleSuperadminDomain) Decode(code string) error {
	if code != SuperadminDomain {
		return ErrIsNotSuperAdmin
	}
	return nil
}
