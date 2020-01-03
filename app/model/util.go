package model

import (
	"encoding/json"
	"log"
	"net/http"
)

type Foul struct {
	Message       string      `json:"message"`
	Type     	  string      `json:"type"`
}

const CONTENT_TYPE  = "content-type"
const APPLICATION_JSON = "application/json"

// Writer for any type
func JsonResponseAny(w http.ResponseWriter, output interface{}) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}
}