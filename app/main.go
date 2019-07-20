package main

import (
	"fmt"
	"log"
	"net/http"

	"../app/controller"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/restaurant", controller.GetRestaurants).Methods("GET")
	r.HandleFunc("/restaurant", controller.SetRestaurant).Methods("POST")
	r.HandleFunc("/item", controller.GetItems).Methods("GET")
	r.HandleFunc("/item", controller.SetItem).Methods("POST")
	r.HandleFunc("/attribute", controller.SetAttribute).Methods("POST")
	r.HandleFunc("/orders", controller.GetOrders).Methods("GET")
	r.HandleFunc("/order", controller.SetOrder).Methods("POST")
	http.Handle("/", r)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
