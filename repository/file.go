package repository

import (
	"context"
	"time"

	"github.com/Caknoooo/golang-clean_template/entities"
	"gorm.io/gorm"
)

type (
	FileRepository interface {
		Create(ctx context.Context, tx *gorm.DB, file entities.File) (entities.File, error)
		GetAllFile(ctx context.Context, tx *gorm.DB) ([]entities.File, error)
		GetAllFileByUserId(ctx context.Context, tx *gorm.DB, userId string) ([]entities.File, error)
		GetLastSubmittedFilesByUserId(ctx context.Context, tx *gorm.DB, userId string, fileTypes []string) ([]entities.File, error)
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

func (r *fileRepository) GetLastSubmittedFilesByUserId(ctx context.Context, tx *gorm.DB, userId string, fileTypes []string) ([]entities.File, error) {
    if tx == nil {
        tx = r.db
    }

    var files []entities.File

    var subqueryResults []struct {
        UserID       string
        FileType     string
        MaxCreatedAt time.Time
    }

    if err := tx.Model(&entities.File{}).
        Select("user_id, file_type, MAX(created_at) as max_created_at").
        Where("user_id = ?", userId).
        Where("file_type IN (?)", fileTypes).
        Group("user_id, file_type").
        Find(&subqueryResults).
        Error; err != nil {
        return nil, err
    }

    if len(subqueryResults) > 0 {
        var maxCreatedAtMap = make(map[string]map[string]time.Time)

        for _, result := range subqueryResults {
            if _, ok := maxCreatedAtMap[result.UserID]; !ok {
                maxCreatedAtMap[result.UserID] = make(map[string]time.Time)
            }
            maxCreatedAtMap[result.UserID][result.FileType] = result.MaxCreatedAt
        }

        if err := tx.
            Where("user_id = ?", userId).
            Where("file_type IN (?)", fileTypes).
            Find(&files).
            Error; err != nil {
            return nil, err
        }

        var selectedFiles []entities.File
        for _, fileType := range fileTypes {
            for i := len(files) - 1; i >= 0; i-- {
                file := files[i]
                if file.FileType == fileType {
                    selectedFiles = append(selectedFiles, file)
                    break
                }
            }
        }

        return selectedFiles, nil
    }

    return nil, nil
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
