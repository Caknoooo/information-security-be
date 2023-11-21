package dto

import (
	"errors"
	"mime/multipart"
)

const (
	MESSAGE_FAILED_UPLOAD_FILE  = "failed upload file"
	MESSAGE_FAILED_GET_ALL_FILE = "failed get all file"
	MESSAGE_FAILED_GET_FILE     = "failed get file"

	MESSAGE_SUCCESS_UPLOAD_FILE  = "success upload file"
	MESSAGE_SUCCESS_GET_ALL_FILE = "success get all file"
)

var (
	ErrKeyInvalid = errors.New("invalid key")
)

type (
	UploadFileRequest struct {
		File     *multipart.FileHeader `form:"file" binding:"required"`
		FileType string                `form:"file_type"`
	}

	UploadFileResponse struct {
		Path             string `json:"path"`
		Filename         string `json:"file_name"`
		FileType         string `json:"file_type"`
		Encryption       string `json:"encryption"`
		EncryptionMode   string `json:"mode"`
		AES_KEY          string `json:"aes_key,omitempty"`
		AES_PLAIN_TEXT   string `json:"aes_plain_text,omitempty"`
		AES_BLOCK_CHIPER string `json:"aes_block_chiper,omitempty"`
		AES_GCM          string `json:"aes_gcm,omitempty"`
		AES_CIPHERTEXT   string `json:"aes_ciphertext,omitempty"`
		AES_NONCE        string `json:"aes_nonce,omitempty"`
		AES_RESULT       string `json:"aes_result,omitempty"`
		ELAPSEDTIME      string `json:"elapsed_time,omitempty"`
	}
)
