package category

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/database"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type Repository interface {
	Create(ctx context.Context, category *models.Category) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	FindBySlug(ctx context.Context, slug string) (*models.Category, error)
	FindAll(ctx context.Context) ([]models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db database.Database
}

func NewRepository(db database.Database) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, category *models.Category) error {
	if err := r.db.GetDB().WithContext(ctx).Create(category).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var cat models.Category
	if err := r.db.GetDB().WithContext(ctx).Where("id = ?", id).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &cat, nil
}

func (r *repository) FindBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var cat models.Category
	if err := r.db.GetDB().WithContext(ctx).Where("slug = ?", slug).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &cat, nil
}

func (r *repository) FindAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.GetDB().WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) Update(ctx context.Context, category *models.Category) error {
	if err := r.db.GetDB().WithContext(ctx).Save(category).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.GetDB().WithContext(ctx).Delete(&models.Category{}, id).Error; err != nil {
		return err
	}
	return nil
}