package controller

import (
	"../model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// GetOrders : List of orders
func GetOrders(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	output, err := order.GetData()
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

// SetOrder : Insert order
func SetOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	err := order.SetData(r.Body)
	if err != nil {
		log.Fatal("Insertion error: ", err)
	}
}

// UpdateOrderState : update order
func UpdateOrderState(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	err := order.SetNewState(r.Body)
	if err != nil {
		log.Fatal("Update error: ", err)
	}
}

// GetOrdersByState : update order
func GetOrdersByState(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	vars := mux.Vars(r)

	output, err := order.GetDataByRestIDAndState(vars["rest_id"], vars["state"])
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

// UpdateOrderStateIncrement : update order status increment
func UpdateOrderStateIncrement(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	vars := mux.Vars(r)

	err := order.UpdateStatusByOne(vars["order_id"])
	if err != nil {
		log.Fatal("Internal error: ", err)
	}
}

// GetOrdersByClientAndRest : get open orders by client and rest
func GetOrdersByClientAndRest(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	vars := mux.Vars(r)

	output, err := order.GetDataByClientAndRestID(vars["email"], vars["rest_id"])
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