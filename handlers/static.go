package handlers

import (
	"net/http"

	"github.com/lucasmcclean/limitlink/assets"
)

func Static() http.Handler {
	return http.FileServer(http.FS(assets.StaticFS()))
}
