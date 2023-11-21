package dto

import "errors"

const (
	MESSAGE_FAILED_CREATE_PRIVATE_ACCESS          = "failed create private access"
	MESSAGE_FAILED_GET_ALL_PRIVATE_ACCESS_REQUEST = "failed get all private access request"
	MESSAGE_FAILED_GET_ALL_PRIVATE_ACCESS_OWNER   = "failed get all private access owner"
	MESSAGE_FAILED_UPDATE_PRIVATE_ACCESS          = "failed update private access"
	MESSAGE_FAILED_SEND_ENCRYPTION_KEY            = "failed send encryption key"

	MESSAGE_SUCCESS_CREATE_PRIVATE_ACCESS          = "success create private access"
	MESSAGE_SUCCESS_GET_ALL_PRIVATE_ACCESS_REQUEST = "success get all private access request"
	MESSAGE_SUCCESS_GET_ALL_PRIVATE_ACCESS_OWNER   = "success get all private access owner"
	MESSAGE_SUCCESS_UPDATE_PRIVATE_ACCESS          = "success update private access"
	MESSAGE_SUCCESS_SEND_ENCRYPTION_KEY            = "success send encryption key"
)

var (
	ErrCreatePrivateAccess = errors.New("failed to create private access")
	ErrStatusNotFound      = errors.New("status not found")
)

type (
	PrivateAccessRequest struct {
		UserReqId   string `json:"user_req_id"`
		UserOwnerId string `json:"user_owner_id"`
	}

	PrivateAccessResponse struct {
		ID          string `json:"id"`
		UserReqId   string `json:"user_req_id"`
		UserOwnerId string `json:"user_owner_id"`
		Status      string `json:"status"`
	}

	GetPrivateAccessResponse struct {
		ID      string `json:"id"`
		UserID  string `json:"user_id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Status  string `json:"status"`
		IsOwner bool   `json:"is_owner,omitempty"`
	}

	UpdatePrivateAccessRequest struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	UpdatePrivateAccessResponse struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	SendEncryptionKeyRequest struct {
		OwnerId string `json:"owner_id"`
		Key     string `json:"key"`
	}

	SendEncryptionKeyResponse struct {
		OwnerId string         `json:"owner_id"`
		User    map[string]any `json:"user"`
	}
)
