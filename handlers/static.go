package handlers

import "net/http"

func Static() http.Handler {
	return http.FileServer(http.Dir("static"))
}
