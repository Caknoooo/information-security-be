package repository

import (
	"context"

	"github.com/Caknoooo/golang-clean_template/entities"
	"gorm.io/gorm"
)

type (
	DigitalSignatureRepository interface {
		Create(ctx context.Context, tx *gorm.DB, digitalSignature entities.DigitalSignature) (entities.DigitalSignature, error)
		GetAllDigitalSignatureReceiver(ctx context.Context, tx *gorm.DB, receiverID string) ([]entities.DigitalSignature, error)
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

func (r *digitalSignatureRepository) GetAllDigitalSignatureReceiver(ctx context.Context, tx *gorm.DB, receiverID string) ([]entities.DigitalSignature, error) {
	if tx == nil {
		tx = r.db
	}

	var digitalSignature []entities.DigitalSignature

	if err := tx.WithContext(ctx).Where("receiver_id = ?", receiverID).Find(&digitalSignature).Error; err != nil {
		return []entities.DigitalSignature{}, err
	}

	return digitalSignature, nil
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
