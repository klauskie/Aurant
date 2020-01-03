package model

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"strconv"

	"../config"
)

// Item schema for db
type Item struct {
	ID            int         `json:"item_id"`
	RestID     	  int         `json:"rest_id"`
	CategoryID    int         `json:"category_id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Price         string      `json:"price"`
	IsEnabled     bool        `json:"is_enabled"`
}

// Category schema for db
type Category struct {
	ID            int         `json:"category_id"`
	RestID     	  int         `json:"rest_id"`
	Name    	  string      `json:"name"`
}

// Categories and Items
type Bundle struct {
	Categories []*Category `json:"categories"`
	Items 	   []*Item 	   `json:"items"`
}

// GetData : call getAllItems
func (item *Item) GetData() ([]*Item, error) {

	data, err := getAllItems()
	if err != nil {
		log.Fatal("getAllItems error: ", err)
	}
	return data, err
}

// GetDataByRestID : call getAllItemsByRestID
// Append categories and items
func (item *Item) GetDataByRestID(id string) (Bundle, error) {

	usable_id, _ := strconv.Atoi(id)

	items, err := getAllItemsByRestID(usable_id)
	cats, err := getCategoriesByRestID(usable_id)

	data := Bundle {
		Categories:cats,
		Items:items,
	}

	if err != nil {
		log.Fatal("getAllItems error: ", err)
	}
	return data, err
}

// GetDataByID : call getItemByID
func (item *Item) GetDataByID(id string) ([]*Item, error) {

	usableID, _ := strconv.Atoi(id)

	data, err := getItemByID(usableID)
	if err != nil {
		log.Fatal("getAllItems error: ", err)
	}
	return data, err
}

// GetCategoriesByRestaurant
func (item *Item) GetCategoriesByRestaurant(id string) ([]*Category, error) {
	usableID, _ := strconv.Atoi(id)

	data, err := getCategoriesByRestID(usableID)
	if err != nil {
		log.Fatal("getCategoriesByRestID error: ", err)
	}
	return data, err
}

// SetData : post stream into db
func (item *Item) SetData(stream io.Reader) error {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&item)
	if err != nil {
		return err
	}
	log.Println(item.Name)

	err = item.insertIntoDB()
	if err != nil {
		return err
	}
	return nil
}

// SetData : post stream into db
func (cate *Category) SetData(stream io.Reader) error {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&cate)
	if err != nil {
		return err
	}

	err = cate.insertIntoDB()
	if err != nil {
		return err
	}
	return nil
}

// UpdateData : post stream into db
func (item *Item) UpdateData(stream io.Reader) error {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&item)
	if err != nil {
		return err
	}

	err = item.updateItemAction()
	if err != nil {
		return err
	}

	return nil
}

// DeleteData :
func (item *Item) DeleteData(id string) error {

	usableID, _ := strconv.Atoi(id)

	err := deleteItem(usableID)
	if err != nil {
		log.Fatal("deleteItem error: ", err)
	}
	return err
}

// DeleteData :
func (item *Item) Enabletor(item_id string, action bool) error {

	usableID, _ := strconv.Atoi(item_id)
	item.ID = usableID

	err := item.enableOrDisable(action)
	if err != nil {
		log.Fatal("deleteItem error: ", err)
	}
	return err
}

func (item *Item) insertIntoDB() error {
	_, err := config.DB.Exec("INSERT INTO ITEM(rest_id, category_id, name, description, price, is_enabled) VALUES (?,?,?,?,?,?)",
		item.RestID, item.CategoryID, item.Name, item.Description, item.Price, item.IsEnabled)
	if err != nil {
		return err
	}
	return nil
}

func (item *Item) updateItemAction() error {
	_, err := config.DB.Exec("UPDATE ITEM SET category_id = ?, name = ?, description = ?, price = ?, is_enabled = ? WHERE item_id = ?",
		item.CategoryID, item.Name, item.Description, item.Price, item.IsEnabled, item.ID)
	if err != nil {
		return err
	}
	return nil
}

func (cate *Category) insertIntoDB() error {
	_, err := config.DB.Exec("INSERT INTO CATEGORY(rest_id, name) VALUES (?,?)", cate.RestID, cate.Name)
	if err != nil {
		return err
	}
	return nil
}

// getAllItems : return map with items
func getAllItems() ([]*Item, error) {
	var m []*Item
	rows, err := config.DB.Query("SELECT * FROM ITEM ORDER BY category_id")
	if err != nil {
		return m, err
	}

	return scanRows(rows)
}

// getAllItemsByRestID : return map with items
func getAllItemsByRestID(rest_id int) ([]*Item, error) {
	var m []*Item
	rows, err := config.DB.Query("SELECT * FROM ITEM WHERE rest_id = ? ORDER BY category_id", rest_id)
	if err != nil {
		return m, err
	}

	return scanRows(rows)
}

// getItemByID : return map with items
func getItemByID(id int) ([]*Item, error) {
	var m []*Item
	rows, err := config.DB.Query("SELECT * FROM ITEM WHERE item_id = ?", id)
	if err != nil {
		return m, err
	}

	return scanRows(rows)
}

// scanRows : scan rows and pass it into an object
func scanRows(rows *sql.Rows) ([]*Item, error) {
	var m []*Item

	for rows.Next() {
		var item Item

		err := rows.Scan(&item.ID, &item.RestID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.IsEnabled)
		if err != nil {
			return m, err
		}

		m = append(m, &item)
	}

	return m, nil
}


func deleteItem(id int) error {

	_, err := config.DB.Exec("DELETE FROM ITEM where item_id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (item *Item) enableOrDisable(action bool) error {

	_, err := config.DB.Exec("UPDATE ITEM SET is_enabled = ? where item_id = ?", action, item.ID)
	if err != nil {
		return err
	}
	return nil
}

// getCategoriesByRestID : return map with categories
func getCategoriesByRestID(id int) ([]*Category, error) {
	var m []*Category
	rows, err := config.DB.Query("SELECT * from CATEGORY where rest_id = ?", id)
	if err != nil {
		return m, err
	}

	for rows.Next() {
		var cate Category

		err := rows.Scan(&cate.ID, &cate.RestID, &cate.Name)
		if err != nil {
			return m, err
		}

		m = append(m, &cate)
	}

	return m, nil
}
