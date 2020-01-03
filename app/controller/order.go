package controller

import (
	"../config"
	"../model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// GetOrders : List of orders
func GetOrders(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	output, err := order.GetData()
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}

	model.JsonResponseAny(w, output)
}

// SetOrder : Insert order
func SetOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	err := order.SetData(r.Body)
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}
}

// UpdateOrderState : update order
func UpdateOrderState(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	err := order.SetNewState(r.Body)
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}
}

// GetOrdersByState : update order
func GetOrdersByState(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	vars := mux.Vars(r)

	output, err := order.GetDataByRestIDAndState(vars[config.TAG_RESTAURANT_ID], vars[config.TAG_STATE])
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}

	model.JsonResponseAny(w, output)
}

// UpdateOrderStateIncrement : update order status increment
func UpdateOrderStateIncrement(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	vars := mux.Vars(r)

	err := order.UpdateStatusByOne(vars[config.TAG_ORDER_ID])
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}
}

// GetOrdersByClientAndRest : get open orders by client and rest
func GetOrdersByClientAndRest(w http.ResponseWriter, r *http.Request) {
	var order model.Cart

	vars := mux.Vars(r)

	output, err := order.GetDataByClientAndRestID(vars[config.TAG_EMAIL], vars[config.TAG_RESTAURANT_ID])
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}

	model.JsonResponseAny(w, output)
}