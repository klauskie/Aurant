package controller

import (
	"../model"
	"fmt"
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

// SetOrder : Insert order
func SetOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	err := order.SetData(r.Body)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintln(w, "yey")
}


