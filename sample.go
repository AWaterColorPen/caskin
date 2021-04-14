package caskin

import (
	"gorm.io/datatypes"
)

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

type SampleNoPermissionObject struct {
}

func (s *SampleNoPermissionObject) GetID() uint64 {
	return 0
}

func (s *SampleNoPermissionObject) SetID(uint64) {
}

func (s *SampleNoPermissionObject) Encode() string {
	return DefaultNoPermissionObject
}

func (s *SampleNoPermissionObject) Decode(code string) error {
	if code != DefaultNoPermissionObject {
		return ErrInValidObject
	}
	return nil
}

func (s *SampleNoPermissionObject) GetParentID() uint64 {
	return 0
}

func (s *SampleNoPermissionObject) SetParentID(uint64) {
}

func (s *SampleNoPermissionObject) GetObject() Object {
	return nil
}

func (s *SampleNoPermissionObject) SetObjectID(uint64) {
}

func (s *SampleNoPermissionObject) GetDomainID() uint64 {
	return 0
}

func (s *SampleNoPermissionObject) SetDomainID(uint64) {
}

func (s *SampleNoPermissionObject) GetName() string {
	return DefaultNoPermissionObject
}

func (s *SampleNoPermissionObject) SetName(string) {
}

func (s *SampleNoPermissionObject) GetObjectType() ObjectType {
	return ObjectTypeDefault
}

func (s *SampleNoPermissionObject) SetObjectType(ObjectType) {
}

func (s *SampleNoPermissionObject) GetCustomizedData() datatypes.JSON {
	return nil
}

func (s *SampleNoPermissionObject) SetCustomizedData(datatypes.JSON) {
}
