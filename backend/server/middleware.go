package server

import "net/http"

const maxBodyBytes = 1 << 16

// maxBodySizeMiddleware limits the size of request bodies to protect against DoS.
// - If Content-Length > maxBodyBytes: rejects with 413.
// - Otherwise wraps the body so reads beyond the limit fail.
func maxBodySizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > maxBodyBytes {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

		next.ServeHTTP(w, r)
	})
}
