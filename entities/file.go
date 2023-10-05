package entities

import "github.com/google/uuid"

type (
	File struct {
		ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Url    string    `json:"url"`
		FileName string `json:"file_name"`
		UserId uuid.UUID `json:"user_id"`
		User   *User      `json:"user" gorm:"foreignKey:UserId"`
	}
)
