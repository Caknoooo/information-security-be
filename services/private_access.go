package services

import (
	"bytes"
	"context"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/Caknoooo/golang-clean_template/constants"
	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/entities"
	"github.com/Caknoooo/golang-clean_template/repository"
	"github.com/Caknoooo/golang-clean_template/utils"
)

type (
	PrivateAccessService interface {
		Create(ctx context.Context, req dto.PrivateAccessRequest) (dto.PrivateAccessResponse, error)
		GetAllPrivateAccessRequestByUserId(ctx context.Context, userId string) ([]dto.GetPrivateAccessResponse, error)
		GetAllPrivateAccessOwnerByUserId(ctx context.Context, userId string) ([]dto.GetPrivateAccessResponse, error)
		UpdatePrivateAccess(ctx context.Context, req dto.UpdatePrivateAccessRequest, userId string) (dto.UpdatePrivateAccessResponse, error)
		SendEncryptionKey(ctx context.Context, req dto.SendEncryptionKeyRequest, userId string) (dto.SendEncryptionKeyResponse, error)
	}

	privateAccessService struct {
		userRepo          repository.UserRepository
		privateAccessRepo repository.PrivateAccessRepository
		fileRepo          repository.FileRepository
	}
)

func NewPrivateAccessService(userRepo repository.UserRepository, privateAccessRepo repository.PrivateAccessRepository, fileRepo repository.FileRepository) PrivateAccessService {
	return &privateAccessService{
		userRepo:          userRepo,
		privateAccessRepo: privateAccessRepo,
		fileRepo:          fileRepo,
	}
}

func (s *privateAccessService) Create(ctx context.Context, req dto.PrivateAccessRequest) (dto.PrivateAccessResponse, error) {
	userReq, err := s.userRepo.GetUserById(ctx, req.UserReqId)
	if err != nil {
		return dto.PrivateAccessResponse{}, dto.ErrGetUserById
	}

	userOwn, err := s.userRepo.GetUserById(ctx, req.UserOwnerId)
	if err != nil {
		return dto.PrivateAccessResponse{}, dto.ErrGetUserById
	}

	privAccess, err := s.privateAccessRepo.GetPrivateAccessRequestByUserId(ctx, nil, req.UserReqId)
	if len(privAccess) != 0 {
		return dto.PrivateAccessResponse{}, dto.ErrPrivateAccessExists
	} else if err != nil {
		return dto.PrivateAccessResponse{}, dto.ErrGetPrivateAccessById
	}

	privateAccess := entities.PrivateAccess{
		UserReqId:   userReq.ID.String(),
		UserOwnerId: userOwn.ID.String(),
	}

	data, err := s.privateAccessRepo.Create(ctx, nil, privateAccess)
	if err != nil {
		return dto.PrivateAccessResponse{}, dto.ErrCreatePrivateAccess
	}

	return dto.PrivateAccessResponse{
		ID:          data.ID.String(),
		UserReqId:   data.UserReqId,
		UserOwnerId: data.UserOwnerId,
		Status:      data.Status,
	}, nil
}

func (s *privateAccessService) GetAllPrivateAccessRequestByUserId(ctx context.Context, userId string) ([]dto.GetPrivateAccessResponse, error) {
	data, err := s.privateAccessRepo.GetPrivateAccessRequestByUserId(ctx, nil, userId)
	if err != nil {
		return nil, err
	}

	var result []dto.GetPrivateAccessResponse
	for _, v := range data {
		user, err := s.userRepo.GetUserById(ctx, v.UserOwnerId)
		if err != nil {
			return nil, err
		}

		decryptName, err := utils.AESDecrypt(user.Name, utils.KEY)
		if err != nil {
			return nil, err
		}

		result = append(result, dto.GetPrivateAccessResponse{
			ID:     v.ID.String(),
			UserID: user.ID.String(),
			Name:   decryptName,
			Email:  user.Email,
			Status: v.Status,
		})
	}

	return result, nil
}

func (s *privateAccessService) GetAllPrivateAccessOwnerByUserId(ctx context.Context, userId string) ([]dto.GetPrivateAccessResponse, error) {
	data, err := s.privateAccessRepo.GetPrivateAccessOwnerByUserId(ctx, nil, userId)
	if err != nil {
		return nil, err
	}

	var result []dto.GetPrivateAccessResponse
	for _, v := range data {
		user, err := s.userRepo.GetUserById(ctx, v.UserReqId)
		if err != nil {
			return nil, err
		}

		decryptName, err := utils.AESDecrypt(user.Name, utils.KEY)
		if err != nil {
			return nil, err
		}

		result = append(result, dto.GetPrivateAccessResponse{
			ID:      v.ID.String(),
			UserID:  user.ID.String(),
			Name:    decryptName,
			Email:   user.Email,
			Status:  v.Status,
			IsOwner: true,
		})
	}

	return result, nil
}

func (s *privateAccessService) UpdatePrivateAccess(ctx context.Context, req dto.UpdatePrivateAccessRequest, userId string) (dto.UpdatePrivateAccessResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return dto.UpdatePrivateAccessResponse{}, dto.ErrUserNotFound
	}

	data, err := s.privateAccessRepo.GetPrivateAccessById(ctx, nil, req.ID, user.ID.String())
	if err != nil {
		return dto.UpdatePrivateAccessResponse{}, err
	}

	if data.Status == constants.ENUM_STATUS_PENDING {
		if req.Status == constants.ENUM_STATUS_REJECTED {
			data.Status = constants.ENUM_STATUS_REJECTED
		} else if req.Status == constants.ENUM_STATUS_APPROVED {
			data.Status = constants.ENUM_STATUS_APPROVED
		} else {
			return dto.UpdatePrivateAccessResponse{}, dto.ErrStatusNotFound
		}
	}

	if data.Status == constants.ENUM_STATUS_APPROVED {
		userReq, err := s.userRepo.GetUserById(ctx, data.UserReqId)
		if err != nil {
			return dto.UpdatePrivateAccessResponse{}, dto.ErrGetUserById
		}

		pubKeyUserReqDecrypted, err := utils.AESDecrypt(userReq.PublicKey, utils.KEY)
		if err != nil {
			return dto.UpdatePrivateAccessResponse{}, err
		}

		symKeyWithDate := user.SymmetricKey + "_" + time.Now().String()

		rsaEncrypt, err := utils.EncryptRSA(symKeyWithDate, pubKeyUserReqDecrypted)
		if err != nil {
			return dto.UpdatePrivateAccessResponse{}, err
		}

		decryptName, err := utils.AESDecrypt(user.Name, utils.KEY)
		if err != nil {
			return dto.UpdatePrivateAccessResponse{}, err
		}

		data := map[string]string{
			"email":      userReq.Email,
			"rsa":        rsaEncrypt,
			"name_owner": decryptName,
		}

		draftEmail, err := privateAccessVerificationEmail(data)
		if err != nil {
			return dto.UpdatePrivateAccessResponse{}, err
		}

		if err := utils.SendMail(userReq.Email, draftEmail["subject"], draftEmail["body"]); err != nil {
			return dto.UpdatePrivateAccessResponse{}, err
		}
	}

	result, err := s.privateAccessRepo.Update(ctx, nil, data)
	if err != nil {
		return dto.UpdatePrivateAccessResponse{}, err
	}

	return dto.UpdatePrivateAccessResponse{
		ID:     result.ID.String(),
		Status: result.Status,
	}, nil
}

func privateAccessVerificationEmail(info map[string]string) (map[string]string, error) {
	readHtml, err := os.ReadFile("utils/email-template/private_access.html")
	if err != nil {
		return nil, err
	}

	data := struct {
		Email      string
		RSAEncrypt string
		NameOwner  string
	}{
		Email:      info["email"],
		RSAEncrypt: info["rsa"],
		NameOwner:  info["name_owner"],
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return nil, err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return nil, err
	}

	draftEmail := map[string]string{
		"subject": "Private Access Permission",
		"body":    strMail.String(),
	}

	return draftEmail, nil
}

func (s *privateAccessService) SendEncryptionKey(ctx context.Context, req dto.SendEncryptionKeyRequest, userId string) (dto.SendEncryptionKeyResponse, error) {
	userReq, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return dto.SendEncryptionKeyResponse{}, err
	}

	decryptPrivateKey, err := utils.AESDecrypt(userReq.PrivateKey, utils.KEY)
	if err != nil {
		return dto.SendEncryptionKeyResponse{}, err
	}

	decryptKey, err := utils.DecryptRSA(req.Key, decryptPrivateKey)
	if err != nil {
		return dto.SendEncryptionKeyResponse{}, err
	}

	decryptKey = strings.Split(decryptKey, "_")[0]

	userSym, err := s.userRepo.GetUserBySymmetricKey(ctx, nil, decryptKey)
	if err != nil {
		return dto.SendEncryptionKeyResponse{}, err
	}

	fileUser, err := s.fileRepo.GetLastSubmittedFilesByUserId(ctx, nil, userSym.ID.String(), []string{constants.ENUM_FILE_TYPE_IMAGE, constants.ENUM_FILE_TYPE_VIDEO, constants.ENUM_FILE_TYPE_FILE})
	if err != nil {
		return dto.SendEncryptionKeyResponse{}, err
	}

	var result []dto.UploadFileResponse
	for _, v := range fileUser {
		result = append(result, dto.UploadFileResponse{
			Path:           v.Path,
			Filename:       v.FileName,
			FileType:       v.FileType,
			Encryption:     v.Encryption,
			EncryptionMode: v.EncryptionMode,
		})
	}

	decryptName, err := utils.AESDecrypt(userSym.Name, utils.KEY)
	if err != nil {
		return dto.SendEncryptionKeyResponse{}, err
	}

	decrypTelp, err := utils.AESDecrypt(userSym.TelpNumber, utils.KEY)
	if err != nil {
		return dto.SendEncryptionKeyResponse{}, err
	}

	user := map[string]any{
		"email": userSym.Email,
		"name":  decryptName,
		"telp":  decrypTelp,
		"file":  result,
	}

	return dto.SendEncryptionKeyResponse{
		OwnerId: userSym.ID.String(),
		User:    user,
	}, nil
}
