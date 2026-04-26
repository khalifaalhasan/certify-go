package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/certificate"
	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/cache"
	"github.com/base-go/backend/pkg/config"
	"github.com/base-go/backend/pkg/database"
)

// CertificateWorker processes certificate delivery jobs from the Redis queue
type CertificateWorker struct {
	cache cache.Cache
	db    database.Database
}

// NewCertificateWorker creates a new CertificateWorker
func NewCertificateWorker(cache cache.Cache, db database.Database) *CertificateWorker {
	return &CertificateWorker{cache: cache, db: db}
}

// Start launches the worker loop in a goroutine and returns immediately.
// It blocks on BRPop to avoid busy-waiting.
func (w *CertificateWorker) Start(ctx context.Context) {
	go func() {
		logrus.Info("CertificateWorker: started, listening on queue: ", certificate.DeliveryQueueKey)
		for {
			select {
			case <-ctx.Done():
				logrus.Info("CertificateWorker: context cancelled, shutting down")
				return
			default:
				w.processNext(ctx)
			}
		}
	}()
}

// processNext blocks waiting for the next job on the Redis queue (5-second timeout),
// then processes it.
func (w *CertificateWorker) processNext(ctx context.Context) {
	result, err := w.cache.BRPop(ctx, 5*time.Second, certificate.DeliveryQueueKey)
	if err != nil {
		logrus.WithError(err).Error("CertificateWorker: BRPop error")
		return
	}
	if result == nil {
		// timeout, no job available – loop again
		return
	}

	// BRPop returns [key, value]
	if len(result) < 2 {
		return
	}
	raw := result[1]

	var payload certificate.DeliveryJobPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		logrus.WithError(err).Error("CertificateWorker: failed to unmarshal job payload")
		return
	}

	logrus.Infof("CertificateWorker: processing delivery for recipient %s (%s)",
		payload.RecipientName, payload.RecipientEmail)

	w.markProcessing(ctx, payload.RecipientID.String())

	if err := w.sendEmail(payload); err != nil {
		logrus.WithError(err).Errorf("CertificateWorker: failed to send email to %s", payload.RecipientEmail)
		w.markFailed(ctx, payload.RecipientID.String(), err.Error())
		return
	}

	logrus.Infof("CertificateWorker: email sent to %s", payload.RecipientEmail)
	w.markSent(ctx, payload.RecipientID.String())
}

// sendEmail sends the certificate email via SMTP.
// If SMTP is not configured the delivery is logged and treated as successful (no-op).
func (w *CertificateWorker) sendEmail(payload certificate.DeliveryJobPayload) error {
	cfg := config.GetConfig()
	if cfg.Email.Host == "" {
		logrus.Infof("CertificateWorker: SMTP not configured, skipping send to %s", payload.RecipientEmail)
		return nil
	}

	subject := fmt.Sprintf("Your Certificate: %s", payload.CertificateTitle)
	body := buildEmailBody(payload)

	msg := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"From: " + cfg.Email.From + "\r\n" +
		"To: " + payload.RecipientEmail + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body

	addr := fmt.Sprintf("%s:%s", cfg.Email.Host, cfg.Email.Port)
	var auth smtp.Auth
	if cfg.Email.Username != "" {
		auth = smtp.PlainAuth("", cfg.Email.Username, cfg.Email.Password, cfg.Email.Host)
	}

	return smtp.SendMail(addr, auth, cfg.Email.From, []string{payload.RecipientEmail}, []byte(msg))
}

func buildEmailBody(p certificate.DeliveryJobPayload) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Dear %s,\n\n", p.RecipientName))
	sb.WriteString("Congratulations! You have been awarded the certificate:\n\n")
	sb.WriteString(fmt.Sprintf("  Title    : %s\n", p.CertificateTitle))
	sb.WriteString(fmt.Sprintf("  Issued By: %s\n", p.IssuedBy))
	if p.TemplateURL != "" {
		sb.WriteString(fmt.Sprintf("  View     : %s\n", p.TemplateURL))
	}
	sb.WriteString("\nBest regards,\n")
	sb.WriteString(p.IssuedBy)
	return sb.String()
}

// --- Database status helpers ---

func (w *CertificateWorker) markProcessing(ctx context.Context, recipientID string) {
	now := time.Now()
	if err := w.db.GetDB().WithContext(ctx).
		Model(&models.CertificateRecipient{}).
		Where("id = ?", recipientID).
		Updates(map[string]interface{}{
			"delivery_status": "processing",
			"updated_at":      now,
		}).Error; err != nil {
		logrus.WithError(err).Errorf("CertificateWorker: failed to mark recipient %s as processing", recipientID)
	}
}

func (w *CertificateWorker) markSent(ctx context.Context, recipientID string) {
	now := time.Now()
	if err := w.db.GetDB().WithContext(ctx).
		Model(&models.CertificateRecipient{}).
		Where("id = ?", recipientID).
		Updates(map[string]interface{}{
			"delivery_status": "sent",
			"sent_at":         now,
			"updated_at":      now,
		}).Error; err != nil {
		logrus.WithError(err).Errorf("CertificateWorker: failed to mark recipient %s as sent", recipientID)
	}
}

func (w *CertificateWorker) markFailed(ctx context.Context, recipientID, errMsg string) {
	now := time.Now()
	if err := w.db.GetDB().WithContext(ctx).
		Model(&models.CertificateRecipient{}).
		Where("id = ?", recipientID).
		Updates(map[string]interface{}{
			"delivery_status": "failed",
			"error_message":   errMsg,
			"retry_count":     gorm.Expr("retry_count + 1"),
			"updated_at":      now,
		}).Error; err != nil {
		logrus.WithError(err).Errorf("CertificateWorker: failed to mark recipient %s as failed", recipientID)
	}
}
