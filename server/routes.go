package server

import (
	"net/http"

	"github.com/lucasmcclean/url-shortener/handler"
)

func registerRoutes(mux *http.ServeMux) {
	handler.ServeViews(mux)
}
