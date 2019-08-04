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
	r.HandleFunc("/item/restaurant/{rest_id}", controller.GetItemsByRestaurant).Methods("GET")
	r.HandleFunc("/item/detail/{item_id}", controller.GetItemByID).Methods("GET")
	r.HandleFunc("/item/update", controller.UpdateItem).Methods("POST")
	r.HandleFunc("/item/delete/{item_id}", controller.DeleteItem).Methods("GET")
	r.HandleFunc("/item/categories/{rest_id}", controller.GetCategoriesByRestaurant).Methods("GET")

	r.HandleFunc("/order", controller.GetOrders).Methods("GET")
	r.HandleFunc("/order", controller.SetOrder).Methods("POST")
	r.HandleFunc("/order/update", controller.UpdateOrderState).Methods("POST")
	r.HandleFunc("/order/update/increment/{order_id}", controller.UpdateOrderStateIncrement).Methods("GET")
	r.HandleFunc("/order/restaurant/{rest_id}/state/{state}", controller.GetOrdersByState).Methods("GET")
	r.HandleFunc("/order/client/{email}/restaurant/{rest_id}", controller.GetOrdersByClientAndRest).Methods("GET")

	http.Handle("/", r)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
