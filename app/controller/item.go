package controller

import (
	"../model"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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

// SetItem : Insert item
func SetItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	item.SetData(r.Body)
	fmt.Fprintln(w, "yey")
}

// UpdateItem : Update item by ID
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	err := item.UpdateData(r.Body, vars["item_id"])

	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintln(w, "yey")
}

// SetAttribute : Insert attrubute
func SetAttribute(w http.ResponseWriter, r *http.Request) {
	var att model.Attribute

	err := att.SetData(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	fmt.Fprintln(w, "yey")
}

// UpdateAttribute : Update attrubute
func UpdateAttribute(w http.ResponseWriter, r *http.Request) {
	var att model.Attribute

	err := att.UpdateData(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	fmt.Fprintln(w, "yey")
}
