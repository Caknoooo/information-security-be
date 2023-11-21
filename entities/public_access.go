package entities

import "github.com/google/uuid"

type PublicAccess struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PublicSymmetricKey string    `json:"public_symmetric_key"`
	RequesterId        uuid.UUID `json:"requester_id" gorm:"foreignKey"`
	Requester          *User     `json:"requester" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OwnerId            uuid.UUID `json:"owner_id" gorm:"foreignKey"`
	Owner              *User     `json:"owner" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
