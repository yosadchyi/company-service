package auth

import (
	"context"
	"encoding/json"
	"net/http"
)

type Service interface {
	Login(ctx context.Context, username, password string) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Login and receive JWT tokens
// @Summary Login and receive JWT tokens
// @Accept json
// @Produce json
// @Param credentials body loginRequest true "Username and Password"
// @Success 200 {object} loginResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	access, refresh, err := h.svc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		println(err.Error())
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(loginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// Refresh refreshes access token
// @Summary Refresh access token
// @Accept json
// @Produce json
// @Param refreshToken body refreshRequest true "Refresh Token"
// @Success 200 {object} map[string]string
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"access_token": token})
}
