package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"app/model"

	"github.com/gorilla/mux"
)

type Message struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Greet struct {
	Greeting string `json:"greeting"`
	Name     string `json:"name"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Hello).Methods("GET")
	r.HandleFunc("/post", Post).Methods("POST")
	r.HandleFunc("/res", PostRes).Methods("GET")
	r.HandleFunc("/res", GetRes).Methods("POST")
	http.Handle("/", r)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello Aurant!")
}

func Post(w http.ResponseWriter, req *http.Request) {

	msg, err := handleRequestJSON(req.Body)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	greet := Greet{
		"Hello",
		msg.Name}

	output, err2 := handleResponseJSON(greet)
	if err2 != nil {
		log.Fatal("Encoding error: ", err2)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func GetRes(w http.ResponseWriter, req *http.Request) {
	var res model.Restaurant
	res.GetData(2)
}

func PostRes(w http.ResponseWriter, req *http.Request) {
	var res model.Restaurant
	res.SetData(req.Body)
}

func handleRequestJSON(body io.Reader) (Message, error) {
	var msg Message
	jsn, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}
	err = json.Unmarshal(jsn, &msg)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}
	return msg, err
}

func handleResponseJSON(grt Greet) ([]byte, error) {
	output, err := json.Marshal(grt)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}
	return output, err
}
