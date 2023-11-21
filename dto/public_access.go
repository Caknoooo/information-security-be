package dto

type (
	PublicAccessRequest struct {
		PublicShareKey string `json:"public_share_key" form:"public_share_key" binding:"required"`
		PublicKey      string `json:"public_key" form:"public_key" binding:"required"`
		RequesterId    string `json:"requester_id" form:"requester_id"`
		OwnerId        string `json:"owner_id" form:"owner_id"`
	}
)
