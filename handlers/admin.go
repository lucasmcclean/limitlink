package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/lucasmcclean/limitlink/link"
)

var tmpl = template.Must(template.ParseFiles("templates/admin.html"))

func Admin(repo link.Repository) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
    adminToken := strings.TrimPrefix(r.URL.Path, "/admin/")
		if adminToken == "" {
			http.Error(w, "missing admin token", http.StatusBadRequest)
			return
		}

		lnk, err := repo.GetByToken(r.Context(), adminToken)
		if err != nil {
			http.Error(w, "link not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, lnk)
		if err != nil {
			http.Error(w, "failed to render template", http.StatusInternalServerError)
		}
	}
}
