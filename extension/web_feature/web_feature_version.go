package web_feature

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	MetaData  Relations      `gorm:"column:metadata"        json:"metadata,omitempty"`
}

type Relations caskin.InheritanceRelations

// Scan scan value into Jsonb, implements sql.Scanner interface
func (r *Relations) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, r)
}

// Value return json value, implement driver.Valuer interface
func (r Relations) Value() (driver.Value, error) {
	bytes, err := json.Marshal(r)
	return string(bytes), err
}

func (r Relations) Version() string {
	h := sha256.New()
	b, _ := json.Marshal(r)
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

type FeatureRelation = caskin.InheritanceRelation
type FeatureRelations = caskin.InheritanceRelations

type VersionedDomain interface {
	caskin.Domain
	GetVersion() string
	SetVersion(string)
}
