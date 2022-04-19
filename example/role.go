package example

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Role sample for caskin.Role interface
type Role struct {
	ID        uint64         `gorm:"column:id;primaryKey"                   json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"                      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"                      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"                 json:"-"`
	Name      string         `gorm:"column:name;index:idx_role,unique"      json:"name,omitempty"`
	DomainID  uint64         `gorm:"column:domain_id;index:idx_role,unique" json:"domain_id,omitempty"`
	ObjectID  uint64         `gorm:"column:object_id"                       json:"object_id,omitempty"`
}

func (r *Role) GetID() uint64 {
	return r.ID
}

func (r *Role) SetID(id uint64) {
	r.ID = id
}

func (r *Role) Encode() string {
	return fmt.Sprintf("role_%v", r.ID)
}

func (r *Role) Decode(code string) error {
	_, err := fmt.Sscanf(code, "role_%v", &r.ID)
	return err
}

func (r *Role) GetObjectID() uint64 {
	return r.ObjectID
}

func (r *Role) SetObjectID(objectId uint64) {
	r.ObjectID = objectId
}

func (r *Role) GetDomainID() uint64 {
	return r.DomainID
}

func (r *Role) SetDomainID(did uint64) {
	r.DomainID = did
}
