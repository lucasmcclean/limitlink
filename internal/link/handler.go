package link

import (
	"encoding/json"
	"net/http"
	"github.com/lucasmcclean/url-shortener/internal/logger"
	"github.com/google/uuid"
)

// Handler encapsulates the HTTP handlers for managing links
type Handler struct {
	Service *Service
	Log     logger.Logger
}

// NewHandler creates a new Handler instance
func NewHandler(service *Service, log logger.Logger) *Handler {
	return &Handler{
		Service: service,
		Log:     log,
	}
}

// Create handles the creation of a new shortened link
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log.Error("Failed to decode request", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	link, err := h.Service.CreateLink(r.Context(), req)
	if err != nil {
		h.Log.Error("Failed to create link", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(link); err != nil {
		h.Log.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Visit handles visiting a shortened link
func (h *Handler) Visit(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path[len("/v/"):]

	original, err := h.Service.VisitLink(r.Context(), short)
	if err != nil {
		h.Log.Error("Failed to resolve link", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}

// AdminView allows the admin to view a link by its admin token
func (h *Handler) AdminView(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Path[len("/admin/"):]

	token, err := uuid.Parse(tokenStr)
	if err != nil {
		h.Log.Error("Invalid admin token", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	link, err := h.Service.GetByAdminToken(r.Context(), token)
	if err != nil {
		h.Log.Error("Failed to fetch link by admin token", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(link); err != nil {
		h.Log.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Delete removes a link by its admin token
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Path[len("/delete/"):]

	token, err := uuid.Parse(tokenStr)
	if err != nil {
		h.Log.Error("Invalid admin token", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteByAdminToken(r.Context(), token); err != nil {
		h.Log.Error("Failed to delete link by admin token", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes registers the routes with the http.ServeMux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create", h.Create)
	mux.HandleFunc("/v/", h.Visit)
	mux.HandleFunc("/admin/", h.AdminView)
	mux.HandleFunc("/delete/", h.Delete)
}
