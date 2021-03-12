package web_feature

import (
    "github.com/awatercolorpen/caskin"
    "time"

    "gorm.io/gorm"
)

// Backend
type Backend struct {
    ID          uint64         `gorm:"column:id;primaryKey"                   json:"id,omitempty"`
    CreatedAt   time.Time      `gorm:"column:created_at"                      json:"created_at,omitempty"`
    UpdatedAt   time.Time      `gorm:"column:updated_at"                      json:"updated_at,omitempty"`
    DeletedAt   gorm.DeletedAt `gorm:"column:delete_at;index"                 json:"-"`
    Path        string         `gorm:"column:path;index:idx_backend,unique"   json:"path,omitempty"`
    Method      string         `gorm:"column:method;index:idx_backend,unique" json:"method,omitempty"`
    Description string         `gorm:"column:description"                     json:"description,omitempty"`
    Group       string         `gorm:"column:group"                           json:"group,omitempty"`
    DomainID    uint64         `gorm:"column:domain_id"                       json:"domain_id,omitempty"`
    ObjectID    uint64         `gorm:"column:object_id"                       json:"object_id,omitempty"`
    Object      caskin.Object  `gorm:"foreignKey:ObjectID"                    json:"object"`
}