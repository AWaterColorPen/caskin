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
	Object    string         `gorm:"column:object"                          json:"object,omitempty"`
	DomainID  uint64         `gorm:"column:tenant_id;index:idx_role,unique" json:"tenant_id,omitempty"`
	ParentID  uint64         `gorm:"-"                                      json:"parent_id"`
}

func (r *Role) GetID() uint64 {
	return r.ID
}

func (r *Role) Encode() string {
	return fmt.Sprintf("role_%v", r.ID)
}

func (r *Role) Decode(code string) error {
	_, err := fmt.Sscanf(code, "role_%v", &r.ID)
	return err
}

func (r *Role) IsObject() bool {
	return true
}

func (r *Role) GetObject() string {
	return r.Object
}

func (r *Role) GetParentID() uint64 {
	return r.ParentID
}

func (r *Role) SetParentID(pid uint64) {
	r.ParentID = pid
}
