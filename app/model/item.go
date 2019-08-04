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

// GetData : call getAllItems
func (item *Item) GetData() ([]byte, error) {

	data, err := getAllItems()
	if err != nil {
		log.Fatal("getAllItems error: ", err)
	}
	output, err2 := json.Marshal(data)
	if err2 != nil {
		log.Fatal("Encoding error: ", err2)
	}
	return output, err2
}

// GetDataByRestID : call getAllItemsByRestID
func (item *Item) GetDataByRestID(id string) ([]byte, error) {

	usable_id, _ := strconv.Atoi(id)

	data, err := getAllItemsByRestID(usable_id)
	if err != nil {
		log.Fatal("getAllItems error: ", err)
	}
	output, err2 := json.Marshal(data)
	if err2 != nil {
		log.Fatal("Encoding error: ", err2)
	}
	return output, err2
}

// GetDataByID : call getItemByID
func (item *Item) GetDataByID(id string) ([]byte, error) {

	usableID, _ := strconv.Atoi(id)

	data, err := getItemByID(usableID)
	if err != nil {
		log.Fatal("getAllItems error: ", err)
	}
	output, err2 := json.Marshal(data)
	if err2 != nil {
		log.Fatal("Encoding error: ", err2)
	}
	return output, err2
}

// GetCategoriesByRestaurant
func (item *Item) GetCategoriesByRestaurant(id string) ([]byte, error) {
	usableID, _ := strconv.Atoi(id)

	data, err := getCategoriesByRestID(usableID)
	if err != nil {
		log.Fatal("getCategoriesByRestID error: ", err)
	}
	output, err2 := json.Marshal(data)
	if err2 != nil {
		log.Fatal("Encoding error: ", err2)
	}
	return output, err2
}

// SetData : post stream into db
func (item *Item) SetData(stream io.Reader) {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&item)
	if err != nil {
		panic(err)
	}
	log.Println(item.Name)

	err = item.insertIntoDB()
	if err != nil {
		panic(err)
	}
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
func (item *Item) DeleteData(id string) ([]byte, error) {

	usableID, _ := strconv.Atoi(id)

	data, err := deleteItem(usableID)
	if err != nil {
		log.Fatal("deleteItem error: ", err)
	}
	output, err2 := json.Marshal(data)
	if err2 != nil {
		log.Fatal("Encoding error: ", err2)
	}
	return output, err2
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

// getAllItems : return map with items
func getAllItems() (map[string]Item, error) {
	m := make(map[string]Item)
	rows, err := config.DB.Query("SELECT * FROM ITEM ORDER BY category_id")
	if err != nil {
		return m, err
	}

	return scanRows(rows)
}

// getAllItemsByRestID : return map with items
func getAllItemsByRestID(rest_id int) (map[string]Item, error) {
	m := make(map[string]Item)
	rows, err := config.DB.Query("SELECT * FROM ITEM WHERE rest_id = ? ORDER BY category_id", rest_id)
	if err != nil {
		return m, err
	}

	return scanRows(rows)
}

// getItemByID : return map with items
func getItemByID(id int) (map[string]Item, error) {
	m := make(map[string]Item)
	rows, err := config.DB.Query("SELECT * FROM ITEM WHERE item_id = ?", id)
	if err != nil {
		return m, err
	}

	return scanRows(rows)
}

// scanRows : scan rows and pass it into an object
func scanRows(rows *sql.Rows) (map[string]Item, error) {
	m := make(map[string]Item)

	for rows.Next() {
		var item Item

		err := rows.Scan(&item.ID, &item.RestID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.IsEnabled)
		if err != nil {
			return m, err
		}

		strID := strconv.Itoa(item.ID)
		m[strID] = item
	}

	return m, nil
}


func deleteItem(id int) (map[string]string,error) {

	message := make(map[string]string)

	_, err := config.DB.Exec("DELETE FROM ITEM where item_id = ?", id)
	if err != nil {
		message["success"] = "false"
		return message, err
	}
	message["success"] = "true"
	return message, nil
}

// getCategoriesByRestID : return map with categories
func getCategoriesByRestID(id int) (map[string]Category, error) {
	m := make(map[string]Category)
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

		strID := strconv.Itoa(cate.ID)
		m[strID] = cate
	}

	return m, nil
}
