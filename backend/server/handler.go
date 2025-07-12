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

// RedirectHandler redirects GET requests to their matching target.
// It will first verify that the link is available and fail if it can't
// increment the hit count.
func RedirectHandler(links link.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		path := strings.Trim(r.URL.Path, "/")
		slug := strings.Split(path, "/")[0]

		if slug == "" {
			http.Error(w, "Missing slug", http.StatusBadRequest)
			return
		}

		if len(slug) < link.MinSlugLen || len(slug) > link.MaxSlugLen {
			http.Error(w, "Invalid slug length", http.StatusBadRequest)
			return
		}

		lnk, err := links.GetBySlug(r.Context(), slug)
		if err != nil || lnk == nil {
			http.Error(w, "Link not found", http.StatusNotFound)
			return
		}

		if !lnk.IsAvailable(time.Now()) {
			http.Error(w, "Link not found", http.StatusNotFound)
			return
		}

		if lnk.PasswordHash != nil {
			password := r.Header.Get("X-Link-Password")
			if password == "" {
				http.Error(w, "Password required", http.StatusUnauthorized)
				return
			}
			valid, err := lnk.IsCorrectPassword(password)
			if err != nil {
				http.Error(w, "Error validating password", http.StatusInternalServerError)
				return
			}
			if !valid {
				http.Error(w, "Invalid password", http.StatusUnauthorized)
				return
			}
		}

		err = links.IncBySlug(r.Context(), slug)
		if err != nil {
			http.Error(w, "Error retrieving link", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, lnk.Target, http.StatusFound)
	}
}

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

	lnk := validated.Link()

	resp := struct {
		Slug        string `json:"slug"`
		AdminToken  string `json:"admin_token"`
		RedirectURL string `json:"redirect_url"`
		AdminURL    string `json:"admin_url"`
	}{
		Slug:        lnk.Slug,
		AdminToken:  lnk.AdminToken,
		RedirectURL: "https://limitl.ink/" + lnk.Slug,
		AdminURL:    "https://limitl.ink/admin/" + lnk.AdminToken,
	}

	w.Header().Set("Location", "/admin/"+lnk.AdminToken)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error encoding created link: %v", err)
		http.Error(w, "Error encoding created link", http.StatusInternalServerError)
	}
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
		return "", errors.New("missing or invalid Authorization header")
	}
	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return "", errors.New("missing admin token")
	}
	return token, nil
}
