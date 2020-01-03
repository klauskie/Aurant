package controller

import (
	"../model"
	"log"
	"net/http"
)

// GetRestaurants : List of restaurants
func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	var res model.Restaurant

	output, err := res.GetData()
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}

	model.JsonResponseAny(w, output)
}

// SetRestaurant : Insert restaurant
func SetRestaurant(w http.ResponseWriter, r *http.Request) {
	var res model.Restaurant

	err := res.SetData(r.Body)
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}
}
