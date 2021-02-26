package example

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// User sample for caskin.User interface
type User struct {
	ID          uint64         `gorm:"column:id;primaryKey"       json:"id,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at"          json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"          json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:delete_at;index"     json:"-"`
	PhoneNumber string         `gorm:"column:phone_number;unique" json:"phone_number,omitempty"`
	Email       string         `gorm:"column:email;unique"        json:"email,omitempty"`
}

func (u *User) GetID() uint64 {
	return u.ID
}

func (u *User) SetID(id uint64) {
	u.ID = id
}

func (u *User) Encode() string {
	return fmt.Sprintf("user_%v", u.ID)
}

func (u *User) Decode(code string) error {
	_, err := fmt.Sscanf(code, "user_%v", &u.ID)
	return err
}

func (u *User) IsObject() bool {
	return false
}

func (u *User) GetObject() string {
	return ""
}
