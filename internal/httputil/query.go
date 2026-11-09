package httputil

import (
	"fmt"
	"net/http"
	"strings"
)

func GetRequiredQuery(w http.ResponseWriter, r *http.Request, param string) (string, bool) {
	val := strings.TrimSpace(r.URL.Query().Get(param))
	if val == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Query parameter '%s' is reuired"}`, param)
		return "", false
	}

	return val, true
}
