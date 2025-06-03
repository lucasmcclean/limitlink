package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/lucasmcclean/limitlink/link"
)


func Redirect(ctx context.Context, repo link.Repository) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
    slug := strings.TrimPrefix(r.URL.Path, "/")
		if slug == "" {
			http.Error(w, "missing slug", http.StatusBadRequest)
			return
		}

		lnk, err := repo.GetAndInc(ctx, slug)
		if err != nil {
			http.Error(w, "link not found", http.StatusNotFound)
			return
		}

		target := lnk.Target
		if target == "" {
			http.Error(w, "invalid target", http.StatusInternalServerError)
		}

		http.Redirect(w, r, target, http.StatusFound)
	}
}
