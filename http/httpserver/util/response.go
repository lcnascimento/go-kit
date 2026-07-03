package util

import (
	"encoding/json"
	"net/http"
)

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
