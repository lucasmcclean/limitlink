package handler

import "net/http"

func ServeViews(mux *http.ServeMux) {
	staticFS := http.FileServer(http.Dir("./static"))

  mux.Handle("/static/", http.StripPrefix("/static/", staticFS))

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
}
