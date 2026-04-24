package category

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// REQUESTS
// ============================================================================

// CreateCategoryRequest digunakan saat admin membuat kategori baru.
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Slug        string `json:"slug" validate:"required,min=2,max=100"`
	// Description tidak wajib, jadi kita hilangkan tag 'required'
	Description string `json:"description" validate:"omitempty,max=500"`
}

// UpdateCategoryRequest menggunakan Pointer (*) agar mendukung method PATCH.
// Jika frontend hanya mengirim "name", maka field "slug" dan "description" 
// akan bernilai nil dan tidak ikut ter-update di database.
type UpdateCategoryRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=100"`
	Slug        *string `json:"slug" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// CategoryIDParamRequest untuk menangkap ID dari parameter URL (misal: /categories/:id)
type CategoryIDParamRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}

// ============================================================================
// RESPONSES
// ============================================================================

// CategoryResponse adalah format balikan standar untuk data kategori tunggal
// maupun array dari kategori (List).
type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"` // omitempty agar rapi jika kosong
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}