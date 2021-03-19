package web_feature

import (
	"time"

	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

type WebFeatureVersion struct {
	ID        uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	SHA256    string         `gorm:"column:sha256;unique"   json:"sha256,omitempty"`
	MetaData  *DumpRelation  `gorm:"column:metadata"        json:"metadata,omitempty"`
}

func (w *WebFeatureVersion) IsCompatible(dump *Dump) bool {
	if !isCompatible(w.MetaData.FeatureTree, dump.FeatureTree) {
		return false
	}
	if !isCompatible(w.MetaData.FrontendTree, dump.FrontendTree) {
		return false
	}
	if !isCompatible(w.MetaData.BackendTree, dump.BackendTree) {
		return false
	}
	return true
}

type Relation = caskin.InheritanceRelation
type Relations = caskin.InheritanceRelations

type VersionedDomain interface {
	caskin.Domain
	GetVersion() string
	SetVersion(string)
}
