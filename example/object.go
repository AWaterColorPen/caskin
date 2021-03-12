package example

import (
	"fmt"
	"time"

	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

// Object sample for caskin.Object interface
type Object struct {
	ID        uint64            `gorm:"column:id;primaryKey"                     json:"id,omitempty"`
	CreatedAt time.Time         `gorm:"column:created_at"                        json:"created_at,omitempty"`
	UpdatedAt time.Time         `gorm:"column:updated_at"                        json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt    `gorm:"column:delete_at;index"                   json:"-"`
	Name      string            `gorm:"column:name;index:idx_object,unique"      json:"name,omitempty"`
	Type      caskin.ObjectType `gorm:"column:type"                              json:"type,omitempty"`
	DomainID  uint64            `gorm:"column:domain_id;index:idx_object,unique" json:"domain_id,omitempty"`
	ObjectID  uint64            `gorm:"column:object_id"                         json:"object_id,omitempty"`
	ParentID  uint64            `gorm:"column:parent_id"                         json:"parent_id"`
}

func (o *Object) GetID() uint64 {
	return o.ID
}

func (o *Object) SetID(id uint64) {
	o.ID = id
}

func (o *Object) Encode() string {
	return fmt.Sprintf("object_%v", o.ID)
}

func (o *Object) Decode(code string) error {
	_, err := fmt.Sscanf(code, "object_%v", &o.ID)
	return err
}

func (o *Object) GetObject() caskin.Object {
	return &Object{ID: o.ObjectID}
}

func (o *Object) SetObjectId(objectId uint64) {
	o.ObjectID = objectId
}

func (o *Object) GetParentID() uint64 {
	return o.ParentID
}

func (o *Object) SetParentID(pid uint64) {
	o.ParentID = pid
}

func (o *Object) SetDomainID(did uint64) {
	o.DomainID = did
}

func (o *Object) GetObjectType() caskin.ObjectType {
	return o.Type
}

func (o *Object) SetObjectType(objectType caskin.ObjectType) {
	o.Type = objectType
}
