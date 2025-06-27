package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/lucasmcclean/limitlink/assets"
	"github.com/lucasmcclean/limitlink/link"
)

func Links(repo link.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postLinks(w, r, repo)
		case http.MethodPatch:
			patchLinks(w, r, repo)
		default:
			http.Error(w, "Sorry, this action is not allowed.", http.StatusMethodNotAllowed)
		}
	}
}

func postLinks(w http.ResponseWriter, r *http.Request, repo link.Repository) {
	tmpl := template.Must(template.ParseFS(assets.TemplateFS(), "link-info.partial.html"))

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "There was a problem processing your form. Please try again.", http.StatusBadRequest)
		return
	}

	lnk, err := link.FromForm(r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repo.Create(r.Context(), lnk); err != nil {
		log.Printf("error storing link: %v", err)
		http.Error(w, "Something went wrong while saving your link. Please try again later.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := tmpl.ExecuteTemplate(w, "new-link", lnk.IntoDisplay()); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "There was an error loading the page. Please try again.", http.StatusInternalServerError)
	}
}

func patchLinks(w http.ResponseWriter, r *http.Request, repo link.Repository) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Problem processing form", http.StatusBadRequest)
		return
	}

	adminToken := r.FormValue("admin-token")
	if adminToken == "" {
		http.Error(w, "Missing admin token", http.StatusUnauthorized)
		return
	}

	patch, err := link.PatchFromForm(r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repo.PatchByToken(r.Context(), adminToken, patch); err != nil {
		http.Error(w, "Error updating link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
