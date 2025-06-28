package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/lucasmcclean/limitlink/assets"
	"github.com/lucasmcclean/limitlink/link"
)

func Redirect(repo link.Repository) http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(assets.TemplateFS(), "password.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/")
		if slug == "" {
			http.Error(w, "missing slug", http.StatusBadRequest)
			return
		}

		lnk, err := repo.GetBySlug(r.Context(), slug)
		if err != nil {
			http.Error(w, "link not found", http.StatusNotFound)
			return
		}

		if lnk.PasswordHash != nil {
			switch r.Method {
			case http.MethodGet:
				repo.IncBySlug(r.Context(), slug)
				tmpl.Execute(w, lnk)
				return

			case http.MethodPost:
				err := r.ParseForm()
				if err != nil {
					http.Error(w, "Invalid form", http.StatusBadRequest)
					return
				}

				password := r.FormValue("password")
				valid, err := link.VerifyPassword(*lnk.PasswordHash, password)
				if !valid {
					http.Error(w, "invalid password", http.StatusUnauthorized)
					return
				}
				if err != nil {
					http.Error(w, "error processing password", http.StatusInternalServerError)
					return
				}
			}
		}

		target := lnk.Target
		if target == "" {
			http.Error(w, "invalid target", http.StatusInternalServerError)
		}

		repo.IncBySlug(r.Context(), slug)
		http.Redirect(w, r, target, http.StatusFound)
	}
}
