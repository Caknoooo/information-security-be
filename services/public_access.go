package services

import (
	"context"
	"time"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/entities"
	"github.com/Caknoooo/golang-clean_template/repository"
	"github.com/Caknoooo/golang-clean_template/utils"
	"github.com/google/uuid"
)

type (
	PublicAccessService interface {
		PublicAccessUserFiles(ctx context.Context, req dto.PublicAccessRequest) ([]dto.UploadFileResponse, error)
	}

	publicAccessService struct {
		pbr repository.PublicAccessRepository
		fr  repository.FileRepository
		ur  repository.UserRepository
	}
)

func NewPublicAccessService(pbr repository.PublicAccessRepository, fr repository.FileRepository, ur repository.UserRepository) PublicAccessService {
	return &publicAccessService{
		pbr: pbr,
		fr:  fr,
		ur:  ur,
	}
}

func (s *publicAccessService) PublicAccessUserFiles(ctx context.Context, req dto.PublicAccessRequest) ([]dto.UploadFileResponse, error) {
	pubSymKey, err := utils.DecryptRSA(req.PublicShareKey, req.PublicKey)
	if err != nil {
		return []dto.UploadFileResponse{}, err
	}

	owner, err := s.ur.GetUserById(ctx, req.OwnerId)
	if err != nil {
		return []dto.UploadFileResponse{}, dto.ErrUserNotFound
	}

	if owner.PublicSymmetricKey != pubSymKey || owner.ActivePeriod.Before(time.Now()) {
		return []dto.UploadFileResponse{}, dto.ErrKeyInvalid
	}

	accesses, err := s.pbr.GetAllPublicAccessByIDs(ctx, nil, req.OwnerId, req.RequesterId)
	if err != nil {
		return []dto.UploadFileResponse{}, err
	}

	cnt := 0
	for _, access := range accesses {
		decryptedPubSymKey, err := utils.AESDecrypt(access.PublicSymmetricKey, utils.KEY)
		if err != nil {
			return []dto.UploadFileResponse{}, err
		}

		if decryptedPubSymKey == pubSymKey {
			cnt++
			if cnt >= 5 {
				return []dto.UploadFileResponse{}, dto.ErrKeyInvalid
			}
		}
	}

	encryptedPubSymKey, _, err := utils.AESEncrypt(pubSymKey, utils.KEY)
	if err != nil {
		return []dto.UploadFileResponse{}, err
	}

	pubAccess := entities.PublicAccess{
		PublicSymmetricKey: encryptedPubSymKey,
		RequesterId:        uuid.MustParse(req.RequesterId),
		OwnerId:            uuid.MustParse(req.OwnerId),
	}

	_, err = s.pbr.Create(ctx, nil, pubAccess)
	if err != nil {
		return []dto.UploadFileResponse{}, err
	}

	result, err := s.fr.GetAllFileByUserId(ctx, nil, owner.ID.String())
	if err != nil {
		return []dto.UploadFileResponse{}, err
	}

	var files []dto.UploadFileResponse
	for _, file := range result {
		files = append(files, dto.UploadFileResponse{
			Path:           file.Path,
			Filename:       file.FileName,
			FileType:       file.FileType,
			Encryption:     file.Encryption,
			EncryptionMode: file.EncryptionMode,
		})
	}

	return files, nil
}
