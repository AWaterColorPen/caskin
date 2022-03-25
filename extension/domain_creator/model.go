package domain_creator

import (
	"fmt"
	"time"

	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

var (
	DefaultDomainCreatorObjectTableName        = "caskin_domain_creator_objects"
	DefaultDomainCreatorRoleTableName          = "caskin_domain_creator_roles"
	DefaultDomainCreatorPolicyTableName        = "caskin_domain_creator_policies"
	ErrRelativeIDAndAbsoluteRoleIDInCompatible = fmt.Errorf("relative_id and absolute_id is in compatible")
	ErrNotSupport                              = fmt.Errorf("not sopport")
	ErrRelativeIDOutOfIndex                    = fmt.Errorf("relative_id out of index")
)

type DomainCreatorObject struct {
	Name             string            `gorm:"column:name;unique"        json:"name,omitempty"`
	Type             caskin.ObjectType `gorm:"column:type"               json:"type,omitempty"`
	RelativeObjectID uint64            `gorm:"column:relative_object_id" json:"relative_object_id,omitempty"`
	RelativeParentID uint64            `gorm:"column:relative_parent_id" json:"relative_parent_id"`
	AbsoluteObjectID uint64            `gorm:"column:absolute_object_id" json:"absolute_object_id,omitempty"`
	AbsoluteParentID uint64            `gorm:"column:absolute_parent_id" json:"absolute_parent_id,omitempty"`
}

func (d *DomainCreatorObject) TableName() string {
	return DefaultDomainCreatorObjectTableName
}

func (d *DomainCreatorObject) IsValid() error {
	if d.RelativeObjectID == 0 && d.AbsoluteObjectID == 0 {
		return ErrRelativeIDAndAbsoluteRoleIDInCompatible
	}
	if d.RelativeObjectID != 0 && d.AbsoluteObjectID != 0 {
		return ErrRelativeIDAndAbsoluteRoleIDInCompatible
	}
	if d.RelativeParentID != 0 && d.AbsoluteParentID != 0 {
		return ErrRelativeIDAndAbsoluteRoleIDInCompatible
	}
	if d.RelativeParentID != 0 || d.AbsoluteParentID != 0 {
		return ErrNotSupport
	}
	return nil
}

type DomainCreatorRole struct {
	ID               uint64         `gorm:"column:id;primaryKey"      json:"id,omitempty"`
	CreatedAt        time.Time      `gorm:"column:created_at"         json:"created_at,omitempty"`
	UpdatedAt        time.Time      `gorm:"column:updated_at"         json:"updated_at,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"column:delete_at;index"    json:"-"`
	Name             string         `gorm:"column:name;unique"        json:"name,omitempty"`
	RelativeObjectID uint64         `gorm:"column:relative_object_id" json:"relative_object_id,omitempty"`
	RelativeParentID uint64         `gorm:"column:relative_parent_id" json:"relative_parent_id"`
	AbsoluteObjectID uint64         `gorm:"column:absolute_object_id" json:"absolute_object_id,omitempty"`
	AbsoluteParentID uint64         `gorm:"column:absolute_parent_id" json:"absolute_parent_id,omitempty"`
}

func (d *DomainCreatorRole) TableName() string {
	return DefaultDomainCreatorRoleTableName
}

func (d *DomainCreatorRole) IsValid() error {
	if d.RelativeObjectID == 0 && d.AbsoluteObjectID == 0 {
		return ErrRelativeIDAndAbsoluteRoleIDInCompatible
	}
	if d.RelativeObjectID != 0 && d.AbsoluteObjectID != 0 {
		return ErrRelativeIDAndAbsoluteRoleIDInCompatible
	}
	if d.RelativeParentID != 0 && d.AbsoluteParentID != 0 {
		return ErrRelativeIDAndAbsoluteRoleIDInCompatible
	}
	if d.RelativeParentID != 0 || d.AbsoluteParentID != 0 {
		return ErrNotSupport
	}
	return nil
}

type DomainCreatorPolicy struct {
	ID               uint64         `gorm:"column:id;primaryKey"      json:"id,omitempty"`
	CreatedAt        time.Time      `gorm:"column:created_at"         json:"created_at,omitempty"`
	UpdatedAt        time.Time      `gorm:"column:updated_at"         json:"updated_at,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"column:delete_at;index"    json:"-"`
	RelativeRoleID   uint64         `gorm:"column:relative_role_id"   json:"relative_role_id,omitempty"`
	RelativeObjectID uint64         `gorm:"column:relative_object_id" json:"relative_object_id,omitempty"`
	AbsoluteRoleID   uint64         `gorm:"column:absolute_role_id"   json:"absolute_role_id,omitempty"`
	AbsoluteObjectID uint64         `gorm:"column:absolute_object_id" json:"absolute_object_id,omitempty"`
	Action           caskin.Action  `gorm:"column:action"             json:"action"`
}

func (d *DomainCreatorPolicy) TableName() string {
	return DefaultDomainCreatorPolicyTableName
}

func (d *DomainCreatorPolicy) IsValid() error {
	if d.RelativeRoleID == 0 && d.AbsoluteRoleID == 0 {
		return ErrRelativeIDAndAbsoluteRoleIDInCompatible
	}
	if d.RelativeRoleID != 0 && d.AbsoluteRoleID != 0 {
		return ErrRelativeIDAndAbsoluteRoleIDInCompatible
	}
	if d.RelativeObjectID == 0 && d.AbsoluteObjectID == 0 {
		return ErrRelativeIDAndAbsoluteRoleIDInCompatible
	}
	if d.RelativeObjectID != 0 && d.AbsoluteObjectID != 0 {
		return ErrRelativeIDAndAbsoluteRoleIDInCompatible
	}
	return nil
}

type relativeIDAndAbsoluteID interface {
	IsValid() error
}
