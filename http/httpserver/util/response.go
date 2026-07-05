package util

import (
	"encoding/json"
	"net/http"
)

func WriteMessage(rw http.ResponseWriter, status int, message string) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(status)
	rw.Write([]byte(message))
}

func WriteResponse(rw http.ResponseWriter, status int, response any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)

	if response == nil {
		return
	}

	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(response); err != nil {
		http.Error(rw, http.StatusText(status), status)

		return
	}
}
