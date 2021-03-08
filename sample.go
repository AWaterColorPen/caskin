package caskin

type sampleSuperadminRole struct {
}

func (s *sampleSuperadminRole) GetID() uint64 {
	return DefaultSuperadminRoleID
}

func (s *sampleSuperadminRole) SetID(uint64) {
}

func (s *sampleSuperadminRole) Encode() string {
	return SuperadminRole
}

func (s *sampleSuperadminRole) Decode(code string) error {
	if code != SuperadminRole {
		return ErrIsNotSuperAdmin
	}
	return nil
}

func (s *sampleSuperadminRole) IsObject() bool {
	return false
}

func (s *sampleSuperadminRole) GetObject() Object {
	return nil
}

func (s *sampleSuperadminRole) SetObjectId(objectId uint64) {
}

func (s *sampleSuperadminRole) GetParentID() uint64 {
	return 0
}

func (s *sampleSuperadminRole) SetParentID(uint64) {
}

func (s *sampleSuperadminRole) SetDomainID(uint64) {
}

type sampleSuperAdminDomain struct {
}

func (s *sampleSuperAdminDomain) GetID() uint64 {
	return DefaultSuperadminDomainID
}

func (s *sampleSuperAdminDomain) SetID(uint64) {
}

func (s *sampleSuperAdminDomain) Encode() string {
	return SuperadminDomain
}

func (s *sampleSuperAdminDomain) Decode(code string) error {
	if code != SuperadminDomain {
		return ErrIsNotSuperAdmin
	}
	return nil
}

func (s *sampleSuperAdminDomain) IsObject() bool {
	return false
}

func (s *sampleSuperAdminDomain) GetObject() Object {
	return nil
}

func (s *sampleSuperAdminDomain) SetObjectId(u uint64) {
	panic("implement me")
}
