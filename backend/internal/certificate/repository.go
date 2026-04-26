package certificate

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/database"
)

var (
	ErrCertificateNotFound = errors.New("certificate not found")
	ErrRecipientNotFound   = errors.New("recipient not found")
)

// Repository defines all database operations for the certificate module
type Repository interface {
	// Certificate CRUD
	CreateCertificate(ctx context.Context, cert *models.Certificate) error
	GetCertificateByID(ctx context.Context, id uuid.UUID) (*models.Certificate, error)
	UpdateCertificate(ctx context.Context, cert *models.Certificate) error
	DeleteCertificate(ctx context.Context, id uuid.UUID) error
	ListCertificates(ctx context.Context, req ListCertificatesRequest) ([]models.Certificate, int64, error)

	// Recipient CRUD
	CreateRecipients(ctx context.Context, recipients []models.CertificateRecipient) error
	GetRecipientByID(ctx context.Context, id uuid.UUID) (*models.CertificateRecipient, error)
	ListRecipients(ctx context.Context, certID uuid.UUID, req ListRecipientsRequest) ([]models.CertificateRecipient, int64, error)
	UpdateRecipient(ctx context.Context, recipient *models.CertificateRecipient) error

	// Delivery queue operations
	GetPendingRecipients(ctx context.Context, certID uuid.UUID) ([]models.CertificateRecipient, error)
	GetDeliveryStats(ctx context.Context, certID uuid.UUID) (*DeliveryStats, error)
	BatchMarkQueued(ctx context.Context, ids []uuid.UUID, queuedAt time.Time) error
}

type repository struct {
	db database.Database
}

func NewRepository(db database.Database) Repository {
	return &repository{db: db}
}

// --- Certificate CRUD ---

func (r *repository) CreateCertificate(ctx context.Context, cert *models.Certificate) error {
	return r.db.GetDB().WithContext(ctx).Create(cert).Error
}

func (r *repository) GetCertificateByID(ctx context.Context, id uuid.UUID) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.db.GetDB().WithContext(ctx).
		Preload("Creator").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&cert).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCertificateNotFound
		}
		return nil, err
	}
	return &cert, nil
}

func (r *repository) UpdateCertificate(ctx context.Context, cert *models.Certificate) error {
	return r.db.GetDB().WithContext(ctx).Save(cert).Error
}

func (r *repository) DeleteCertificate(ctx context.Context, id uuid.UUID) error {
	return r.db.GetDB().WithContext(ctx).
		Model(&models.Certificate{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *repository) ListCertificates(ctx context.Context, req ListCertificatesRequest) ([]models.Certificate, int64, error) {
	query := r.db.GetDB().WithContext(ctx).
		Model(&models.Certificate{}).
		Preload("Creator").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		pattern := "%" + req.Search + "%"
		query = query.Where("title ILIKE ? OR issued_by ILIKE ?", pattern, pattern)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	var certs []models.Certificate
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&certs).Error; err != nil {
		return nil, 0, err
	}

	return certs, total, nil
}

// --- Recipient CRUD ---

func (r *repository) CreateRecipients(ctx context.Context, recipients []models.CertificateRecipient) error {
	return r.db.GetDB().WithContext(ctx).Create(&recipients).Error
}

func (r *repository) GetRecipientByID(ctx context.Context, id uuid.UUID) (*models.CertificateRecipient, error) {
	var rec models.CertificateRecipient
	err := r.db.GetDB().WithContext(ctx).
		Where("id = ?", id).
		First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecipientNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func (r *repository) ListRecipients(ctx context.Context, certID uuid.UUID, req ListRecipientsRequest) ([]models.CertificateRecipient, int64, error) {
	query := r.db.GetDB().WithContext(ctx).
		Model(&models.CertificateRecipient{}).
		Where("certificate_id = ?", certID)

	if req.DeliveryStatus != "" {
		query = query.Where("delivery_status = ?", req.DeliveryStatus)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	var recipients []models.CertificateRecipient
	if err := query.Order("created_at ASC").Offset(offset).Limit(req.PageSize).Find(&recipients).Error; err != nil {
		return nil, 0, err
	}

	return recipients, total, nil
}

func (r *repository) UpdateRecipient(ctx context.Context, recipient *models.CertificateRecipient) error {
	return r.db.GetDB().WithContext(ctx).Save(recipient).Error
}

// --- Delivery helpers ---

func (r *repository) GetPendingRecipients(ctx context.Context, certID uuid.UUID) ([]models.CertificateRecipient, error) {
	var recipients []models.CertificateRecipient
	err := r.db.GetDB().WithContext(ctx).
		Where("certificate_id = ? AND delivery_status IN ('pending', 'failed')", certID).
		Find(&recipients).Error
	return recipients, err
}

func (r *repository) GetDeliveryStats(ctx context.Context, certID uuid.UUID) (*DeliveryStats, error) {
	type row struct {
		Status string
		Count  int64
	}
	var rows []row
	err := r.db.GetDB().WithContext(ctx).
		Model(&models.CertificateRecipient{}).
		Select("delivery_status as status, COUNT(*) as count").
		Where("certificate_id = ?", certID).
		Group("delivery_status").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	stats := &DeliveryStats{}
	for _, rw := range rows {
		stats.Total += rw.Count
		switch rw.Status {
		case "pending":
			stats.Pending = rw.Count
		case "queued":
			stats.Queued = rw.Count
		case "processing":
			stats.Processing = rw.Count
		case "sent":
			stats.Sent = rw.Count
		case "failed":
			stats.Failed = rw.Count
		}
	}
	return stats, nil
}

func (r *repository) BatchMarkQueued(ctx context.Context, ids []uuid.UUID, queuedAt time.Time) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDB().WithContext(ctx).
		Model(&models.CertificateRecipient{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"delivery_status": "queued",
			"queued_at":       queuedAt,
			"updated_at":      queuedAt,
		}).Error
}
