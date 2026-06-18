package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data any
	Err  string
}

func SendJSON(w http.ResponseWriter, status int, res Response) {

	if res.Err != "" {
		SendJSON(w, status, Response{
			Err: res.Err,
		})
	}

	data, err := json.Marshal(res.Data)

	if err != nil {
		SendJSON(w, 500, Response{
			Err: "Internal server error",
		})
	}

	w.Write(data)
	w.WriteHeader(status)

}
