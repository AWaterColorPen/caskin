package caskin

type sampleSuperAdminRole struct {
}

func (s *sampleSuperAdminRole) GetID() uint64 {
	return DefaultSuperAdminRoleID
}

func (s *sampleSuperAdminRole) Encode() string {
	return SuperAdminRole
}

func (s *sampleSuperAdminRole) Decode(code string) error {
	if code != SuperAdminRole {
		return ErrIsNotSuperAdmin
	}
	return nil
}

func (s *sampleSuperAdminRole) IsObject() bool {
	return false
}

func (s *sampleSuperAdminRole) GetObject() string {
	return ""
}

func (s *sampleSuperAdminRole) GetParentID() uint64 {
	return 0
}

func (s *sampleSuperAdminRole) SetParentID(uint64) {
}

type sampleSuperAdminDomain struct {
}

func (s *sampleSuperAdminDomain) GetID() uint64 {
	return DefaultSuperAdminDomainID
}

func (s *sampleSuperAdminDomain) Encode() string {
	return SuperAdminDomain
}

func (s *sampleSuperAdminDomain) Decode(code string) error {
	if code != SuperAdminDomain {
		return ErrIsNotSuperAdmin
	}
	return nil
}

func (s *sampleSuperAdminDomain) IsObject() bool {
	return false
}

func (s *sampleSuperAdminDomain) GetObject() string {
	return ""
}
