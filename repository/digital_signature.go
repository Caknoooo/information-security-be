package repository

import (
	"context"
	"math"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/entities"
	"gorm.io/gorm"
)

type (
	DigitalSignatureRepository interface {
		Create(ctx context.Context, tx *gorm.DB, digitalSignature entities.DigitalSignature) (entities.DigitalSignature, error)
		GetAllDigitalSignatureReceiver(ctx context.Context, tx *gorm.DB, receiverID string, req dto.PaginationRequest) (dto.GetAllNotificationsRepository, error)
		GetDigitalSignatureById(ctx context.Context, tx *gorm.DB, id string) (entities.DigitalSignature, error)
		Update(ctx context.Context, tx *gorm.DB, digitalSignature entities.DigitalSignature) (entities.DigitalSignature, error)
		Delete(ctx context.Context, tx *gorm.DB, digitalSignature entities.DigitalSignature) error
	}

	digitalSignatureRepository struct {
		db *gorm.DB
	}
)

func NewDigitalSignatureRepository(db *gorm.DB) DigitalSignatureRepository {
	return &digitalSignatureRepository{
		db: db,
	}
}

func (r *digitalSignatureRepository) Create(ctx context.Context, tx *gorm.DB, digitalSignature entities.DigitalSignature) (entities.DigitalSignature, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(digitalSignature).Error; err != nil {
		return entities.DigitalSignature{}, err
	}

	return digitalSignature, nil
}
func (r *digitalSignatureRepository) GetAllDigitalSignatureReceiver(ctx context.Context, tx *gorm.DB, receiverID string, req dto.PaginationRequest) (dto.GetAllNotificationsRepository, error) {
	if tx == nil {
		tx = r.db
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}

	var count int64
	var err error
	var digitalSignatures []entities.DigitalSignature

	stmt := tx.Table("digital_signatures").
        Select("digital_signatures.*, users.email").
        Joins("JOIN users ON digital_signatures.sender_id = users.id").
        Or("digital_signatures.receiver_id = ?", receiverID)

	if req.Search != "" {
		err = stmt.Where("users.email ILIKE ?", "%"+req.Search+"%").Count(&count).Error
		if err != nil {
			return dto.GetAllNotificationsRepository{}, err
		}
	} else {
		err = stmt.Count(&count).Error
		if err != nil {
			return dto.GetAllNotificationsRepository{}, err
		}
	}

	maxPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	if req.PerPage <= 0 {
		stmt.Find(&digitalSignatures)
		return dto.GetAllNotificationsRepository{
			DigitalSignature: digitalSignatures,
			PaginationResponse: dto.PaginationResponse{
				Page:    req.Page,
				PerPage: req.PerPage,
				MaxPage: maxPage,
				Count:   count,
			},
		}, nil
	}

	offset := (req.Page - 1) * req.PerPage
	err = stmt.Offset(offset).Limit(req.PerPage).Find(&digitalSignatures).Error
	if err != nil {
		return dto.GetAllNotificationsRepository{}, err
	}

	return dto.GetAllNotificationsRepository{
		DigitalSignature: digitalSignatures,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}


func (r *digitalSignatureRepository) GetDigitalSignatureById(ctx context.Context, tx *gorm.DB, id string) (entities.DigitalSignature, error) {
	if tx == nil {
		tx = r.db
	}

	var digitalSignature entities.DigitalSignature

	if err := tx.WithContext(ctx).Where("id = ?", id).First(&digitalSignature).Error; err != nil {
		return entities.DigitalSignature{}, err
	}

	return digitalSignature, nil
}

func (r *digitalSignatureRepository) Update(ctx context.Context, tx *gorm.DB, digitalSignature entities.DigitalSignature) (entities.DigitalSignature, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(digitalSignature).Error; err != nil {
		return entities.DigitalSignature{}, err
	}

	return digitalSignature, nil
}

func (r *digitalSignatureRepository) Delete(ctx context.Context, tx *gorm.DB, digitalSignature entities.DigitalSignature) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(digitalSignature).Error; err != nil {
		return err
	}

	return nil
}
