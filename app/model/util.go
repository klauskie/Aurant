package model

import (
	"../config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type Foul struct {
	Type     	string      `json:"type"`
	Message		string      `json:"message"`
	From		string		`json:"from"`
}

type FoulError interface {
	Render()
}

const CONTENT_TYPE  = "content-type"
const APPLICATION_JSON = "application/json"

// Get Foul for different setups [develop, prod]
func (foul *Foul) Render() {
	if config.IsProductionEnabled() {
		foul.From = ""
		foul.Type = ""
	}
}

// TraceCall : returns info from caller
func TraceCall() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function)
}

// Writer for any type
func JsonResponseAny(w http.ResponseWriter, output interface{}) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}
}

// Writer for errors
func JsonResponseError(w http.ResponseWriter, output FoulError) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	w.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}
}