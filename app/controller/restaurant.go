package controller

import (
	"app/model"
	"fmt"
	"log"
	"net/http"
)

// GetRestaurants : List of restaurants
func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	var res model.Restaurant

	output, err := res.GetData()
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

// SetRestaurant : Insert restaurant
func SetRestaurant(w http.ResponseWriter, r *http.Request) {
	var res model.Restaurant

	res.SetData(r.Body)
	fmt.Fprintln(w, "yey")
}
