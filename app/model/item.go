package model

import (
	"encoding/json"
	"io"
	"log"
	"strconv"

	"app/config"
)

// Item schema for db
type Item struct {
	ID         int         `json:"item_id"`
	RestID     int         `json:"rest_id"`
	Name       string      `json:"name"`
	Price      string      `json:"price"`
	Attributes []Attribute `json:"attributes"`
}

// Attribute schema
type Attribute struct {
	ItemID int    `json:"item_id"`
	AttID  int    `json:"att_id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
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

func (item *Item) insertIntoDB() error {
	_, err := config.DB.Exec("INSERT INTO Item(rest_id, name, price) VALUES (?,?,?)", item.RestID, item.Name, item.Price)
	if err != nil {
		return err
	}
	return nil
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

// getAllItems : return map with items
func getAllItems() (map[string]Item, error) {
	m := make(map[string]Item)
	rows, err := config.DB.Query("SELECT * FROM Item")
	if err != nil {
		return m, err
	}
	for rows.Next() {
		var item Item

		err := rows.Scan(&item.ID, &item.RestID, &item.Name, &item.Price)
		if err != nil {
			return m, err
		}

		item.Attributes, err = getAttributes(item.ID)
		if err != nil {
			return m, err
		}

		strID := strconv.Itoa(item.ID)
		m[strID] = item
	}

	return m, nil
}

func getAttributes(itemID int) ([]Attribute, error) {
	var s []Attribute

	rows, err := config.DB.Query("SELECT item_id, a.att_id, a.value, item_attribute.value FROM item_attribute JOIN attribute_value a ON item_attribute.att_id = a.att_id WHERE item_id = ?", itemID)
	if err != nil {
		return s, err
	}
	for rows.Next() {
		var att Attribute

		err := rows.Scan(&att.ItemID, &att.AttID, &att.Label, &att.Value)
		if err != nil {
			return s, err
		}

		s = append(s, att)
	}
	return s, nil
}
