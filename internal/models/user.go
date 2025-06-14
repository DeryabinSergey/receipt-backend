package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        []byte         `gorm:"type:binary(16);primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	GoogleID  *uint64        `gorm:"uniqueIndex" json:"google_id,omitempty"`
}

// BeforeCreate generates UUID v7 for new users
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if len(u.ID) == 0 {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.ID = id[:]
	}
	return nil
}

// GetUUID returns the UUID as string
func (u *User) GetUUID() string {
	if len(u.ID) == 16 {
		id, _ := uuid.FromBytes(u.ID)
		return id.String()
	}
	return ""
}

// SetUUID sets the UUID from string
func (u *User) SetUUID(uuidStr string) error {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return err
	}
	u.ID = id[:]
	return nil
}