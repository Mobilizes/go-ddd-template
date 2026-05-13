package persistence

import (
	"fmt"
	"mob/ddd-template/internal/domain/entity"
	"mob/ddd-template/internal/domain/repository"
	vo "mob/ddd-template/internal/domain/valueobject"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FilePersistence struct {
	db *gorm.DB
}

func NewFilePersistence(db *gorm.DB) repository.FileRepository {
	return &FilePersistence{db: db}
}

func (p *FilePersistence) Create(file *entity.File) error {
	return p.db.Create(file).Error
}

func (p *FilePersistence) GetById(id string) (*entity.File, error) {
	var file entity.File
	if err := p.db.First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &file, nil
}

func (p *FilePersistence) GetAllByUser(userId string, opts *vo.PaginateOptions) (*vo.PaginatedResult[*entity.File], error) {
	var files []*entity.File
	var total int64

	query := p.db.Where("user_id = ?", userId)

	if opts.Filter != "" && opts.FilterBy != "" {
		query = query.Where(
			clause.Like{
				Column: clause.Column{Name: opts.FilterBy},
				Value:  "%" + opts.Filter + "%",
			},
		)
	}

	if err := query.Session(&gorm.Session{}).Model(&entity.File{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := opts.Page * opts.Limit
	order := fmt.Sprintf("%s %s", opts.SortBy, opts.Sort)

	query = query.Limit(opts.Limit).Offset(offset).Order(order)
	if err := query.Find(&files).Error; err != nil {
		return nil, err
	}

	return &vo.PaginatedResult[*entity.File]{
		Data:      files,
		Page:      opts.Page,
		Limit:     opts.Limit,
		TotalData: total,
		TotalPage: int((total + int64(opts.Limit) - 1) / int64(opts.Limit)),
	}, nil
}

func (p *FilePersistence) Update(file *entity.File) error {
	return p.db.Model(&entity.File{}).
		Where("id = ?", file.ID).
		Updates(map[string]any{
			"user_id":     file.UserID,
			"name":        file.Name,
			"storage_key": file.StorageKey,
			"mime_type":   file.MimeType,
			"size_bytes":  file.SizeBytes,
			"updated_at":  time.Now(),
			"deleted_at":  file.DeletedAt,
		}).Error
}

func (p *FilePersistence) Delete(id string) error {
	return p.db.Where("id = ?", id).Delete(&entity.File{}).Error
}
