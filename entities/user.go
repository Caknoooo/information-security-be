package entities

import (
	"time"

	"github.com/Caknoooo/golang-clean_template/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name               string    `json:"name"`
	TelpNumber         string    `json:"telp_number"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	Role               string    `json:"role"`
	Work               string    `json:"work"`
	SymmetricKey       string    `json:"symmetric_key"`
	PublicSymmetricKey string    `json:"public_symmetric_key"`
	PublicKey          string    `json:"public_key"`
	PrivateKey         string    `json:"private_key"`
	ActivePeriod       time.Time `json:"active_period"`
	IsVerified         bool      `json:"is_verified"`

	Timestamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
}
