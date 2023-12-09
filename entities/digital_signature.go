package entities

import "github.com/google/uuid"

type (
	DigitalSignature struct {
		ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
		SenderID   string    `json:"sender_id"`
		ReceiverID string    `json:"receiver_id"`
		Subject    string    `json:"subject"`
		Content    string    `json:"content"`
		Filepath   string    `json:"filepath"`
		IsSigned   bool      `gprm:"default:false" json:"is_signed"`
	}
)
