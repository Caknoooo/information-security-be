package dto

import "mime/multipart"

const (
	MESSAGE_FAILED_UPLOAD_FILE  = "failed upload file"
	MESSAGE_FAILED_GET_ALL_FILE = "failed get all file"
	MESSAGE_FAILED_GET_FILE     = "failed get file"

	MESSAGE_SUCCESS_UPLOAD_FILE  = "success upload file"
	MESSAGE_SUCCESS_GET_ALL_FILE = "success get all file"
)

type (
	UploadFileRequest struct {
		File *multipart.FileHeader `form:"file" binding:"required"`
	}

	UploadFileResponse struct {
		Path             string `json:"path"`
		Filename         string `json:"file_name"`
		Encryption       string `json:"encryption"`
		AES_KEY          string `json:"aes_key,omitempty"`
		AES_PLAIN_TEXT   string `json:"aes_plain_text,omitempty"`
		AES_BLOCK_CHIPER string `json:"aes_block_chiper,omitempty"`
		AES_GCM          string `json:"aes_gcm,omitempty"`
		AES_NONCE        string `json:"aes_nonce,omitempty"`
		AES_RESULT       string `json:"aes_result,omitempty"`
	}
)
