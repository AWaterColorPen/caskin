package db

import (
	"time"

	"github.com/awatercolorpen/caskin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	DefaultDomainCreatorObjectTableName = "caskin_domain_creator_objects"
	DefaultDomainCreatorRoleTableName   = "caskin_domain_creator_roles"
)

type DomainCreatorObject struct {
	ID             uint64            `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt      time.Time         `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt      time.Time         `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt      gorm.DeletedAt    `gorm:"column:delete_at;index" json:"-"`
	Name           string            `gorm:"column:name;unique"     json:"name,omitempty"`
	Type           caskin.ObjectType `gorm:"column:type"            json:"type,omitempty"`
	ObjectID       uint64            `gorm:"column:object_id"       json:"object_id,omitempty"`
	ParentID       uint64            `gorm:"column:parent_id"       json:"parent_id"`
	CustomizedData datatypes.JSON    `gorm:"column:customized_data" json:"customized_data"`
}

func (d *DomainCreatorObject) TableName() string {
	return DefaultDomainCreatorObjectTableName
}

type DomainCreatorRole struct {
	ID        uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	Name      string         `gorm:"column:name;unique"     json:"name,omitempty"`
	ObjectID  uint64         `gorm:"column:object_id"       json:"object_id,omitempty"`
	ParentID  uint64         `gorm:"column:parent_id"       json:"parent_id"`
}

func (d *DomainCreatorRole) TableName() string {
	return DefaultDomainCreatorRoleTableName
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

type MetaDB interface {
	AutoMigrate(...interface{}) error

	Create(interface{}) error
	Recover(interface{}) error
	Update(interface{}) error
	Upsert(interface{}) error
	Take(interface{}) error
	TakeUnscoped(interface{}) error
	Find(items interface{}, cond ...interface{}) error
	DeleteByID(interface{}, uint64) error
}


