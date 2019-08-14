package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorWrapper struct {
	Status string `json:"status"`
	Msg    string `json:"message"`
	Err    string `json:"error"`
}

// Facilitates handlers in the sending of helpful error responses over http
func RespondWithError(msg string, err error, statusCode int, w http.ResponseWriter) {
	jsonObj, _ := json.Marshal(ErrorWrapper{
		Status: http.StatusText(statusCode),
		Msg:    msg,
		Err:    err.Error(),
	})

	w.WriteHeader(statusCode)
	w.Write(jsonObj)
}
