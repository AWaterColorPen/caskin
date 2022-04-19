package caskin

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

func GetSuperadminRole() Role {
	return &SampleSuperadminRole{
		ID:   0,
		Name: DefaultSuperadminRoleName,
	}
}

func GetSuperadminDomain() Domain {
	return &SampleSuperadminDomain{
		ID:   0,
		Name: DefaultSuperadminDomainName,
	}
}
