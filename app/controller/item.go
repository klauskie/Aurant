package controller

import (
	"app/model"
	"fmt"
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

// SetItem : Insert restaurant
func SetItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	item.SetData(r.Body)
	fmt.Fprintln(w, "yey")
}
