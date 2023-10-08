package repository

import (
	"context"

	"github.com/Caknoooo/golang-clean_template/entities"
	"gorm.io/gorm"
)

type (
	FileRepository interface {
		Create(ctx context.Context, tx *gorm.DB, file entities.File) (entities.File, error)
		GetAllFile(ctx context.Context, tx *gorm.DB) ([]entities.File, error)
		GetAllFileByUserId(ctx context.Context, tx *gorm.DB, userId string) ([]entities.File, error)
		GetFileById(ctx context.Context, tx *gorm.DB, fileId string) (entities.File, error)
	}

	fileRepository struct {
		db *gorm.DB
	}
)

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{
		db: db,
	}
}

func (r *fileRepository) Create(ctx context.Context, tx *gorm.DB, file entities.File) (entities.File, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.Create(&file).Error; err != nil {
		return entities.File{}, err
	}
	return file, nil
}

func (r *fileRepository) GetAllFile(ctx context.Context, tx *gorm.DB) ([]entities.File, error) {
	if tx == nil {
		tx = r.db
	}

	var file []entities.File
	if err := tx.Find(&file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (r *fileRepository) GetAllFileByUserId(ctx context.Context, tx *gorm.DB, userId string) ([]entities.File, error) {
	if tx == nil {
		tx = r.db
	}

	var file []entities.File
	if err := tx.Where("user_id = ?", userId).Find(&file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (r *fileRepository) GetFileById(ctx context.Context, tx *gorm.DB, fileId string) (entities.File, error) {
	if tx == nil {
		tx = r.db
	}

	var file entities.File
	if err := tx.Where("id = ?", fileId).Take(&file).Error; err != nil {
		return entities.File{}, err
	}
	return file, nil
}