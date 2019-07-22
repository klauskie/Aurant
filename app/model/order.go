package model

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"time"

	"../config"
)

// Order schema for db
type Order struct {
	ID       int     `json:"order_id"`
	ItemID 	 int     `json:"item_id"`
	RestID 	 int     `json:"rest_id"`
	ClientID int     `json:"client_id"`
	State    int     `json:"state"`
	Date     string  `json:"date"`
}

var timeLayout string = "2006-01-02 15:04:05"

const (
	new         = iota
	on_progress = iota
	delivered   = iota
	closed      = iota
)

// GetData : get all orders
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

// GetDataByRestIDAndState : call getAllItemsByRestID
func (order *Order) GetDataByRestIDAndState(rest_id string, state string) ([]byte, error) {

	usableID, _ := strconv.Atoi(rest_id)

	data, err := getAllOrdersByRestIDAndState(usableID, state)
	if err != nil {
		log.Fatal("getAllItems error: ", err)
	}
	output, err2 := json.Marshal(data)
	if err2 != nil {
		log.Fatal("Encoding error: ", err2)
	}
	return output, err2
}

// UpdateStatusByOne : update Status by one
func (order *Order) UpdateStatusByOne(orderID string) ([]byte, error) {

	usableID, _ := strconv.Atoi(orderID)

	data, err := incrementStateByOne(usableID)
	if err != nil {
		log.Fatal("incrementStateByOne error: ", err)
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

// UpdateState : post stream into order table
func (order *Order) SetNewState(stream io.Reader) error {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&order)
	if err != nil {
		panic(err)
	}

	err = order.updateState()
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
	return scanRowsOrder(rows)
}

// getAllOrdersByRestIDAndState : return map with orders
func getAllOrdersByRestIDAndState(rest_id int, state string) (map[string]Order, error) {
	m := make(map[string]Order)
	rows, err := config.DB.Query("SELECT * FROM `order` WHERE rest_id = ? AND state = ?", rest_id, state)
	if err != nil {
		return m, err
	}
	return scanRowsOrder(rows)
}


// scanRowsOrder : scan rows
func scanRowsOrder(rows *sql.Rows) (map[string]Order, error) {
	m := make(map[string]Order)

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
	_, err := config.DB.Exec("INSERT INTO `order`(item_id, rest_id, client_id, state, date) VALUES (?,?,?,?,?)", order.ItemID, order.RestID, order.ClientID, new, time.Now().Format(timeLayout))
	if err != nil {
		return err
	}
	return nil
}

func (order *Order) updateState() error {
	_, err := config.DB.Exec("UPDATE `order` SET state = ?, date = ? WHERE order_id = ?", order.State, time.Now().Format(timeLayout), order.ID)
	if err != nil {
		return err
	}
	return nil
}

func incrementStateByOne (orderID int) (map[string]string, error) {

	var actualState int
	message := make(map[string]string)

	row := config.DB.QueryRow("SELECT state FROM `order` WHERE order_id = ?", orderID)
	err := row.Scan(&actualState)

	if err != nil && err.Error() != "sql: no rows in result set" {
		message["state"] = ""
		return message, err
	}

	if actualState <= closed {
		actualState += 1
		_, err := config.DB.Exec("UPDATE `order` SET state = ?, date = ? WHERE order_id = ?", actualState, time.Now().Format(timeLayout), orderID)
		if err != nil {
			message["state"] = ""
			return message, err
		}
	}
	message["state"] = strconv.Itoa(actualState)
	return message,nil
}