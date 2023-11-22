package entities

import "github.com/google/uuid"

type (
	File struct {
		ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Path           string    `json:"path"`
		FileName       string    `json:"file_name"`
		Encryption     string    `json:"encryption"`
		EncryptionMode string    `json:"mode"`
		FileType       string    `json:"file_type"`
		UserId         uuid.UUID `json:"user_id"`
		User           *User     `json:"user" gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

		Timestamp
	}
)
