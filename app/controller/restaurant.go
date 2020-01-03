package controller

import (
	"../model"
	"encoding/json"
	"log"
	"net/http"
)

// GetRestaurants : List of restaurants
func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	var res model.Restaurant

	output, err := res.GetData()
	if err != nil {
		log.Fatal("Internal error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}
}

// SetRestaurant : Insert restaurant
func SetRestaurant(w http.ResponseWriter, r *http.Request) {
	var res model.Restaurant

	err := res.SetData(r.Body)
	if err != nil {
		log.Fatal("Insertion error: ", err)
	}
}
