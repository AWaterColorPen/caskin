package example

import (
	"time"

	"gorm.io/gorm"
)

// OneObjectData one kind of object_data
type OneObjectData struct {
	ID        uint64         `gorm:"column:id;primaryKey"                              json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"                                 json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"                                 json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"                            json:"-"`
	Name      string         `gorm:"column:name;index:idx_one_object_data,unique"      json:"name,omitempty"`
	DomainID  uint64         `gorm:"column:domain_id;index:idx_one_object_data,unique" json:"domain_id,omitempty"`
	ObjectID  uint64         `gorm:"column:object_id"                                  json:"object_id,omitempty"`
}
