package web_feature

import (
	"time"

	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

// Feature
type Feature struct {
	ID          uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	Name        string         `gorm:"column:name;unique"     json:"name,omitempty"`
	Description string         `gorm:"column:description"     json:"description,omitempty"`
	Group       string         `gorm:"column:group"           json:"group,omitempty"`
	DomainID    uint64         `gorm:"column:domain_id"       json:"domain_id,omitempty"`
	ObjectID    uint64         `gorm:"column:object_id"       json:"object_id,omitempty"`
	ParentID    uint64         `gorm:"-"                      json:"parent_id,omitempty"`
}

const (
	ObjectTypeFeature  caskin.ObjectType = "feature"
	ObjectTypeFrontend caskin.ObjectType = "frontend"
	ObjectTypeBackend  caskin.ObjectType = "backend"
)
