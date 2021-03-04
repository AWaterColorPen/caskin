package example

import (
	"fmt"
	"time"

	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

// Domain sample for caskin.Domain interface
type Domain struct {
	ID        uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	Name      string         `gorm:"column:name;unique"     json:"name,omitempty"`
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

func (d *Domain) IsObject() bool {
	return false
}

func (d *Domain) GetObject() caskin.Object {
	return nil
}

func (d *Domain) SetObjectId(objcetId uint64) {
}
