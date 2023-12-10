package services

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/entities"
	"github.com/Caknoooo/golang-clean_template/repository"
	"github.com/Caknoooo/golang-clean_template/utils"
	"github.com/google/uuid"
)

type (
	DigitalSignatureService interface {
		CreateDigitalSignature(ctx context.Context, req dto.DigitalSignatureRequest) (dto.DigitalSignatureResponse, error)
		VerifyDigitalSignature(ctx context.Context, req dto.VerifyDigitalSignatureRequest) (dto.VerifyDigitalSignatureResponse, error)
		GetAllNotifications(ctx context.Context, userId string, req dto.PaginationRequest) (dto.GetAllNotificationsWithPaginationResponse, error)
	}

	digitalSignatureService struct {
		digitalSignatureRepo repository.DigitalSignatureRepository
		userRepo             repository.UserRepository
	}
)

func NewDigitalSignatureService(digitalSignatureRepo repository.DigitalSignatureRepository, userRepo repository.UserRepository) DigitalSignatureService {
	return &digitalSignatureService{
		digitalSignatureRepo: digitalSignatureRepo,
		userRepo:             userRepo,
	}
}

const (
	DataCommentKey      = "DataKeyF-02_"
	SignatureCommentKey = "SignatureKeyF-02_"
	PublicKeyCommentKey = "PublicKeyKeyF-02_"
	API_URL             = "www.isf.sre-its.com/static"
)

type (
	Signing struct {
		Name    string
		Email   string
		Release string
	}

	ReadContents struct {
		Data      []byte
		Signature []byte
		PublicKey []byte
	}
)

func (s *digitalSignatureService) CreateDigitalSignature(ctx context.Context, req dto.DigitalSignatureRequest) (dto.DigitalSignatureResponse, error) {
	from, err := s.userRepo.GetUserByEmail(ctx, req.From)
	if err != nil {
		return dto.DigitalSignatureResponse{}, dto.ErrUserNotFound
	}

	to, err := s.userRepo.GetUserByEmail(ctx, req.To)
	if err != nil {
		return dto.DigitalSignatureResponse{}, dto.ErrUserNotFound
	}

	ext := utils.GetExtension(req.BodyFiles.Filename)

	fileId := uuid.New()
	filename := fmt.Sprintf("%s/digital_signature/%s.%s", from.ID, fileId, ext)
	if err := utils.UploadFileSuccess(req.BodyFiles, filename); err != nil {
		return dto.DigitalSignatureResponse{}, err
	}

	fullPath := utils.PATH + "/" + filename
	// Read file to get the content
	file, err := utils.Read(fullPath)
	if err != nil {
		return dto.DigitalSignatureResponse{}, err
	}

	messageAdd, err := WriteContent(file, from)
	if err != nil {
		return dto.DigitalSignatureResponse{}, err
	}

	// Write to new file
	err = utils.Write(fullPath, messageAdd)
	if err != nil {
		return dto.DigitalSignatureResponse{}, err
	}

	// encryption to save the path
	encryption, _, err := utils.AESEncrypt(filename, utils.FILE_KEY_AES)
	if err != nil {
		return dto.DigitalSignatureResponse{}, err
	}

	// Save to database
	digitalSignature := entities.DigitalSignature{
		ID:         uuid.New(),
		SenderID:   from.ID.String(),
		ReceiverID: to.ID.String(),
		Subject:    req.Subject,
		Content:    req.BodyContent,
		Filepath:   encryption,
		IsSigned:   true,
	}

	digitalSignature, err = s.digitalSignatureRepo.Create(ctx, nil, digitalSignature)
	if err != nil {
		return dto.DigitalSignatureResponse{}, err
	}

	decryptName, err := utils.AESDecrypt(from.Name, utils.KEY)
	if err != nil {
		return dto.DigitalSignatureResponse{}, err
	}

	data := map[string]string{
		"email":        to.Email,
		"subject":      req.Subject,
		"body_content": req.BodyContent,
		"filepath":     API_URL + "/" + filename,
		"name_owner":   decryptName,
	}

	draftEmail, err := SendDigitalSignatureMail(data)
	if err != nil {
		return dto.DigitalSignatureResponse{}, err
	}

	if err := utils.SendMail(to.Email, draftEmail["subject"], draftEmail["body"], fullPath); err != nil {
		return dto.DigitalSignatureResponse{}, err
	}

	return dto.DigitalSignatureResponse{
		ID:         digitalSignature.ID.String(),
		SenderID:   digitalSignature.SenderID,
		ReceiverID: digitalSignature.ReceiverID,
		Subject:    digitalSignature.Subject,
		Content:    digitalSignature.Content,
		Filepath:   digitalSignature.Filepath,
		IsSigned:   digitalSignature.IsSigned,
	}, nil
}

func SendDigitalSignatureMail(info map[string]string) (map[string]string, error) {
	readHtml, err := utils.Read("utils/email-template/digital_signature.html")
	if err != nil {
		return nil, err
	}

	data := struct {
		Email       string
		Subject     string
		BodyContent string
		Filepath    string
		NameOwner   string
	}{
		Email:       info["email"],
		Subject:     info["subject"],
		BodyContent: info["body_content"],
		Filepath:    info["filepath"],
		NameOwner:   info["name_owner"],
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
		"subject": "Digital Signature",
		"body":    strMail.String(),
	}

	return draftEmail, nil
}

func WriteContent(content []byte, from entities.User) ([]byte, error) {
	var messagesAdded []byte
	hash := sha256.Sum256(content)

	decryptName, err := utils.AESDecrypt(from.Name, utils.KEY)
	if err != nil {
		return nil, err
	}

	signing := Signing{
		Name:    decryptName,
		Email:   from.Email,
		Release: time.Now().Format("2006-01-02 15:04:05"),
	}

	signingBytes, err := json.Marshal(signing)
	if err != nil {
		return nil, err
	}

	decryptPrivateKey, err := utils.AESDecrypt(from.PrivateKey, utils.KEY)
	if err != nil {
		return nil, err
	}

	decryptPublicKey, err := utils.AESDecrypt(from.PublicKey, utils.KEY)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, err := utils.ParsePrivateKeyFromPEM(decryptPrivateKey)
	if err != nil {
		return nil, err
	}

	data_signature, err := rsa.SignPKCS1v15(nil, rsaPrivateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}

	messagesAdded = append(messagesAdded, content...)

	messagesAdded = append(messagesAdded, []byte("\n%"+DataCommentKey)...)
	messagesAdded = append(messagesAdded, data_signature...)

	messagesAdded = append(messagesAdded, []byte("\n%"+SignatureCommentKey)...)
	messagesAdded = append(messagesAdded, signingBytes...)

	messagesAdded = append(messagesAdded, []byte("\n%"+PublicKeyCommentKey)...)
	messagesAdded = append(messagesAdded, decryptPublicKey...)

	return messagesAdded, nil
}

func (s *digitalSignatureService) VerifyDigitalSignature(ctx context.Context, req dto.VerifyDigitalSignatureRequest) (dto.VerifyDigitalSignatureResponse, error) {
	var data Signing

	ext := utils.GetExtension(req.Files.Filename)
	fileId := uuid.New()
	filename := fmt.Sprintf("%s/verify_digital_signature/%s.%s", req.UserId, fileId, ext)
	if err := utils.UploadFileSuccess(req.Files, filename); err != nil {
		return dto.VerifyDigitalSignatureResponse{}, err
	}

	fullPath := utils.PATH + "/" + filename

	// Read file to get the content
	file, err := utils.Read(fullPath)
	if err != nil {
		return dto.VerifyDigitalSignatureResponse{}, err
	}

	readContent, err := ReadContent(file)
	if err != nil {
		return dto.VerifyDigitalSignatureResponse{}, err
	}

	decryptPubKey := readContent.PublicKey

	// Get Public Key From Embedded
	pubKey, err := utils.ParsePublicKeyFromPEM(string(decryptPubKey))
	if err != nil {
		return dto.VerifyDigitalSignatureResponse{}, err
	}

	parts := bytes.Split(file, []byte("\n%"+DataCommentKey))
	if len(parts) < 2 {
		return dto.VerifyDigitalSignatureResponse{}, dto.ErrInvalidDigitalSignature
	}

	content := parts[0]
	signature := readContent.Data

	hash := sha256.Sum256(content)

	// Unmarshal Data
	err = json.Unmarshal(readContent.Signature, &data)
	if err != nil {
		return dto.VerifyDigitalSignatureResponse{}, err
	}

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return dto.VerifyDigitalSignatureResponse{}, dto.ErrPdfFileDifferentContent
	}

	return dto.VerifyDigitalSignatureResponse{
		IsVerified: true,
		SenderVerifyResponse: dto.SenderVerifyResponse{
			Name:  data.Name,
			Email: data.Email,
			Date:  data.Release,
		},
	}, nil
}

func ReadContent(content []byte) (ReadContents, error) {
	var buffer []byte

	tokens := []string{
		"%" + DataCommentKey,
		"%" + SignatureCommentKey,
		"%" + PublicKeyCommentKey,
	}

	results := [3]string{}
	ctr := 0
	idx := 0

	stringFileContent := string(content)
	for idx < len(stringFileContent) {
		if ctr <= 2 &&
			idx+len(tokens[ctr]) < len(stringFileContent) &&
			stringFileContent[idx:idx+len(tokens[ctr])] == tokens[ctr] {

			if ctr > 0 && len(buffer) > 0 {
				results[ctr-1] = string(buffer[:len(buffer)-1])
				buffer = []byte{}
			}

			idx += len(stringFileContent[idx : idx+len(tokens[ctr])])
			ctr++
		}

		if ctr > 0 {
			buffer = append(buffer, stringFileContent[idx])
		}

		if idx == len(stringFileContent)-1 && ctr != 0 {
			results[ctr-1] = string(buffer)
		}

		idx++
	}

	return ReadContents{
		Data:      []byte(results[0]),
		Signature: []byte(results[1]),
		PublicKey: []byte(results[2]),
	}, nil
}

func (s *digitalSignatureService) GetAllNotifications(ctx context.Context, userId string, req dto.PaginationRequest) (dto.GetAllNotificationsWithPaginationResponse, error) {
	digitalSignatures, err := s.digitalSignatureRepo.GetAllDigitalSignatureReceiver(ctx, nil, userId, req)
	if err != nil {
		return dto.GetAllNotificationsWithPaginationResponse{}, err
	}

	var notifications []dto.GetAllNotificationsResponse
	for _, digitalSignature := range digitalSignatures.DigitalSignature {
		user, err := s.userRepo.GetUserById(ctx, digitalSignature.SenderID)
		if err != nil {
			return dto.GetAllNotificationsWithPaginationResponse{}, err
		}

		decryptName, err := utils.AESDecrypt(user.Name, utils.KEY)
		if err != nil {
			return dto.GetAllNotificationsWithPaginationResponse{}, err
		}

		notification := dto.GetAllNotificationsResponse{
			Name:        decryptName,
			Email:       user.Email,
			Subject:     digitalSignature.Subject,
			BodyContent: digitalSignature.Content,
			Filepath:    digitalSignature.Filepath,
		}

		notifications = append(notifications, notification)
	}

	return dto.GetAllNotificationsWithPaginationResponse{
		Notifications:      notifications,
		PaginationResponse: digitalSignatures.PaginationResponse,
	}, nil
}
