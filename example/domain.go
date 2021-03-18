package example

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Domain sample for caskin.Domain interface
type Domain struct {
	ID        uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	Name      string         `gorm:"column:name;unique"     json:"name,omitempty"`
	Meta      *DomainMeta    `gorm:"column:meta"            json:"meta,omitempty"`
}

func (d *Domain) GetID() uint64 {
	return d.ID
}

func (d *Domain) SetID(id uint64) {
	d.ID = id
}

func (d *Domain) Encode() string {
	return fmt.Sprintf("domain_%v", d.ID)
}

func (d *Domain) Decode(code string) error {
	_, err := fmt.Sscanf(code, "domain_%v", &d.ID)
	return err
}

func (d *Domain) GetVersion() string {
	if d.Meta == nil {
		return ""
	}
	return d.Meta.WebFeatureVersion
}

func (d *Domain) SetVersion(version string) {
	if d.Meta == nil {
		d.Meta = &DomainMeta{}
	}
	d.Meta.WebFeatureVersion = version
}

type DomainMeta struct {
	WebFeatureVersion string `json:"web_feature_version,omitempty"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (d *DomainMeta) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, d)
}

// Value return json value, implement driver.Valuer interface
func (d DomainMeta) Value() (driver.Value, error) {
	bytes, err := json.Marshal(d)
	return string(bytes), err
}
