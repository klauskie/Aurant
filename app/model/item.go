package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strconv"

	"../config"
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
	RestID  int    `json:"rest_id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
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
func (item *Item) UpdateData(stream io.Reader, itemID string) error {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&item)
	if err != nil {
		return err
	}

	tempID,_ := strconv.Atoi(itemID)

	if item.ID != tempID {
		return errors.New("Item IDs does not match.")
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

// SetData :  set data for new attribute
func (att *Attribute) SetData(stream io.Reader) error {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&att)
	if err != nil {
		panic(err)
	}
	log.Println(att.Label)

	err = att.insertIntoDB()
	if err != nil {
		return err
	}
	return nil
}

// UpdateData :  update data for given attribute
func (att *Attribute) UpdateData(stream io.Reader) error {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&att)
	if err != nil {
		panic(err)
	}

	err = att.updateAttributeData()
	if err != nil {
		return err
	}
	return nil
}

// insertIntoDB : see if the attribute exists (creates a new one if it doesn't), then maps it to the item
func (att *Attribute) insertIntoDB() error {

	var tempLabel sql.NullString

	row := config.DB.QueryRow("SELECT att_id FROM attribute_value WHERE att_id = ?", att.AttID)
	err := row.Scan(&tempLabel)

	if err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}

	if !tempLabel.Valid {
		_, err = config.DB.Exec("INSERT INTO attribute_value(label, rest_id) VALUES (?,?)", att.Label, att.RestID)
		if err != nil {
			return err
		}

		row = config.DB.QueryRow("SELECT att_id FROM attribute_value where rest_id = ? AND label = ? ORDER BY att_id DESC;", att.RestID, att.Label)
		err = row.Scan(&att.AttID)

		if err != nil && err.Error() != "sql: no rows in result set" {
			return err
		}
	}

	_, err = config.DB.Exec("INSERT INTO item_attribute(item_id, att_id, value) VALUES (?,?,?)", att.ItemID, att.AttID, att.Value)
	if err != nil {
		return err
	}
	return nil
}

// updateAttributeData : update the fields if necessary
func (att *Attribute) updateAttributeData() error {

	var actualAtt Attribute

	row := config.DB.QueryRow("SELECT item_id, ia.att_id, rest_id, label, value from attribute_value join item_attribute ia ON attribute_value.att_id = ia.att_id where item_id = ? AND ia.att_id = ?", att.ItemID, att.AttID)
	err := row.Scan(&actualAtt.ItemID, &actualAtt.AttID, &actualAtt.RestID, &actualAtt.Label, &actualAtt.Value)
	if err != nil {
		return err
	}

	if actualAtt.Label != att.Label {
		_, err = config.DB.Exec("UPDATE attribute_value SET label = ? WHERE att_id = ?", att.Label, att.AttID)
		if err != nil {
			return err
		}
	}

	if actualAtt.Value != att.Value {
		_, err = config.DB.Exec("UPDATE item_attribute SET value = ? WHERE item_id = ? AND att_id = ?", att.Value, att.ItemID, att.AttID)
		if err != nil {
			return err
		}
	}

	return nil

}

func (item *Item) insertIntoDB() error {
	_, err := config.DB.Exec("INSERT INTO Item(rest_id, name, price) VALUES (?,?,?)", item.RestID, item.Name, item.Price)
	if err != nil {
		return err
	}
	return nil
}

func (item *Item) updateItemAction() error {
	_, err := config.DB.Exec("UPDATE Item SET name = ?, price = ? WHERE item_id = ?", item.Name, item.Price, item.ID)
	if err != nil {
		return err
	}
	return nil
}

// getAllItems : return map with items
func getAllItems() (map[string]Item, error) {
	m := make(map[string]Item)
	rows, err := config.DB.Query("SELECT * FROM Item")
	if err != nil {
		return m, err
	}

	return scanRows(rows)
}

// getAllItemsByRestID : return map with items
func getAllItemsByRestID(rest_id int) (map[string]Item, error) {
	m := make(map[string]Item)
	rows, err := config.DB.Query("SELECT * FROM Item WHERE rest_id = ?", rest_id)
	if err != nil {
		return m, err
	}

	return scanRows(rows)
}

// getAllItemsByRestID : return map with items
func getItemByID(id int) (map[string]Item, error) {
	m := make(map[string]Item)
	rows, err := config.DB.Query("SELECT * FROM Item WHERE item_id = ?", id)
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

	rows, err := config.DB.Query("SELECT item_id, a.att_id, a.label, item_attribute.value FROM item_attribute JOIN attribute_value a ON item_attribute.att_id = a.att_id WHERE item_id = ?", itemID)
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

func deleteItem(id int) (map[string]string,error) {

	message := make(map[string]string)

	_, err := config.DB.Exec("DELETE FROM Item where item_id = ?", id)
	if err != nil {
		message["success"] = "false"
		return message, err
	}
	message["success"] = "true"
	return message, nil
}
