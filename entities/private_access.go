package entities

import "github.com/google/uuid"

type (
	PrivateAccess struct {
		ID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`

		UserReqId string `json:"user_req_id"`
		UserReq   *User  `json:"user_req" gorm:"foreignKey:UserReqId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

		UserOwnerId string `json:"user_owner_id"`
		UserOwner   *User  `json:"user_owner" gorm:"foreignKey:UserOwnerId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		
		Status string `gorm:"default:pending" json:"status"`
	}
)