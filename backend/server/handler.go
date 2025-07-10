package server

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lucasmcclean/limitlink/link"
)

// LinkHandler routes POST, PATCH, and GET requests to the appropriate handlers.
func LinkHandler(links link.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postLink(w, r, links)
		case http.MethodGet:
			getLink(w, r, links)
		case http.MethodPatch:
			patchLink(w, r, links)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// postLink handles HTTP POST requests to create a new shortened link.
// It expects a JSON body containing all required fields and possibly optional fields.
func postLink(w http.ResponseWriter, r *http.Request, links link.Repository) {
	var validated *link.Validated
	var err error

	validated, err = link.FromJSON(r.Body, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := links.Create(r.Context(), validated); err != nil {
		log.Printf("error storing link: %v", err)
		http.Error(w, "Something went wrong while saving your link. Please try again later.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// patchLink handles PATCH requests for updating a link.
// It expects a JSON body with optional fields to modify, and a Bearer token for authentication.
func patchLink(w http.ResponseWriter, r *http.Request, links link.Repository) {
	adminToken, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	original, err := links.GetByToken(r.Context(), adminToken)
	if err != nil {
		http.Error(w, "Link not found or invalid admin token", http.StatusNotFound)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	patch, err := link.PatchFromJSON(data, original)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := links.PatchByToken(r.Context(), adminToken, patch); err != nil {
		http.Error(w, "Error updating link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getLink returns the current state of a link.
// It expects a Bearer token to authorize the request.
func getLink(w http.ResponseWriter, r *http.Request, links link.Repository) {
	adminToken, err := extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	link, err := links.GetByToken(r.Context(), adminToken)
	if err != nil {
		http.Error(w, "Link not found or invalid admin token", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(link); err != nil {
		http.Error(w, "Error serializing link", http.StatusInternalServerError)
	}
}

// extractBearerToken parses the Authorization header and extracts the Bearer token.
func extractBearerToken(r *http.Request) (string, error) {
	const prefix = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("Missing or invalid Authorization header")
	}
	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return "", errors.New("Missing admin token")
	}
	return token, nil
}
