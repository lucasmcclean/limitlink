package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/lucasmcclean/limitlink/assets"
	"github.com/lucasmcclean/limitlink/link"
)

func Admin(repo link.Repository) http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(assets.TemplateFS(), "admin.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.URL.Path, "/admin/")
		if token == "" {
			http.Error(w, "missing admin token", http.StatusBadRequest)
			return
		}

		lnk, err := repo.GetByToken(r.Context(), token)
		if err != nil {
			http.Error(w, "link not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, lnk.IntoDisplay())
		if err != nil {
			log.Printf("template render error: %v", err)
			http.Error(w, "failed to render template", http.StatusInternalServerError)
		}
	}
}
