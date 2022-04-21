package example

import (
	"time"

	"gorm.io/gorm"
)

// OneObjectData one kind of object_data
type OneObjectData struct {
	ID        uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	Name      string         `gorm:"column:name"            json:"name,omitempty"`
	DomainID  uint64         `gorm:"column:domain_id"       json:"domain_id,omitempty"`
	ObjectID  uint64         `gorm:"column:object_id"       json:"object_id,omitempty"`
}

func (o *OneObjectData) GetID() uint64 {
	return o.ID
}

func (o *OneObjectData) SetID(id uint64) {
	o.ID = id
}

func (o *OneObjectData) GetObjectID() uint64 {
	return o.ObjectID
}

func (o *OneObjectData) SetObjectID(oid uint64) {
	o.ObjectID = oid
}

func (o *OneObjectData) GetDomainID() uint64 {
	return o.DomainID
}

func (o *OneObjectData) SetDomainID(did uint64) {
	o.DomainID = did
}
