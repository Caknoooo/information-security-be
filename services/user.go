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
	"github.com/Caknoooo/golang-clean_template/helpers"
	"github.com/Caknoooo/golang-clean_template/repository"
	"github.com/Caknoooo/golang-clean_template/utils"
)

type UserService interface {
	RegisterUser(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error)
	GetAllUser(ctx context.Context, adminId string) ([]dto.UserResponse, error)
	GetUserById(ctx context.Context, userId string) (dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error)
	GetUserByAdmin(ctx context.Context, adminId string, userId string) (dto.UserResponse, error)
	UpdateStatusIsVerified(ctx context.Context, req dto.UpdateStatusIsVerifiedRequest, adminId string) (dto.UserResponse, error)
	SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error
	VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error)
	CheckUser(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error)
	DeleteUser(ctx context.Context, userId string) error
	Verify(ctx context.Context, email string, password string) (bool, error)
}

const (
	LOCAL_URL          = "http://localhost:3000"
	VERIFY_EMAIL_ROUTE = "register/verify_email"
)

type userService struct {
	userRepo repository.UserRepository
	fileRepo repository.FileRepository
}

func NewUserService(ur repository.UserRepository, fr repository.FileRepository) UserService {
	return &userService{
		userRepo: ur,
		fileRepo: fr,
	}
}

func (s *userService) RegisterUser(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error) {
	email, _ := s.userRepo.CheckEmail(ctx, req.Email)
	if email {
		return dto.UserResponse{}, dto.ErrEmailAlreadyExists
	}

	encryptedName, _, err := utils.AESEncrypt(req.Name, utils.KEY)
	if err != nil {
		return dto.UserResponse{}, err
	}

	encryptedTelp, _, err := utils.AESEncrypt(req.TelpNumber, utils.KEY)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user := entities.User{
		Name:       encryptedName,
		TelpNumber: encryptedTelp,
		Role:       constants.ENUM_ROLE_USER,
		Email:      req.Email,
		Password:   req.Password,
		IsVerified: false,
	}

	userReg, err := s.userRepo.RegisterUser(ctx, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	draftEmail, err := makeVerificationEmail(req.Email)
	if err != nil {
		return dto.UserResponse{}, err
	}

	err = utils.SendMail(req.Email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:           userReg.ID.String(),
		Name:         userReg.Name,
		TelpNumber:   userReg.TelpNumber,
		Role:         userReg.Role,
		Email:        userReg.Email,
		IsVerified:   userReg.IsVerified,
		PublicKey:    userReg.PublicKey,
		PrivateKey:   userReg.PrivateKey,
		SymmetricKey: userReg.SymmetricKey,
	}, nil
}

func makeVerificationEmail(receiverEmail string) (map[string]string, error) {
	expired := time.Now().Add(time.Hour * 24).Format("2006-01-02 15:04:05")
	plainText := receiverEmail + "_" + expired
	token, datas, err := utils.AESEncrypt(plainText, utils.KEY)
	if err != nil {
		return nil, err
	}

	verifyLink := LOCAL_URL + "/" + VERIFY_EMAIL_ROUTE + "?token=" + token

	readHtml, err := os.ReadFile("utils/email-template/base_mail.html")
	if err != nil {
		return nil, err
	}

	data := struct {
		Email          string
		Verify         string
		AES_KEY        string
		AES_PLAIN_TEXT string
		AES_BLOCK      string
		AES_GCM        string
		AES_NONCE      string
		AES_RESULT     string
	}{
		Email:          receiverEmail,
		Verify:         verifyLink,
		AES_KEY:        datas["key"].(string),
		AES_PLAIN_TEXT: datas["plaintext"].(string),
		AES_BLOCK:      datas["block"].(string),
		AES_GCM:        datas["aes-gcm"].(string),
		AES_NONCE:      datas["nonce"].(string),
		AES_RESULT:     token,
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
		"subject": "Information Security F - Verification Email",
		"body":    strMail.String(),
	}

	return draftEmail, nil
}

func (s *userService) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	draftEmail, err := makeVerificationEmail(req.Email)
	if err != nil {
		return err
	}

	err = utils.SendMail(req.Email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
	decryptedToken, err := utils.AESDecrypt(req.Token, utils.KEY)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	if !strings.Contains(decryptedToken, "_") {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	decryptedTokenSplit := strings.Split(decryptedToken, "_")
	email := decryptedTokenSplit[0]
	expired := decryptedTokenSplit[1]

	now := time.Now()
	expiredTime, err := time.Parse("2006-01-02 15:04:05", expired)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
	}

	if expiredTime.Sub(now) < 0 {
		return dto.VerifyEmailResponse{
			Email:      email,
			IsVerified: false,
		}, dto.ErrTokenExpired
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUserNotFound
	}

	symKey, err := utils.GenerateKey(32)
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrCreateUser
	}

	privKey, pubKey, err := utils.GenerateRSAKey()
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrCreateUser
	}

	keys := [3]string{pubKey, privKey, symKey}
	var encryptedKeys []string
	for i := 0; i < 3; i++ {
		encryptedKey, _, err := utils.AESEncrypt(keys[i], utils.KEY)
		if err != nil {
			return dto.VerifyEmailResponse{}, err
		}
		encryptedKeys = append(encryptedKeys, encryptedKey)
	}

	updatedUser, err := s.userRepo.UpdateUser(ctx, entities.User{
		ID:           user.ID,
		IsVerified:   true,
		PublicKey:    encryptedKeys[0],
		PrivateKey:   encryptedKeys[1],
		SymmetricKey: encryptedKeys[2],
	})
	if err != nil {
		return dto.VerifyEmailResponse{}, dto.ErrUpdateUser
	}

	return dto.VerifyEmailResponse{
		Email:      email,
		IsVerified: updatedUser.IsVerified,
	}, nil
}

func (s *userService) GetAllUser(ctx context.Context, adminId string) ([]dto.UserResponse, error) {
	admin, err := s.userRepo.GetUserById(ctx, adminId)
	if err != nil {
		return []dto.UserResponse{}, dto.ErrUserNotFound
	}

	if admin.Role != constants.ENUM_ROLE_ADMIN {
		return []dto.UserResponse{}, dto.ErrUserNotAdmin
	}

	users, err := s.userRepo.GetAllUser(ctx)
	if err != nil {
		return []dto.UserResponse{}, dto.ErrGetAllUser
	}

	var userResponse []dto.UserResponse
	for _, user := range users {
		datas := []string{user.Name, user.TelpNumber, user.PublicKey}
		var encDatas []string
		for i := 0; i < len(datas); i++ {
			encData, err := utils.AESDecrypt(datas[i], utils.KEY)
			if err != nil {
				return []dto.UserResponse{}, err
			}
			encDatas = append(encDatas, encData)
		}

		userResponse = append(userResponse, dto.UserResponse{
			ID:         user.ID.String(),
			Name:       encDatas[0],
			TelpNumber: encDatas[1],
			Role:       user.Role,
			Email:      user.Email,
			IsVerified: user.IsVerified,
			PublicKey:  encDatas[2],
			CreatedAt:  string(user.CreatedAt.Format("2006-01-02 15:04:05")),
		})
	}

	return userResponse, nil
}

func (s *userService) GetUserByAdmin(ctx context.Context, adminId string, userId string) (dto.UserResponse, error) {
	admin, err := s.userRepo.GetUserById(ctx, adminId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}

	if admin.Role != constants.ENUM_ROLE_ADMIN {
		return dto.UserResponse{}, dto.ErrUserNotAdmin
	}

	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}

	fileUser, err := s.fileRepo.GetAllFileByUserId(ctx, nil, user.ID.String())
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetAllFileByUserId
	}

	var files []dto.UploadFileResponse
	for _, file := range fileUser {
		data := dto.UploadFileResponse{
			Path:           file.Path,
			Filename:       file.FileName,
			FileType:       file.FileType,
			Encryption:     file.Encryption,
			EncryptionMode: file.EncryptionMode,
		}

		files = append(files, data)
	}

	datas := []string{user.Name, user.TelpNumber, user.PublicKey, user.PrivateKey, user.SymmetricKey}
	var encDatas []string
	for i := 0; i < len(datas); i++ {
		encData, err := utils.AESDecrypt(datas[i], utils.KEY)
		if err != nil {
			return dto.UserResponse{}, err
		}
		encDatas = append(encDatas, encData)
	}

	return dto.UserResponse{
		ID:           user.ID.String(),
		Name:         encDatas[0],
		TelpNumber:   encDatas[1],
		Role:         user.Role,
		Email:        user.Email,
		IsVerified:   user.IsVerified,
		PublicKey:    encDatas[2],
		PrivateKey:   encDatas[3],
		SymmetricKey: encDatas[4],
		Files:        files,
		Work:         user.Work,
		CreatedAt:    string(user.CreatedAt.Format("2006-01-02 15:04:05")),
	}, nil
}

func (s *userService) UpdateStatusIsVerified(ctx context.Context, req dto.UpdateStatusIsVerifiedRequest, adminId string) (dto.UserResponse, error) {
	admin, err := s.userRepo.GetUserById(ctx, adminId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}

	if admin.Role != constants.ENUM_ROLE_ADMIN {
		return dto.UserResponse{}, dto.ErrUserNotAdmin
	}

	user, err := s.userRepo.GetUserById(ctx, req.UserId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}

	data := entities.User{
		ID:         user.ID,
		IsVerified: req.IsVerified,
	}

	userUpdate, err := s.userRepo.UpdateUser(ctx, data)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUpdateUser
	}

	return dto.UserResponse{
		ID:         user.ID.String(),
		Name:       user.Name,
		TelpNumber: user.TelpNumber,
		Role:       user.Role,
		Email:      user.Email,
		IsVerified: userUpdate.IsVerified,
	}, nil
}

func (s *userService) GetUserById(ctx context.Context, userId string) (dto.UserResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserById
	}

	datas := []string{user.Name, user.TelpNumber, user.PublicKey, user.PrivateKey, user.SymmetricKey}
	var encDatas []string
	for i := 0; i < len(datas); i++ {
		encData, err := utils.AESDecrypt(datas[i], utils.KEY)
		if err != nil {
			return dto.UserResponse{}, err
		}
		encDatas = append(encDatas, encData)
	}

	return dto.UserResponse{
		ID:           user.ID.String(),
		Name:         encDatas[0],
		TelpNumber:   encDatas[1],
		Role:         user.Role,
		Email:        user.Email,
		PublicKey:    encDatas[2],
		PrivateKey:   encDatas[3],
		SymmetricKey: encDatas[4],
		IsVerified:   user.IsVerified,
	}, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserByEmail
	}

	datas := []string{user.Name, user.TelpNumber, user.PublicKey}
	var encDatas []string
	for i := 0; i < len(datas); i++ {
		encData, err := utils.AESDecrypt(datas[i], utils.KEY)
		if err != nil {
			return dto.UserResponse{}, err
		}
		encDatas = append(encDatas, encData)
	}

	return dto.UserResponse{
		ID:         user.ID.String(),
		Name:       encDatas[0],
		TelpNumber: encDatas[1],
		Role:       user.Role,
		Email:      user.Email,
		PublicKey:  encDatas[2],
		IsVerified: user.IsVerified,
	}, nil
}

func (s *userService) CheckUser(ctx context.Context, email string) (bool, error) {
	res, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if res.Email == "" {
		return false, err
	}
	return true, nil
}

func (s *userService) UpdateUser(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUserNotFound
	}

	data := entities.User{
		ID:         user.ID,
		Name:       req.Name,
		TelpNumber: req.TelpNumber,
		Work:       req.Work,
		Role:       user.Role,
		Email:      req.Email,
		Password:   req.Password,
	}

	userUpdate, err := s.userRepo.UpdateUser(ctx, data)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUpdateUser
	}

	return dto.UserUpdateResponse{
		ID:         userUpdate.ID.String(),
		Name:       userUpdate.Name,
		TelpNumber: userUpdate.TelpNumber,
		Role:       userUpdate.Role,
		Email:      userUpdate.Email,
		IsVerified: userUpdate.IsVerified,
		Work:       userUpdate.Work,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, userId string) error {
	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return dto.ErrUserNotFound
	}

	err = s.userRepo.DeleteUser(ctx, user.ID.String())
	if err != nil {
		return dto.ErrDeleteUser
	}

	return nil
}

func (s *userService) Verify(ctx context.Context, email string, password string) (bool, error) {
	res, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return false, dto.ErrUserNotFound
	}

	if !res.IsVerified {
		return false, dto.ErrAccountNotVerified
	}

	checkPassword, err := helpers.CheckPassword(res.Password, []byte(password))
	if err != nil {
		return false, dto.ErrPasswordNotMatch
	}

	if res.Email == email && checkPassword {
		return true, nil
	}

	return false, dto.ErrEmailOrPassword
}
