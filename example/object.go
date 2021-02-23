package example

import (
	"fmt"
	"github.com/awatercolorpen/caskin"
	"time"

	"gorm.io/gorm"
)

// Object sample for caskin.Object interface
type Object struct {
	ID        uint64            `gorm:"column:id;primaryKey"                   json:"id,omitempty"`
	CreatedAt time.Time         `gorm:"column:created_at"                      json:"created_at,omitempty"`
	UpdatedAt time.Time         `gorm:"column:updated_at"                      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt    `gorm:"column:delete_at;index"                 json:"-"`
	Name      string            `gorm:"column:name;index:idx_role,unique"      json:"name,omitempty"`
	Type      caskin.ObjectType `gorm:"column:type"                            json:"type,omitempty"`
	Object    string            `gorm:"column:object"                          json:"object,omitempty"`
	DomainID  uint64            `gorm:"column:tenant_id;index:idx_role,unique" json:"tenant_id,omitempty"`
	ParentID  uint64            `gorm:"-"                                      json:"parent_id"`
}

const (
	ObjectTypeDefault caskin.ObjectType = "default"
	ObjectTypeObject  caskin.ObjectType = "object"
	ObjectTypeRole    caskin.ObjectType = "role"
)

func (o *Object) GetID() uint64 {
	return o.ID
}

func (o *Object) Encode() string {
	return fmt.Sprintf("object_%v", o.ID)
}

func (o *Object) Decode(code string) error {
	_, err := fmt.Sscanf(code, "object_%v", &o.ID)
	return err
}

func (o *Object) IsObject() bool {
	return true
}

func (o *Object) GetObject() string {
	return o.Object
}

func (o *Object) GetParentID() uint64 {
	return o.ParentID
}

func (o *Object) SetParentID(pid uint64) {
	o.ParentID = pid
}

func (o *Object) GetObjectType() caskin.ObjectType {
	return o.Type
}
