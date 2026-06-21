package utils

import (
	"encoding/json"
	"net/http"
)

func SetJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type Response struct {
	Data any    `json:"data,omitempty"`
	Err  string `json:"err,omitempty"`
}

func SendJSON(w http.ResponseWriter, status int, res Response) {
	w.WriteHeader(status)

	data, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Write(data)
}
