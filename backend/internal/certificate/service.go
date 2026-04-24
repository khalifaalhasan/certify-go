package certificate

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/cache"
	"github.com/base-go/backend/pkg/validator"
)

const DeliveryQueueKey = "certify:delivery:queue"

// Service defines business logic for the certificate module
type Service interface {
	CreateCertificate(ctx context.Context, createdBy uuid.UUID, req CreateCertificateRequest) (*CertificateResponse, int, error)
	GetCertificateByID(ctx context.Context, id uuid.UUID) (*CertificateResponse, int, error)
	UpdateCertificate(ctx context.Context, id uuid.UUID, req UpdateCertificateRequest) (*CertificateResponse, int, error)
	DeleteCertificate(ctx context.Context, id uuid.UUID) (int, error)
	ListCertificates(ctx context.Context, req ListCertificatesRequest) (*CertificateListResponse, int, error)

	AddRecipients(ctx context.Context, certID uuid.UUID, req AddRecipientsRequest) ([]RecipientResponse, int, error)
	ListRecipients(ctx context.Context, certID uuid.UUID, req ListRecipientsRequest) (*RecipientListResponse, int, error)

	DispatchDeliveries(ctx context.Context, certID uuid.UUID) (*DispatchResponse, int, error)
}

type service struct {
	repo  Repository
	cache cache.Cache
}

func NewService(repo Repository, cache cache.Cache) Service {
	return &service{repo: repo, cache: cache}
}

// calcTotalPages computes the total number of pages given a total item count and page size.
func calcTotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	pages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		pages++
	}
	return pages
}

// --- Certificate CRUD ---

func (s *service) CreateCertificate(ctx context.Context, createdBy uuid.UUID, req CreateCertificateRequest) (*CertificateResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	cert := &models.Certificate{
		Title:       req.Title,
		Description: req.Description,
		TemplateURL: req.TemplateURL,
		EventDate:   req.EventDate,
		IssuedBy:    req.IssuedBy,
		Status:      req.Status,
		CreatedBy:   createdBy,
	}

	if err := s.repo.CreateCertificate(ctx, cert); err != nil {
		logrus.WithError(err).Error("Failed to create certificate")
		return nil, http.StatusInternalServerError, errors.New("failed to create certificate")
	}

	// Fetch with preloaded relations
	created, err := s.repo.GetCertificateByID(ctx, cert.ID)
	if err != nil {
		resp := s.mapCertToResponse(cert, nil)
		return &resp, http.StatusCreated, nil
	}

	resp := s.mapCertToResponse(created, nil)
	return &resp, http.StatusCreated, nil
}

func (s *service) GetCertificateByID(ctx context.Context, id uuid.UUID) (*CertificateResponse, int, error) {
	cert, err := s.repo.GetCertificateByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get certificate")
	}

	stats, _ := s.repo.GetDeliveryStats(ctx, id)
	resp := s.mapCertToResponse(cert, stats)
	return &resp, http.StatusOK, nil
}

func (s *service) UpdateCertificate(ctx context.Context, id uuid.UUID, req UpdateCertificateRequest) (*CertificateResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	cert, err := s.repo.GetCertificateByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get certificate")
	}

	if req.Title != nil {
		cert.Title = *req.Title
	}
	if req.Description != nil {
		cert.Description = *req.Description
	}
	if req.TemplateURL != nil {
		cert.TemplateURL = *req.TemplateURL
	}
	if req.EventDate != nil {
		cert.EventDate = req.EventDate
	}
	if req.IssuedBy != nil {
		cert.IssuedBy = *req.IssuedBy
	}
	if req.Status != nil {
		cert.Status = *req.Status
	}
	cert.UpdatedAt = time.Now()

	if err := s.repo.UpdateCertificate(ctx, cert); err != nil {
		logrus.WithError(err).Error("Failed to update certificate")
		return nil, http.StatusInternalServerError, errors.New("failed to update certificate")
	}

	stats, _ := s.repo.GetDeliveryStats(ctx, id)
	resp := s.mapCertToResponse(cert, stats)
	return &resp, http.StatusOK, nil
}

func (s *service) DeleteCertificate(ctx context.Context, id uuid.UUID) (int, error) {
	if _, err := s.repo.GetCertificateByID(ctx, id); err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, errors.New("failed to get certificate")
	}

	if err := s.repo.DeleteCertificate(ctx, id); err != nil {
		logrus.WithError(err).Error("Failed to delete certificate")
		return http.StatusInternalServerError, errors.New("failed to delete certificate")
	}
	return http.StatusOK, nil
}

func (s *service) ListCertificates(ctx context.Context, req ListCertificatesRequest) (*CertificateListResponse, int, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	certs, total, err := s.repo.ListCertificates(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Failed to list certificates")
		return nil, http.StatusInternalServerError, errors.New("failed to list certificates")
	}

	items := make([]CertificateResponse, len(certs))
	for i, c := range certs {
		items[i] = s.mapCertToResponse(&c, nil)
	}

	return &CertificateListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: calcTotalPages(total, req.PageSize),
	}, http.StatusOK, nil
}

// --- Recipients ---

func (s *service) AddRecipients(ctx context.Context, certID uuid.UUID, req AddRecipientsRequest) ([]RecipientResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Ensure certificate exists
	if _, err := s.repo.GetCertificateByID(ctx, certID); err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get certificate")
	}

	recipients := make([]models.CertificateRecipient, len(req.Recipients))
	for i, r := range req.Recipients {
		recipients[i] = models.CertificateRecipient{
			CertificateID:  certID,
			RecipientName:  r.Name,
			RecipientEmail: r.Email,
			DeliveryStatus: "pending",
		}
	}

	if err := s.repo.CreateRecipients(ctx, recipients); err != nil {
		logrus.WithError(err).Error("Failed to add recipients")
		return nil, http.StatusInternalServerError, errors.New("failed to add recipients")
	}

	responses := make([]RecipientResponse, len(recipients))
	for i, rec := range recipients {
		responses[i] = s.mapRecipientToResponse(&rec)
	}
	return responses, http.StatusCreated, nil
}

func (s *service) ListRecipients(ctx context.Context, certID uuid.UUID, req ListRecipientsRequest) (*RecipientListResponse, int, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	// Ensure certificate exists
	if _, err := s.repo.GetCertificateByID(ctx, certID); err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get certificate")
	}

	recipients, total, err := s.repo.ListRecipients(ctx, certID, req)
	if err != nil {
		logrus.WithError(err).Error("Failed to list recipients")
		return nil, http.StatusInternalServerError, errors.New("failed to list recipients")
	}

	items := make([]RecipientResponse, len(recipients))
	for i, rec := range recipients {
		items[i] = s.mapRecipientToResponse(&rec)
	}

	return &RecipientListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: calcTotalPages(total, req.PageSize),
	}, http.StatusOK, nil
}

// --- Dispatch ---

// DispatchDeliveries enqueues all pending/failed recipients onto the Redis delivery queue.
func (s *service) DispatchDeliveries(ctx context.Context, certID uuid.UUID) (*DispatchResponse, int, error) {
	cert, err := s.repo.GetCertificateByID(ctx, certID)
	if err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get certificate")
	}

	pending, err := s.repo.GetPendingRecipients(ctx, certID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get pending recipients")
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve pending recipients")
	}

	stats, _ := s.repo.GetDeliveryStats(ctx, certID)
	alreadySent := 0
	if stats != nil {
		alreadySent = int(stats.Sent)
	}

	now := time.Now()
	var successfullyQueuedIDs []uuid.UUID

	for i := range pending {
		payload := DeliveryJobPayload{
			RecipientID:      pending[i].ID,
			CertificateID:    certID,
			RecipientName:    pending[i].RecipientName,
			RecipientEmail:   pending[i].RecipientEmail,
			CertificateTitle: cert.Title,
			IssuedBy:         cert.IssuedBy,
			TemplateURL:      cert.TemplateURL,
		}

		data, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			logrus.WithError(marshalErr).Errorf("Failed to marshal payload for recipient %s", pending[i].ID)
			continue
		}

		if pushErr := s.cache.LPush(ctx, DeliveryQueueKey, string(data)); pushErr != nil {
			logrus.WithError(pushErr).Errorf("Failed to enqueue recipient %s", pending[i].ID)
			continue
		}

		successfullyQueuedIDs = append(successfullyQueuedIDs, pending[i].ID)
	}

	// Batch-update all successfully queued recipients to "queued" status in one query
	if len(successfullyQueuedIDs) > 0 {
		if updateErr := s.repo.BatchMarkQueued(ctx, successfullyQueuedIDs, now); updateErr != nil {
			logrus.WithError(updateErr).Error("Failed to batch-update queued recipient statuses")
		}
	}

	return &DispatchResponse{
		CertificateID: certID,
		TotalQueued:   len(successfullyQueuedIDs),
		AlreadySent:   alreadySent,
		Message:       "Certificates dispatched to delivery queue",
	}, http.StatusOK, nil
}

// --- Mappers ---

func (s *service) mapCertToResponse(cert *models.Certificate, stats *DeliveryStats) CertificateResponse {
	resp := CertificateResponse{
		ID:          cert.ID,
		Title:       cert.Title,
		Description: cert.Description,
		TemplateURL: cert.TemplateURL,
		EventDate:   cert.EventDate,
		IssuedBy:    cert.IssuedBy,
		Status:      cert.Status,
		CreatedBy:   cert.CreatedBy,
		CreatedAt:   cert.CreatedAt,
		UpdatedAt:   cert.UpdatedAt,
		Stats:       stats,
	}
	if cert.Creator != nil {
		resp.Creator = &CreatorSummary{
			ID:   cert.Creator.ID,
			Name: cert.Creator.Name,
		}
	}
	return resp
}

func (s *service) mapRecipientToResponse(rec *models.CertificateRecipient) RecipientResponse {
	return RecipientResponse{
		ID:             rec.ID,
		CertificateID:  rec.CertificateID,
		RecipientName:  rec.RecipientName,
		RecipientEmail: rec.RecipientEmail,
		DeliveryStatus: rec.DeliveryStatus,
		QueuedAt:       rec.QueuedAt,
		SentAt:         rec.SentAt,
		ErrorMessage:   rec.ErrorMessage,
		RetryCount:     rec.RetryCount,
		CreatedAt:      rec.CreatedAt,
	}
}
