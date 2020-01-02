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

	// query: by restaurant, by id,
	r.HandleFunc("/item", controller.ItemRouterGet).Methods("GET").
		Queries("item_id", "{item_id:[0-9]+}",
			"action", "{action:detail|delete|enable|disable}")
	r.HandleFunc("/item", controller.GetItems).Methods("GET")
	r.HandleFunc("/item", controller.SetItem).Methods("POST")
	r.HandleFunc("/item/restaurant/{rest_id:[0-9]+}", controller.GetItemsByRestaurant).Methods("GET")
	r.HandleFunc("/item/update", controller.UpdateItem).Methods("POST")
	r.HandleFunc("/item/category", controller.SetCategory).Methods("POST")
	r.HandleFunc("/item/category/restaurant/{rest_id:[0-9]+}", controller.GetCategoriesByRestaurant).Methods("GET")

	r.HandleFunc("/test", testHandler).Methods("GET").Queries("foo", "{foo:[a-z]+}", "id", "{id:[0-9]+|[A-Z]+|id}")

	r.HandleFunc("/order", controller.GetOrders).Methods("GET")
	r.HandleFunc("/order", controller.SetOrder).Methods("POST")
	r.HandleFunc("/order/update", controller.UpdateOrderState).Methods("POST")
	r.HandleFunc("/order/update/increment/{order_id:[0-9]+}", controller.UpdateOrderStateIncrement).Methods("GET")
	r.HandleFunc("/order/restaurant/{rest_id:[0-9]+}/state/{state:[1-4]}", controller.GetOrdersByState).Methods("GET")
	r.HandleFunc("/order/client/{email}/restaurant/{rest_id:[0-9]+}", controller.GetOrdersByClientAndRest).Methods("GET")

	http.Handle("/", r)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	//params := r.URL.Query()

	fmt.Fprint(w, vars)
}
