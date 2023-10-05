package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/entities"
	"github.com/Caknoooo/golang-clean_template/repository"
	"github.com/google/uuid"
)

type (
	FileService interface {
		UploadFile(ctx context.Context, req dto.UploadFileRequest, userId string) (dto.UploadFileResponse, error)
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
	fileName := fmt.Sprintf("%s/%s/%s/%s", PATH, FILES, userId, fileId)

	// Ensure the parent directory exists
	if err := os.MkdirAll(filepath.Join(PATH, FILES, userId), 0666); err != nil {
		return dto.UploadFileResponse{}, err
	}

	out, err := os.Create(fileName)
	if err != nil {
		return dto.UploadFileResponse{}, err
	}
	defer out.Close()

	uploadFile := entities.File{
		ID:       fileId,
		Url:      fileName,
		FileName: req.File.Filename,
		UserId:   uuid.MustParse(userId),
	}

	_, err = s.fileRepo.Create(ctx, nil, uploadFile)
	if err != nil {
		return dto.UploadFileResponse{}, err
	}

	uploadedFile, err := req.File.Open()
	if err != nil {
		return dto.UploadFileResponse{}, err
	}
	defer uploadedFile.Close()

	_, err = io.Copy(out, uploadedFile)
	if err != nil {
		return dto.UploadFileResponse{}, err
	}

	return dto.UploadFileResponse{
		Url:      fileName,
		Filename: req.File.Filename,
	}, nil
}

func (s *fileService) GetAllFile(ctx context.Context) ([]entities.File, error) {
	file, err := s.fileRepo.GetAllFile(ctx, nil)
	if err != nil {
		return nil, err
	}
	return file, nil
}
