package company

import (
	"context"
	"encoding/json"
	"net/http"

	"company-service/pkg/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Service defines company service operations
type Service interface {
	Create(ctx context.Context, c *Company) error
	Get(ctx context.Context, id uuid.UUID) (*Company, error)
	Update(ctx context.Context, id uuid.UUID, patch *PartialCompany) (*Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// Handler handles HTTP requests for the Company entity.
type Handler struct {
	svc Service
}

// NewHandler creates a new Handler with a given Company service.
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes sets up the company-related routes on the given router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/companies", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Patch("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// Create handles the creation of a new company.
// @Summary Create a new company
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param company body Company true "Company to create"
// @Success 201 {object} Company
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /companies [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var c Company
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ID = uuid.New()
	if err := h.svc.Create(r.Context(), &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kafka.PublishEvent("company_created", c)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}

// Get retrieves a company by ID.
// @Summary Get a company
// @Security BearerAuth
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} Company
// @Failure 404 {string} string "Not found"
// @Failure 401 {string} string "Unauthorized"
// @Router /companies/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	c, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(c)
}

// Update Partially update a company
// @Summary Partially update a company
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param patch body PartialCompany true "Partial company fields"
// @Success 200 {object} Company
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Not found"
// @Failure 401 {string} string "Unauthorized"
// @Router /companies/{id} [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	var c PartialCompany
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updated *Company

	if u, err := h.svc.Update(r.Context(), id, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		updated = u
	}
	kafka.PublishEvent("company_updated", updated)
	_ = json.NewEncoder(w).Encode(updated)
}

// Delete removes a company by ID.
// @Summary Delete a company
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 204
// @Failure 404 {string} string "Not found"
// @Failure 401 {string} string "Unauthorized"
// @Router /companies/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kafka.PublishEvent("company_deleted", map[string]string{"id": id.String()})
	w.WriteHeader(http.StatusNoContent)
}
