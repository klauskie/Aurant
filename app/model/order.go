package model

import (
	"encoding/json"
	"io"
	"log"
	"strconv"

	"../config"
)

// Order schema for db
type Order struct {
	ID       int     `json:"order_id"`
	ItemID 	 int     `json:"item_id"`
	RestID 	 int     `json:"rest_id"`
	ClientID int     `json:"client_id"`
	State    string  `json:"state"`
	Date     string  `json:"date"`
}

func (order *Order) GetData() ([]byte, error) {
	data, err := getAllOrders()
	if err != nil {
		log.Fatal("getAllOrders error: ", err)
	}
	output, err2 := json.Marshal(data)
	if err2 != nil {
		log.Fatal("Encoding error: ", err2)
	}
	return output, err2
}

// SetData : post stream into order table
func (order *Order) SetData(stream io.Reader) error {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&order)
	if err != nil {
		panic(err)
	}

	err = order.insertIntoDB()
	if err != nil {
		return err
	}
	return nil
}

// getAllOrders : return map with orders
func getAllOrders() (map[string]Order, error) {
	m := make(map[string]Order)
	rows, err := config.DB.Query("SELECT * FROM `order`")
	if err != nil {
		return m, err
	}
	for rows.Next() {
		var order Order

		err := rows.Scan(&order.ID, &order.ItemID, &order.RestID, &order.ClientID, &order.State, &order.Date)
		if err != nil {
			return m, err
		}

		strID := strconv.Itoa(order.ID)
		m[strID] = order
	}

	return m, nil
}

func (order *Order) insertIntoDB() error {
	_, err := config.DB.Exec("INSERT INTO `order`(item_id, rest_id, client_id, state, date) VALUES (?,?,?,?,?)", order.ItemID, order.RestID, order.ClientID, order.State, order.Date)
	if err != nil {
		return err
	}
	return nil
}