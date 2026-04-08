package middleware

import (
	"net/http"
	"strings"
)

func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			method := r.FormValue("_method")
			method = strings.ToUpper(method)

			switch method {
			case http.MethodPut, http.MethodPatch, http.MethodDelete:
				r.Method = method
			}
		}

		next.ServeHTTP(w, r)
	})
}
