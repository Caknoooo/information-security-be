package repository

import (
	"context"

	"github.com/Caknoooo/golang-clean_template/entities"
	"gorm.io/gorm"
)

type (
	PublicAccessRepository interface {
		Create(ctx context.Context, tx *gorm.DB, publicAccess entities.PublicAccess) (entities.PublicAccess, error)
		GetAllPublicAccess(ctx context.Context, tx *gorm.DB) ([]entities.PublicAccess, error)
		GetAllPublicAccessByOwnerId(ctx context.Context, tx *gorm.DB, ownerId string) ([]entities.PublicAccess, error)
		GetAllPublicAccessByRequesterId(ctx context.Context, tx *gorm.DB, reqId string) ([]entities.PublicAccess, error)
		GetPublicAccessById(ctx context.Context, tx *gorm.DB, publicAccessId string) (entities.PublicAccess, error)
		GetAllPublicAccessByIDs(ctx context.Context, tx *gorm.DB, ownerId string, reqId string) ([]entities.PublicAccess, error)
	}

	publicAccessRepository struct {
		db *gorm.DB
	}
)

func NewPublicAccessRepository(db *gorm.DB) PublicAccessRepository {
	return &publicAccessRepository{
		db: db,
	}
}

func (r *publicAccessRepository) Create(ctx context.Context, tx *gorm.DB, publicAccess entities.PublicAccess) (entities.PublicAccess, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.Create(&publicAccess).Error; err != nil {
		return entities.PublicAccess{}, err
	}
	return publicAccess, nil
}

func (r *publicAccessRepository) GetAllPublicAccess(ctx context.Context, tx *gorm.DB) ([]entities.PublicAccess, error) {
	if tx == nil {
		tx = r.db
	}

	var publicAccess []entities.PublicAccess
	if err := tx.Find(&publicAccess).Error; err != nil {
		return nil, err
	}
	return publicAccess, nil
}

func (r *publicAccessRepository) GetAllPublicAccessByOwnerId(ctx context.Context, tx *gorm.DB, ownerId string) ([]entities.PublicAccess, error) {
	if tx == nil {
		tx = r.db
	}

	var publicAccess []entities.PublicAccess
	if err := tx.Where("owner_id = ?", ownerId).Find(&publicAccess).Error; err != nil {
		return nil, err
	}
	return publicAccess, nil
}

func (r *publicAccessRepository) GetAllPublicAccessByRequesterId(ctx context.Context, tx *gorm.DB, reqId string) ([]entities.PublicAccess, error) {
	if tx == nil {
		tx = r.db
	}

	var publicAccess []entities.PublicAccess
	if err := tx.Where("requester_id = ?", reqId).Find(&publicAccess).Error; err != nil {
		return nil, err
	}
	return publicAccess, nil
}

func (r *publicAccessRepository) GetPublicAccessById(ctx context.Context, tx *gorm.DB, publicAccessId string) (entities.PublicAccess, error) {
	if tx == nil {
		tx = r.db
	}

	var publicAccess entities.PublicAccess
	if err := tx.Where("id = ?", publicAccessId).Take(&publicAccess).Error; err != nil {
		return entities.PublicAccess{}, err
	}
	return publicAccess, nil
}

func (r *publicAccessRepository) GetAllPublicAccessByIDs(ctx context.Context, tx *gorm.DB, ownerId string, reqId string) ([]entities.PublicAccess, error) {
	if tx == nil {
		tx = r.db
	}

	var publicAccess []entities.PublicAccess
	if err := tx.Where("owner_id = ?", ownerId).Where("requester_id = ?", reqId).Find(&publicAccess).Error; err != nil {
		return nil, err
	}
	return publicAccess, nil
}
