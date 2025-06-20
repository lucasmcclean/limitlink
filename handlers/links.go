package handlers

import (
	"context"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/lucasmcclean/limitlink/link"
)

func Links(ctx context.Context, repo link.Repository, templatesFS fs.FS) http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(templatesFS, "new-link.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		lnk, err := link.NewFromForm(r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repo.Create(ctx, lnk); err != nil {
			log.Printf("error storing link: %v", err)
			http.Error(w, "failed to store link", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		if err := tmpl.ExecuteTemplate(w, "new-link", lnk); err != nil {
			log.Printf("template execution error: %v", err)
			http.Error(w, "template error", http.StatusInternalServerError)
		}
	}
}
