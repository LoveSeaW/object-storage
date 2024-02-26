package temp

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == http.MethodGet {
		put(w, r)
	}
	if method == http.MethodPatch {
		patch(w, r)
	}
	if method == http.MethodPost {
		post(w, r)
	}
	if method == http.MethodDelete {
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
