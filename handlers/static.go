package handlers

import (
	"io/fs"
	"net/http"
)

func Static(staticFS fs.FS) http.Handler {
	return http.FileServer(http.FS(staticFS))
}
