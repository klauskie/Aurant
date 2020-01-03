package controller

import (
	"../config"
	"../model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// ItemRouter : Entry point for get requests
func ItemRouterGet(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	qParams := mux.Vars(r)

	itemId := qParams[config.TAG_ITEM_ID]
	action := qParams[config.TAG_ACTION]

	var output []*model.Item
	var err model.FoulError

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
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}

	model.JsonResponseAny(w, output)
}

// GetItems : List of items
func GetItems(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	output, err := item.GetData()
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}

	model.JsonResponseAny(w, output)
}

// GetItemsByRestaurant : List of items by restaurant ID and its categories
func GetItemsByRestaurant(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	output, err := item.GetDataByRestID(vars[config.TAG_RESTAURANT_ID])
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}

	model.JsonResponseAny(w, output)
}

// GetItemByID : item detail by ID
func GetItemByID(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	output, err := item.GetDataByID(vars[config.TAG_ITEM_ID])
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}

	model.JsonResponseAny(w, output)
}

// GetCategoriesByRestaurant : return categories grouped by category_id
func GetCategoriesByRestaurant(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	vars := mux.Vars(r)

	output, err := item.GetCategoriesByRestaurant(vars[config.TAG_RESTAURANT_ID])
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}

	model.JsonResponseAny(w, output)
}


// SetItem : Insert item
func SetItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	err := item.SetData(r.Body)
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}
}

// SetCategory : Insert category
func SetCategory(w http.ResponseWriter, r *http.Request) {
	var cate model.Category

	err := cate.SetData(r.Body)
	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}
}

// UpdateItem : Update item by ID
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	err := item.UpdateData(r.Body)

	if err != nil {
		log.Println(err)
		model.JsonResponseError(w, err)
		return
	}
}