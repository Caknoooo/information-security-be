package entities

import "github.com/google/uuid"

type (
	DigitalSignature struct {
		ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
		SenderID   string    `json:"sender_id"`
		Sender     User      `gorm:"foreignKey:SenderID;references:ID"`
		ReceiverID string    `json:"receiver_id"`
		Receiver   User      `gorm:"foreignKey:ReceiverID;references:ID"`
		Subject    string    `json:"subject"`
		Content    string    `json:"content"`
		Filepath   string    `json:"filepath"`
		IsSigned   bool      `gprm:"default:false" json:"is_signed"`
	}
)
