package models

import (
	"time"

	"github.com/google/uuid"
)

// Certificate represents a certificate template / issuance event
type Certificate struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	TemplateURL string     `gorm:"type:text" json:"template_url"`
	EventDate   *time.Time `gorm:"type:timestamp with time zone" json:"event_date"`
	IssuedBy    string     `gorm:"type:varchar(255);not null" json:"issued_by"`
	Status      string     `gorm:"type:varchar(50);default:'draft'" json:"status"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Creator    *User                  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Recipients []CertificateRecipient `gorm:"foreignKey:CertificateID" json:"recipients,omitempty"`
}

func (Certificate) TableName() string {
	return "certificates"
}

// CertificateRecipient represents one recipient of a certificate delivery
type CertificateRecipient struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CertificateID  uuid.UUID  `gorm:"type:uuid;not null" json:"certificate_id"`
	RecipientName  string     `gorm:"type:varchar(255);not null" json:"recipient_name"`
	RecipientEmail string     `gorm:"type:varchar(255);not null" json:"recipient_email"`
	DeliveryStatus string     `gorm:"type:varchar(50);default:'pending'" json:"delivery_status"`
	QueuedAt       *time.Time `gorm:"type:timestamp with time zone" json:"queued_at,omitempty"`
	SentAt         *time.Time `gorm:"type:timestamp with time zone" json:"sent_at,omitempty"`
	ErrorMessage   *string    `gorm:"type:text" json:"error_message,omitempty"`
	RetryCount     int        `gorm:"default:0" json:"retry_count"`
	CreatedAt      time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"not null;default:now()" json:"updated_at"`

	Certificate *Certificate `gorm:"foreignKey:CertificateID" json:"certificate,omitempty"`
}

func (CertificateRecipient) TableName() string {
	return "certificate_recipients"
}
