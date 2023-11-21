package repository

import (
	"context"

	"github.com/Caknoooo/golang-clean_template/constants"
	"github.com/Caknoooo/golang-clean_template/entities"
	"gorm.io/gorm"
)

type (
	PrivateAccessRepository interface {
		Create(ctx context.Context, tx *gorm.DB, data entities.PrivateAccess) (entities.PrivateAccess, error)
		GetPrivateAccessRequestByUserId(ctx context.Context, tx *gorm.DB, userId string) ([]entities.PrivateAccess, error)
		GetPrivateAccessOwnerByUserId(ctx context.Context, tx *gorm.DB, userId string) ([]entities.PrivateAccess, error)
		GetPrivateAccessById(ctx context.Context, tx *gorm.DB, id string, userId string) (entities.PrivateAccess, error)
		Update(ctx context.Context, tx *gorm.DB, data entities.PrivateAccess) (entities.PrivateAccess, error)
	}

	privateAccessRepository struct {
		db *gorm.DB
	}
)

func NewPrivateAccessRepository(db *gorm.DB) PrivateAccessRepository {
	return &privateAccessRepository{
		db: db,
	}
}

func (r *privateAccessRepository) Create(ctx context.Context, tx *gorm.DB, data entities.PrivateAccess) (entities.PrivateAccess, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.Create(&data).Error; err != nil {
		return entities.PrivateAccess{}, err
	}

	return data, nil
}

func (r *privateAccessRepository) GetPrivateAccessById(ctx context.Context, tx *gorm.DB, id string, userId string) (entities.PrivateAccess, error) {
	if tx == nil {
		tx = r.db
	}

	var result entities.PrivateAccess
	if err := tx.Where("id = ? AND user_owner_id = ?", id, userId).Take(&result).Error; err != nil {
		return entities.PrivateAccess{}, err
	}

	return result, nil
}

func (r *privateAccessRepository) GetPrivateAccessRequestByUserId(ctx context.Context, tx *gorm.DB, userId string) ([]entities.PrivateAccess, error) {
	if tx == nil {
		tx = r.db
	}

	var result []entities.PrivateAccess
	if err := tx.Where("user_req_id = ?", userId).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *privateAccessRepository) GetPrivateAccessOwnerByUserId(ctx context.Context, tx *gorm.DB, userId string) ([]entities.PrivateAccess, error) {
	if tx == nil {
		tx = r.db
	}

	var result []entities.PrivateAccess
	if err := tx.Where("user_owner_id = ? AND status = ?", userId, constants.ENUM_STATUS_PENDING).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *privateAccessRepository) Update(ctx context.Context, tx *gorm.DB, data entities.PrivateAccess) (entities.PrivateAccess, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.Updates(&data).Error; err != nil {
		return entities.PrivateAccess{}, err
	}

	return data, nil
}