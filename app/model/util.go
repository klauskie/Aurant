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

func jsonResponse(w http.ResponseWriter, output []*Restaurant) http.ResponseWriter {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}
	return w
}