package category

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/validator"
)

type Service interface {
	Create(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*CategoryResponse, int, error)
	GetBySlug(ctx context.Context, slug string) (*CategoryResponse, int, error)
	GetAll(ctx context.Context) ([]CategoryResponse, int, error)
	GetActive(ctx context.Context) ([]CategoryResponse, int, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, int, error)
	Delete(ctx context.Context, id uuid.UUID) (int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		logrus.WithError(err).Error("Failed to create category")
		return nil, http.StatusInternalServerError, errors.New("failed to create category")
	}

	resp := s.mapToResponse(category)
	return &resp, http.StatusCreated, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*CategoryResponse, int, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get category by ID")
		return nil, http.StatusInternalServerError, errors.New("failed to get category")
	}

	resp := s.mapToResponse(category)
	return &resp, http.StatusOK, nil
}

func (s *service) GetBySlug(ctx context.Context, slug string) (*CategoryResponse, int, error) {
	category, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get category by slug")
		return nil, http.StatusInternalServerError, errors.New("failed to get category")
	}

	resp := s.mapToResponse(category)
	return &resp, http.StatusOK, nil
}

func (s *service) GetAll(ctx context.Context) ([]CategoryResponse, int, error) {
	categories, err := s.repo.FindAll(ctx)
	if err != nil {
		logrus.WithError(err).Error("Failed to get all categories")
		return nil, http.StatusInternalServerError, errors.New("failed to get categories")
	}

	responses := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = s.mapToResponse(&cat)
	}

	return responses, http.StatusOK, nil
}

func (s *service) GetActive(ctx context.Context) ([]CategoryResponse, int, error) {
	// Categories without IsActive field are 'active' if they are not soft-deleted.
	// GORM's FindAll automatically excludes soft-deleted records.
	return s.GetAll(ctx)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get category")
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	category.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, category); err != nil {
		logrus.WithError(err).Error("Failed to update category")
		return nil, http.StatusInternalServerError, errors.New("failed to update category")
	}

	resp := s.mapToResponse(category)
	return &resp, http.StatusOK, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) (int, error) {
	if err := s.repo.Delete(ctx, id); err != nil {
		logrus.WithError(err).Error("Failed to delete category")
		return http.StatusInternalServerError, errors.New("failed to delete category")
	}
	return http.StatusOK, nil
}

func (s *service) mapToResponse(c *models.Category) CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
