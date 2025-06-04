package handlers

import (
	"io/fs"
	"log"
	"net/http"
)

func Index(staticFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, _ := fs.ReadDir(staticFS, "static")
		for _, entry := range entries {
			log.Println("entry:", entry.Name())
		}

		http.ServeFileFS(w, r, staticFS, "index.html")
	}
}
