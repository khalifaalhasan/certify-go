package blog

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// REQUESTS
// ============================================================================

type CreateBlogRequest struct {
	Title         string    `json:"title" validate:"required,min=3,max=255"`
	Slug          string    `json:"slug" validate:"required,min=3,max=255"`
	Content       string    `json:"content" validate:"required,min=3"`
	Excerpt       string    `json:"excerpt" validate:"required,min=3"`
	CoverImageURL string    `json:"cover_image_url" validate:"required,url"`
	Status        string    `json:"status" validate:"required,oneof=draft published archived"`
	CategoryID    uuid.UUID `json:"category_id" validate:"required"`
	// AuthorID sengaja ditiadakan, akan diinjeksi dari JWT Context di Handler
}

// UpdateBlogRequest menggunakan Pointer (*) untuk mendeteksi field mana yang diubah.
// Jika nilainya nil, berarti frontend tidak mengirim field tersebut (PATCH behavior).
type UpdateBlogRequest struct {
	Title         *string    `json:"title" validate:"omitempty,min=3,max=255"`
	Slug          *string    `json:"slug" validate:"omitempty,min=3,max=255"`
	Content       *string    `json:"content" validate:"omitempty,min=3"`
	Excerpt       *string    `json:"excerpt" validate:"omitempty,min=3"`
	CoverImageURL *string    `json:"cover_image_url" validate:"omitempty,url"`
	Status        *string    `json:"status" validate:"omitempty,oneof=draft published archived"`
	CategoryID    *uuid.UUID `json:"category_id" validate:"omitempty"`
}

// Get dan Delete biasanya mengambil ID dari URL Param (ex: /blogs/:id), 
// tag `param` digunakan jika Anda memakai validator/binder dari framework seperti Echo/Gin/Fiber.
type BlogIDParamRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}

// ============================================================================
// RESPONSES
// ============================================================================

// BlogResponse adalah single source of truth untuk balikan data blog
type BlogResponse struct {
	ID            uuid.UUID        `json:"id"`
	Title         string           `json:"title"`
	Slug          string           `json:"slug"`
	Content       string           `json:"content"`
	Excerpt       string           `json:"excerpt"`
	CoverImageURL string           `json:"cover_image_url"`
	Status        string           `json:"status"`
	PublishedAt   *time.Time       `json:"published_at,omitempty"` // Nullable dari database
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	
	// Nested Response untuk Relasi (memudahkan Frontend)
	Category      *CategorySummary `json:"category,omitempty"`
	Author        AuthorSummary    `json:"author"`
}

// CategorySummary versi ringkas untuk nested response
type CategorySummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

// AuthorSummary versi ringkas, JANGAN sertakan password atau email di sini!
type AuthorSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"` // Sesuaikan dengan field nama di tabel users Anda
}