package blog

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/database"
)

var (
	ErrBlogNotFound = errors.New("blog not found")
)

type Repository interface {
	Create(ctx context.Context, blog *models.Blog) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Blog, error)
	FindBySlug(ctx context.Context, slug string) (*models.Blog, error)
	FindAll(ctx context.Context, limit, offset int, status string) ([]models.Blog, error)
	Update(ctx context.Context, blog *models.Blog) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db database.Database
}

func NewRepository(db database.Database) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, blog *models.Blog) error {
	if err := r.db.GetDB().WithContext(ctx).Create(blog).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Blog, error) {
	var b models.Blog
	if err := r.db.GetDB().WithContext(ctx).
		Preload("Author").
		Preload("Category").
		Where("id = ?", id).First(&b).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBlogNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *repository) FindBySlug(ctx context.Context, slug string) (*models.Blog, error) {
	var b models.Blog
	if err := r.db.GetDB().WithContext(ctx).
		Preload("Author").
		Preload("Category").
		Where("slug = ?", slug).First(&b).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBlogNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *repository) FindAll(ctx context.Context, limit, offset int, status string) ([]models.Blog, error) {
	var blogs []models.Blog
	query := r.db.GetDB().WithContext(ctx).
		Preload("Author").
		Preload("Category")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Order("created_at desc").Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *repository) Update(ctx context.Context, blog *models.Blog) error {
	if err := r.db.GetDB().WithContext(ctx).Save(blog).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.GetDB().WithContext(ctx).Delete(&models.Blog{}, id).Error; err != nil {
		return err
	}
	return nil
}