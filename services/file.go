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
		GetAllFileByUser(ctx context.Context, userId string) ([]entities.File, error)
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

	uploadFile := entities.File{
		ID:       fileId,
		Url:      fileName,
		FileName: req.File.Filename,
		UserId:   uuid.MustParse(userId),
	}

	_, err := s.fileRepo.Create(ctx, nil, uploadFile)
	if err != nil {
		return dto.UploadFileResponse{}, err
	}

	return dto.UploadFileResponse{
		Url:      fileName,
		Filename: req.File.Filename,
	}, nil
}

func (s *fileService) GetAllFileByUser(ctx context.Context, userId string) ([]entities.File, error) {
	result, err := s.fileRepo.GetAllFileByUserId(ctx, nil, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}
