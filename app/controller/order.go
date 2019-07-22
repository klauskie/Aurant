package controller

import (
	"../model"
	"fmt"
	"github.com/gorilla/mux"
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

// UpdateOrderState : update order
func UpdateOrderState(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	err := order.SetNewState(r.Body)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintln(w, "yey")
}

// GetOrdersByState : update order
func GetOrdersByState(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	vars := mux.Vars(r)

	output, err := order.GetDataByRestIDAndState(vars["rest_id"], vars["state"])
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

// UpdateOrderStateIncrement : update order status increment
func UpdateOrderStateIncrement(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	vars := mux.Vars(r)

	output, err := order.UpdateStatusByOne(vars["order_id"])
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}