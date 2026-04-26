package certificate

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// REQUESTS
// ============================================================================

type CreateCertificateRequest struct {
	Title       string     `json:"title" validate:"required,min=3,max=255"`
	Description string     `json:"description" validate:"omitempty"`
	TemplateURL string     `json:"template_url" validate:"omitempty,url"`
	EventDate   *time.Time `json:"event_date" validate:"omitempty"`
	IssuedBy    string     `json:"issued_by" validate:"required,min=2,max=255"`
	Status      string     `json:"status" validate:"required,oneof=draft active archived"`
}

type UpdateCertificateRequest struct {
	Title       *string    `json:"title" validate:"omitempty,min=3,max=255"`
	Description *string    `json:"description" validate:"omitempty"`
	TemplateURL *string    `json:"template_url" validate:"omitempty,url"`
	EventDate   *time.Time `json:"event_date" validate:"omitempty"`
	IssuedBy    *string    `json:"issued_by" validate:"omitempty,min=2,max=255"`
	Status      *string    `json:"status" validate:"omitempty,oneof=draft active archived"`
}

type AddRecipientsRequest struct {
	Recipients []RecipientInput `json:"recipients" validate:"required,min=1,dive"`
}

type RecipientInput struct {
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Email string `json:"email" validate:"required,email"`
}

type ListCertificatesRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
	Status   string `json:"status"`
}

type ListRecipientsRequest struct {
	Page           int    `json:"page"`
	PageSize       int    `json:"page_size"`
	DeliveryStatus string `json:"delivery_status"`
}

// ============================================================================
// RESPONSES
// ============================================================================

type CertificateResponse struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	TemplateURL string          `json:"template_url"`
	EventDate   *time.Time      `json:"event_date,omitempty"`
	IssuedBy    string          `json:"issued_by"`
	Status      string          `json:"status"`
	CreatedBy   uuid.UUID       `json:"created_by"`
	Creator     *CreatorSummary `json:"creator,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Stats       *DeliveryStats  `json:"stats,omitempty"`
}

type CreatorSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type DeliveryStats struct {
	Total      int64 `json:"total"`
	Pending    int64 `json:"pending"`
	Queued     int64 `json:"queued"`
	Processing int64 `json:"processing"`
	Sent       int64 `json:"sent"`
	Failed     int64 `json:"failed"`
}

type CertificateListResponse struct {
	Items      []CertificateResponse `json:"items"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

type RecipientResponse struct {
	ID             uuid.UUID  `json:"id"`
	CertificateID  uuid.UUID  `json:"certificate_id"`
	RecipientName  string     `json:"recipient_name"`
	RecipientEmail string     `json:"recipient_email"`
	DeliveryStatus string     `json:"delivery_status"`
	QueuedAt       *time.Time `json:"queued_at,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	ErrorMessage   *string    `json:"error_message,omitempty"`
	RetryCount     int        `json:"retry_count"`
	CreatedAt      time.Time  `json:"created_at"`
}

type RecipientListResponse struct {
	Items      []RecipientResponse `json:"items"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

type DispatchResponse struct {
	CertificateID uuid.UUID `json:"certificate_id"`
	TotalQueued   int       `json:"total_queued"`
	AlreadySent   int       `json:"already_sent"`
	Message       string    `json:"message"`
}

// DeliveryJobPayload is the JSON structure pushed onto the Redis queue
type DeliveryJobPayload struct {
	RecipientID      uuid.UUID `json:"recipient_id"`
	CertificateID    uuid.UUID `json:"certificate_id"`
	RecipientName    string    `json:"recipient_name"`
	RecipientEmail   string    `json:"recipient_email"`
	CertificateTitle string    `json:"certificate_title"`
	IssuedBy         string    `json:"issued_by"`
	TemplateURL      string    `json:"template_url"`
}
