package web_feature

import (
	"time"

	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

// Frontend
type Frontend struct {
	ID          uint64         `gorm:"column:id;primaryKey"                  json:"id,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at"                     json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"                     json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:delete_at;index"                json:"-"`
	Key         string         `gorm:"column:key;index:idx_frontend,unique"  json:"key,omitempty"`
	Type        FrontendType   `gorm:"column:type;index:idx_frontend,unique" json:"type,omitempty"`
	Description string         `gorm:"column:description"                    json:"description,omitempty"`
	Group       string         `gorm:"column:group"                          json:"group,omitempty"`
	DomainID    uint64         `gorm:"column:domain_id"                      json:"domain_id,omitempty"`
	ObjectID    uint64         `gorm:"column:object_id"                      json:"object_id,omitempty"`
	Object      caskin.Object  `gorm:"foreignKey:ObjectID"                   json:"object"`
}

type FrontendType string

const (
	FrontendTypeMenu        FrontendType = "menu"
	FrontendTypeSubFunction FrontendType = "sub_function"
)
