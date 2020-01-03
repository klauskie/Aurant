package controller

import (
	"../model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// ItemRouter : Entry point for get requests
func ItemRouterGet(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	qParams := mux.Vars(r)

	itemId := qParams["item_id"]
	action := qParams["action"]

	var output []*model.Item
	var err error

	switch action {
	case "detail":
		output, err = item.GetDataByID(itemId)
		break
	case "delete":
		err = item.DeleteData(itemId)
		break
	case "enable":
		err = item.Enabletor(itemId, true)
		break
	case "disable":
		err = item.Enabletor(itemId, false)
		break
	}

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

// GetItems : List of items
func GetItems(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	output, err := item.GetData()
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

// GetItemsByRestaurant : List of items by restaurant ID and its categories
func GetItemsByRestaurant(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	output, err := item.GetDataByRestID(vars["rest_id"])
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

// GetItemByID : item detail by ID
/*func GetItemByID(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	output, err := item.GetDataByID(vars["item_id"])
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}*/

// GetCategoriesByRestaurant : return categories grouped by category_id
func GetCategoriesByRestaurant(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	output, err := item.GetCategoriesByRestaurant(vars["rest_id"])
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


// SetItem : Insert item
func SetItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	err := item.SetData(r.Body)
	if err != nil {
		log.Fatal("Insertion error: ", err)
	}
}

// SetCategory : Insert category
func SetCategory(w http.ResponseWriter, r *http.Request) {
	var cate model.Category

	err := cate.SetData(r.Body)
	if err != nil {
		log.Fatal("Insertion error: ", err)
	}
}

// UpdateItem : Update item by ID
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	err := item.UpdateData(r.Body)

	if err != nil {
		log.Fatal("Update Item error: ", err)
	}
}