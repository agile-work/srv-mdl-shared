package middlewares

import (
	"net/http"

	"github.com/agile-work/srv-mdl-shared/models/translation"
)

// Translation sets field request language code
func Translation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		translation.FieldsRequestLanguageCode = r.Header.Get("Content-Language")

		next.ServeHTTP(w, r)
	})
}
