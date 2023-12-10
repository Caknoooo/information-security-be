package dto

import (
	"errors"
	"mime/multipart"

	"github.com/Caknoooo/golang-clean_template/entities"
)

var (
	ErrInvalidDigitalSignature = errors.New("invalid digital signature")
	ErrPdfFileDifferentContent = errors.New("the PDF files have different content")
)

const (
	MESSAGE_FAILED_CREATE_DIGITAL_SIGNATURE = "failed to create digital signature"
	MESSAGE_FAILED_VERIFY_DIGITAL_SIGNATURE = "failed to verify digital signature"
	MESSAGE_FAILED_GET_ALL_NOTIFICATIONS    = "failed to get all notifications"

	MESSAGE_SUCCESS_CREATE_DIGITAL_SIGNATURE = "success to create digital signature"
	MESSAGE_SUCCESS_VERIFY_DIGITAL_SIGNATURE = "success to verify digital signature"
	MESSAGE_SUCCESS_GET_ALL_NOTIFICATIONS    = "success to get all notifications"
)

type (
	DigitalSignatureRequest struct {
		From        string                `json:"from" form:"from"`
		To          string                `json:"to" form:"to"`
		Subject     string                `json:"subject" form:"subject"`
		BodyContent string                `json:"body_content" form:"body_content"`
		BodyFiles   *multipart.FileHeader `json:"body_files" form:"body_files"`
	}

	BodyRequest struct {
		Content string                `json:"content" form:"content"`
		Files   *multipart.FileHeader `json:"files" form:"files" binding:"required"`
	}

	DigitalSignatureResponse struct {
		ID         string `json:"id"`
		SenderID   string `json:"sender_id"`
		ReceiverID string `json:"receiver_id"`
		Subject    string `json:"subject"`
		Content    string `json:"content"`
		Filepath   string `json:"filepath"`
		IsSigned   bool   `json:"is_signed"`
	}

	GetAllNotificationsResponse struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Subject     string `json:"subject"`
		BodyContent string `json:"body_content"`
		Filepath    string `json:"filepath"`
	}

	GetAllNotificationsWithPaginationResponse struct {
		Notifications []GetAllNotificationsResponse `json:"notifications"`
		PaginationResponse
	}

	GetAllNotificationsRepository struct {
		DigitalSignature []entities.DigitalSignature
		PaginationResponse
	}

	VerifyDigitalSignatureRequest struct {
		UserId string                `json:"user_id" form:"user_id"`
		Files  *multipart.FileHeader `json:"files" form:"files" binding:"required"`
	}

	VerifyDigitalSignatureResponse struct {
		IsVerified bool `json:"is_verified"`
		SenderVerifyResponse
	}

	SenderVerifyResponse struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Date  string `json:"date"`
	}
)
