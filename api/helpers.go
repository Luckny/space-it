package api

import (
	"encoding/json"
	"net/http"

	"github.com/Luckny/space-it/util"
)

func writeJSON(w http.ResponseWriter, r *http.Request, status int, value any) error {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("X-Frame-Options", "DENY")
	w.Header().Add("X-XSS-Protection", "0")
	w.Header().Add("Cache-Control", "no-store")
	w.WriteHeader(status)

	util.InfoLog.Println(" ---> ", r.URL.Path, status)
	return json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, r *http.Request, status int, error error) {
	if status == http.StatusInternalServerError {
		writeJSON(w, r, status, map[string]string{"error": "internal server error"})
		util.ErrorLog.Println(" ---> ", r.URL.Path, status, error)
		return
	}
	writeJSON(w, r, status, map[string]string{"error": error.Error()})
}
