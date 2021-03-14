package web_feature

import (
	"time"

	"gorm.io/gorm"
)

type WebFeatureVersion struct {
	ID        uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	MD5       string         `gorm:"column:md5"             json:"md5,omitempty"`
	MetaData  Relation       `gorm:"column:metadata"        json:"metadata,omitempty"`
}

type Relation = []interface{}
type Relations map[interface{}]Relation

