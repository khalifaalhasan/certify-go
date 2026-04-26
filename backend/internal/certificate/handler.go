package certificate

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/base-go/backend/pkg/response"
)

// Handler holds HTTP handlers for the certificate module
type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

// Create godoc
// POST /v1/admin/certificates
func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := r.Context().Value("user_context").(response.UserContext)
	if !ok {
		response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	createdBy, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req CreateCertificateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, statusCode, err := h.service.CreateCertificate(r.Context(), createdBy, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, resp)
}

// GetByID godoc
// GET /v1/admin/certificates/{id}
func (h Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid certificate ID")
		return
	}

	resp, statusCode, err := h.service.GetCertificateByID(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, resp)
}

// Update godoc
// PUT /v1/admin/certificates/{id}
func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid certificate ID")
		return
	}

	var req UpdateCertificateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, statusCode, err := h.service.UpdateCertificate(r.Context(), id, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, resp)
}

// Delete godoc
// DELETE /v1/admin/certificates/{id}
func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid certificate ID")
		return
	}

	statusCode, err := h.service.DeleteCertificate(r.Context(), id)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, map[string]string{"message": "Certificate deleted successfully"})
}

// List godoc
// GET /v1/admin/certificates
func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	req := ListCertificatesRequest{
		Page:     parseQueryInt(r.URL.Query().Get("page"), 1),
		PageSize: parseQueryInt(r.URL.Query().Get("page_size"), 10),
		Search:   r.URL.Query().Get("search"),
		Status:   r.URL.Query().Get("status"),
	}

	resp, statusCode, err := h.service.ListCertificates(r.Context(), req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, resp)
}

// AddRecipients godoc
// POST /v1/admin/certificates/{id}/recipients
func (h Handler) AddRecipients(w http.ResponseWriter, r *http.Request) {
	certID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid certificate ID")
		return
	}

	var req AddRecipientsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, statusCode, err := h.service.AddRecipients(r.Context(), certID, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, resp)
}

// ListRecipients godoc
// GET /v1/admin/certificates/{id}/recipients
func (h Handler) ListRecipients(w http.ResponseWriter, r *http.Request) {
	certID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid certificate ID")
		return
	}

	req := ListRecipientsRequest{
		Page:           parseQueryInt(r.URL.Query().Get("page"), 1),
		PageSize:       parseQueryInt(r.URL.Query().Get("page_size"), 10),
		DeliveryStatus: r.URL.Query().Get("delivery_status"),
	}

	resp, statusCode, err := h.service.ListRecipients(r.Context(), certID, req)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, resp)
}

// Dispatch godoc
// POST /v1/admin/certificates/{id}/dispatch
func (h Handler) Dispatch(w http.ResponseWriter, r *http.Request) {
	certID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid certificate ID")
		return
	}

	resp, statusCode, err := h.service.DispatchDeliveries(r.Context(), certID)
	if err != nil {
		response.ResponseError(w, statusCode, err.Error())
		return
	}

	response.ResponseJSON(w, statusCode, resp)
}

func parseQueryInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 {
		return defaultVal
	}
	return v
}
