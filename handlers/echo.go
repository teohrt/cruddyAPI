package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]
	fmt.Fprintf(w, word)
}
