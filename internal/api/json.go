package api

import "net/http"

func methodAllowed(w http.ResponseWriter, r *http.Request, methods ...string) bool {
	for _, m := range methods {
		if r.Method == m {
			return true
		}
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	return false
}
func health(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }
