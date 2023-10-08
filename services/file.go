package services

import (
	"context"
	"fmt"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/entities"
	"github.com/Caknoooo/golang-clean_template/repository"
	"github.com/Caknoooo/golang-clean_template/utils"
	"github.com/google/uuid"
)

type (
	FileService interface {
		UploadFile(ctx context.Context, req dto.UploadFileRequest, userId string) (dto.UploadFileResponse, error)
		GetAllFileByUser(ctx context.Context, userId string) ([]dto.UploadFileResponse, error)
		DecryptFileData(ctx context.Context, encryption string) (string, error)
	}

	fileService struct {
		fileRepo repository.FileRepository
	}
)

func NewFileService(fileRepo repository.FileRepository) FileService {
	return &fileService{
		fileRepo: fileRepo,
	}
}

const (
	PATH  = "storage"
	FILES = "files"
)

func (s *fileService) UploadFile(ctx context.Context, req dto.UploadFileRequest, userId string) (dto.UploadFileResponse, error) {
	fileId := uuid.New()

	fileName := fmt.Sprintf("%s/files/%s", userId, fileId)
	if err := utils.UploadFileSuccess(req.File, fileName); err != nil {
		return dto.UploadFileResponse{}, err
	}

	encryption, data, err := utils.AESEncrypt(fileName, utils.FILE_KEY)
	if err != nil {
		return dto.UploadFileResponse{}, err
	}

	uploadFile := entities.File{
		ID:         fileId,
		Path:       fileName,
		Encryption: encryption,
		FileType:   req.FileType,
		FileName:   req.File.Filename,
		UserId:     uuid.MustParse(userId),
	}

	_, err = s.fileRepo.Create(ctx, nil, uploadFile)
	if err != nil {
		return dto.UploadFileResponse{}, err
	}

	return dto.UploadFileResponse{
		Path:             fileName,
		Filename:         req.File.Filename,
		FileType:         req.FileType,
		Encryption:       encryption,
		AES_KEY:          data["key"].(string),
		AES_PLAIN_TEXT:   data["plaintext"].(string),
		AES_BLOCK_CHIPER: data["block"].(string),
		AES_GCM:          data["aes-gcm"].(string),
		AES_NONCE:        data["nonce"].(string),
		AES_RESULT:       encryption,
	}, nil
}

func (s *fileService) GetAllFileByUser(ctx context.Context, userId string) ([]dto.UploadFileResponse, error) {
	result, err := s.fileRepo.GetAllFileByUserId(ctx, nil, userId)
	if err != nil {
		return nil, err
	}

	var files []dto.UploadFileResponse
	for _, file := range result {
		files = append(files, dto.UploadFileResponse{
			Path:       file.Path,
			Filename:   file.FileName,
			FileType:   file.FileType,
			Encryption: file.Encryption,
		})
	}

	return files, nil
}

func (s *fileService) DecryptFileData(ctx context.Context, encryption string) (string, error) {
	result, err := utils.AESDecrypt(encryption, utils.FILE_KEY)
	if err != nil {
		return "", err
	}

	return result, nil
}
