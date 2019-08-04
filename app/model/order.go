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

// Cart schema for db
type Cart struct {
	ID       int     `json:"cart_id"`
	RestID 	 int     `json:"rest_id"`
	ItemID 	 int     `json:"item_id"`
	ClientID int     `json:"email"`
	State    int     `json:"state_id"`
	Date     string  `json:"datetime"`
	AdInfo   string  `json:"additional_info"`
}

var timeLayout string = "2006-01-02 15:04:05"

const (
	NEW         = iota
	ON_PROGRESS = iota
	DELIVERED   = iota
	CLOSED      = iota
)

// GetData : get all orders
func (cart *Cart) GetData() ([]byte, error) {
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
func (order *Cart) GetDataByRestIDAndState(rest_id string, state string) ([]byte, error) {

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

// GetDataByClientAndRestID : call
func (order *Cart) GetDataByClientAndRestID(email string, restID string) ([]byte, error) {

	usableID, _ := strconv.Atoi(restID)

	data, err := getAllOrdersByClientAndRestID(email, usableID)
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
func (order *Cart) UpdateStatusByOne(orderID string) ([]byte, error) {

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
func (order *Cart) SetData(stream io.Reader) error {
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
func (order *Cart) SetNewState(stream io.Reader) error {
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
func getAllOrders() ([]*Cart, error) {
	var m []*Cart
	rows, err := config.DB.Query("SELECT * FROM CART")
	if err != nil {
		return m, err
	}
	return scanRowsOrder(rows)
}

// getAllOrdersByRestIDAndState : return map with orders
func getAllOrdersByRestIDAndState(rest_id int, state string) ([]*Cart, error) {
	var m []*Cart
	rows, err := config.DB.Query("SELECT * FROM CART WHERE rest_id = ? AND state_id = ?", rest_id, state)
	if err != nil {
		return m, err
	}
	return scanRowsOrder(rows)
}

// getAllOrdersByClientAndRestID : return orders
func getAllOrdersByClientAndRestID(email string, restID int) ([]*Cart, error) {
	var m []*Cart
	rows, err := config.DB.Query("SELECT * FROM CART WHERE rest_id = ? AND email = ? AND state_id != ?", restID, email, CLOSED)
	if err != nil {
		return m, err
	}
	return scanRowsOrder(rows)
}


// scanRowsOrder : scan rows
func scanRowsOrder(rows *sql.Rows) ([]*Cart, error) {
	var m []*Cart

	for rows.Next() {
		var cart Cart

		err := rows.Scan(&cart.ID, &cart.RestID, &cart.ItemID, &cart.ClientID, &cart.State, &cart.Date, &cart.AdInfo)
		if err != nil {
			return m, err
		}

		m = append(m, &cart)
	}

	return m, nil
}

func (order *Cart) insertIntoDB() error {
	_, err := config.DB.Exec("INSERT INTO CART(rest_id, item_id, email, state_id, datetime, additional_info) VALUES (?,?,?,?,?,?)",
		order.RestID, order.ItemID, order.ClientID, NEW, time.Now().Format(timeLayout), order.AdInfo)
	if err != nil {
		return err
	}
	return nil
}

func (order *Cart) updateState() error {
	_, err := config.DB.Exec("UPDATE CART SET state_id = ?, date = ? WHERE cart_id = ?", order.State, time.Now().Format(timeLayout), order.ID)
	if err != nil {
		return err
	}
	return nil
}

func incrementStateByOne (cartID int) (map[string]string, error) {

	var actualState int
	message := make(map[string]string)

	row := config.DB.QueryRow("SELECT state_id FROM CART WHERE cart_id = ?", cartID)
	err := row.Scan(&actualState)

	if err != nil && err.Error() != "sql: no rows in result set" {
		message["state"] = ""
		return message, err
	}

	if actualState <= CLOSED {
		actualState += 1
		_, err := config.DB.Exec("UPDATE CART SET state_id = ?, date = ? WHERE cart_id = ?", actualState, time.Now().Format(timeLayout), cartID)
		if err != nil {
			message["state"] = ""
			return message, err
		}
	}
	message["state"] = strconv.Itoa(actualState)
	return message,nil
}