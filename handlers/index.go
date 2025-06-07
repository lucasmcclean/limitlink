package handlers

import (
	"io/fs"
	"net/http"
)

func Index(staticFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, staticFS, "index.html")
	}
}
