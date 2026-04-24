package blog

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
	Create(ctx context.Context, authorID uuid.UUID, req CreateBlogRequest) (*BlogResponse, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*BlogResponse, int, error)
	GetBySlug(ctx context.Context, slug string) (*BlogResponse, int, error)
	GetAll(ctx context.Context, limit, offset int, status string) ([]BlogResponse, int, error)
	GetPublishedBlogs(ctx context.Context, limit, offset int) ([]BlogResponse, int, error)
	GetPublishedBySlug(ctx context.Context, slug string) (*BlogResponse, int, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateBlogRequest) (*BlogResponse, int, error)
	Delete(ctx context.Context, id uuid.UUID) (int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, authorID uuid.UUID, req CreateBlogRequest) (*BlogResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	blog := &models.Blog{
		Title:         req.Title,
		Slug:          req.Slug,
		Content:       req.Content,
		Excerpt:       req.Excerpt,
		CoverImageURL: req.CoverImageURL,
		Status:        req.Status,
		CategoryID:    req.CategoryID,
		AuthorID:      authorID,
	}

	if req.Status == "published" {
		now := time.Now()
		blog.PublishedAt = &now
	}

	if err := s.repo.Create(ctx, blog); err != nil {
		logrus.WithError(err).Error("Failed to create blog")
		return nil, http.StatusInternalServerError, errors.New("failed to create blog")
	}

	// Refetch to get populated Author and Category 
	createdBlog, err := s.repo.FindByID(ctx, blog.ID)
	if err == nil {
		resp := s.mapToResponse(createdBlog)
		return &resp, http.StatusCreated, nil
	}

	resp := s.mapToResponse(blog)
	return &resp, http.StatusCreated, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*BlogResponse, int, error) {
	blog, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrBlogNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get blog by ID")
		return nil, http.StatusInternalServerError, errors.New("failed to get blog")
	}

	resp := s.mapToResponse(blog)
	return &resp, http.StatusOK, nil
}

func (s *service) GetBySlug(ctx context.Context, slug string) (*BlogResponse, int, error) {
	blog, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, ErrBlogNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get blog by slug")
		return nil, http.StatusInternalServerError, errors.New("failed to get blog")
	}

	resp := s.mapToResponse(blog)
	return &resp, http.StatusOK, nil
}

func (s *service) GetAll(ctx context.Context, limit, offset int, status string) ([]BlogResponse, int, error) {
	blogs, err := s.repo.FindAll(ctx, limit, offset, status)
	if err != nil {
		logrus.WithError(err).Error("Failed to get all blogs")
		return nil, http.StatusInternalServerError, errors.New("failed to get blogs")
	}

	responses := make([]BlogResponse, len(blogs))
	for i, b := range blogs {
		responses[i] = s.mapToResponse(&b)
	}

	return responses, http.StatusOK, nil
}

func (s *service) GetPublishedBlogs(ctx context.Context, limit, offset int) ([]BlogResponse, int, error) {
	// Re-uses GetAll with 'published' status hardcoded
	return s.GetAll(ctx, limit, offset, "published")
}

func (s *service) GetPublishedBySlug(ctx context.Context, slug string) (*BlogResponse, int, error) {
	blog, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, ErrBlogNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get blog by slug")
		return nil, http.StatusInternalServerError, errors.New("failed to get blog")
	}

	// Double check if it's really published
	if blog.Status != "published" {
		return nil, http.StatusNotFound, ErrBlogNotFound
	}

	resp := s.mapToResponse(blog)
	return &resp, http.StatusOK, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req UpdateBlogRequest) (*BlogResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	blog, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrBlogNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get blog")
	}

	if req.Title != nil {
		blog.Title = *req.Title
	}
	if req.Slug != nil {
		blog.Slug = *req.Slug
	}
	if req.Content != nil {
		blog.Content = *req.Content
	}
	if req.Excerpt != nil {
		blog.Excerpt = *req.Excerpt
	}
	if req.CoverImageURL != nil {
		blog.CoverImageURL = *req.CoverImageURL
	}
	if req.Status != nil {
		// handle transition to published
		if blog.Status != "published" && *req.Status == "published" && blog.PublishedAt == nil {
			now := time.Now()
			blog.PublishedAt = &now
		}
		blog.Status = *req.Status
	}
	if req.CategoryID != nil {
		blog.CategoryID = *req.CategoryID
	}

	blog.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, blog); err != nil {
		logrus.WithError(err).Error("Failed to update blog")
		return nil, http.StatusInternalServerError, errors.New("failed to update blog")
	}

	resp := s.mapToResponse(blog)
	return &resp, http.StatusOK, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) (int, error) {
	if err := s.repo.Delete(ctx, id); err != nil {
		logrus.WithError(err).Error("Failed to delete blog")
		return http.StatusInternalServerError, errors.New("failed to delete blog")
	}
	return http.StatusOK, nil
}

func (s *service) mapToResponse(b *models.Blog) BlogResponse {
	resp := BlogResponse{
		ID:            b.ID,
		Title:         b.Title,
		Slug:          b.Slug,
		Content:       b.Content,
		Excerpt:       b.Excerpt,
		CoverImageURL: b.CoverImageURL,
		Status:        b.Status,
		PublishedAt:   b.PublishedAt,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}

	if b.Category != nil {
		resp.Category = &CategorySummary{
			ID:   b.Category.ID,
			Name: b.Category.Name,
			Slug: b.Category.Slug,
		}
	}

	if b.Author != nil {
		resp.Author = AuthorSummary{
			ID:   b.Author.ID,
			Name: b.Author.Name,
		}
	}

	return resp
}
