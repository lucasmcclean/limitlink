package handlers

import (
	"net/http"

	"github.com/lucasmcclean/limitlink/assets"
)

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, assets.StaticFS(), "index.html")
	}
}
