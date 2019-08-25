package controller

import (
	"../model"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// GetItems : List of items
func GetItems(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	output, err := item.GetData()
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

// GetItemsByRestaurant : List of items by restaurant
func GetItemsByRestaurant(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	output, err := item.GetDataByRestID(vars["rest_id"])
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

// GetItemByID : item detail by ID
func GetItemByID(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	output, err := item.GetDataByID(vars["item_id"])
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

// GetCategoriesByRestaurant : return categories grouped by category_id
func GetCategoriesByRestaurant(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	output, err := item.GetCategoriesByRestaurant(vars["rest_id"])
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}


// SetItem : Insert item
func SetItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	item.SetData(r.Body)
	fmt.Fprintln(w, "yey")
}

// SetCategory : Insert category
func SetCategory(w http.ResponseWriter, r *http.Request) {
	var cate model.Category

	cate.SetData(r.Body)
	fmt.Fprintln(w, "yey")
}

// UpdateItem : Update item by ID
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	err := item.UpdateData(r.Body)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintln(w, "yey")
}

// DeleteItem : delete item by ID
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	output, err := item.DeleteData(vars["item_id"])
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

// EnableItem : enable item by id
func EnableItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)
	usableID, _ := strconv.Atoi(vars["item_id"])
	item.ID = usableID

	output, err := item.Enabletor(true)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

// DisableItem : disable item by id
func DisableItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)
	usableID, _ := strconv.Atoi(vars["item_id"])
	item.ID = usableID

	output, err := item.Enabletor(false)
	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}
