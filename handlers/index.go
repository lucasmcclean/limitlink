package handlers

import "net/http"

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/root/static/index.html")
	}
}
