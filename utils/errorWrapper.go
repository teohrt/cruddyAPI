package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorWrapper struct {
	Msg string `json:"Message"`
	Err string `json:"Error"`
}

// Facilitates handlers in the sending of helpful error responses over http
func RespondWithError(msg string, err error, statusCode int, w http.ResponseWriter) {
	jsonObj, _ := json.Marshal(ErrorWrapper{
		Msg: msg,
		Err: err.Error(),
	})

	w.WriteHeader(statusCode)
	w.Write(jsonObj)
}
