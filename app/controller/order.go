package controller

import (
	"../model"
	"log"
	"net/http"
)

// GetOrders : List of orders
func GetOrders(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	output, err := order.GetData()
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}


