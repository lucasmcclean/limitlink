package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lucasmcclean/limitlink/link"
)

func Link(links link.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postLink(w, r, links)
		case http.MethodPatch:
			patchLink(w, r, links)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func postLink(w http.ResponseWriter, r *http.Request, links link.Repository) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "There was a problem processing your form. Please try again.", http.StatusBadRequest)
		return
	}

	validated, err := link.FromForm(r.PostForm, time.Now())
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

func patchLink(w http.ResponseWriter, r *http.Request, links link.Repository) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "There was a problem processing your form. Please try again.", http.StatusBadRequest)
		return
	}

	// Expecting header: Authorization: Bearer <admin-token>
	authHeader := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	adminToken := strings.TrimPrefix(authHeader, prefix)
	if adminToken == "" {
		http.Error(w, "Missing admin token", http.StatusUnauthorized)
		return
	}

	original, err := links.GetByToken(r.Context(), adminToken)
	if err != nil {
		http.Error(w, "Link not found or invalid admin token", http.StatusNotFound)
		return
	}

	patch, err := link.PatchFromForm(r.PostForm, original)
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
