package util

import (
	"encoding/json"
	"net/http"
)

type idResponse struct {
	ID string `json:"id"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func WriteID(rw http.ResponseWriter, status int, id string) {
	WriteResponse(rw, status, &idResponse{id})
}

func WriteMessage(rw http.ResponseWriter, status int, msg string) {
	WriteResponse(rw, status, &messageResponse{msg})
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
