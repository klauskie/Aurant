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
	ClientID int     `json:"email"`
	State    int     `json:"state_id"`
	Date     string  `json:"datetime"`
}

// Cart_Item schema for db
type CartItem struct {
	ID       int     `json:"id"`
	CartID 	 int     `json:"cart_id"`
	ItemID 	 int     `json:"item_id"`
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
func getAllOrders() (map[string]Cart, error) {
	m := make(map[string]Cart)
	rows, err := config.DB.Query("SELECT * FROM CART")
	if err != nil {
		return m, err
	}
	return scanRowsOrder(rows)
}

// getAllOrdersByRestIDAndState : return map with orders
func getAllOrdersByRestIDAndState(rest_id int, state string) (map[string]Cart, error) {
	m := make(map[string]Cart)
	rows, err := config.DB.Query("SELECT * FROM CART WHERE rest_id = ? AND state = ?", rest_id, state)
	if err != nil {
		return m, err
	}
	return scanRowsOrder(rows)
}


// scanRowsOrder : scan rows
func scanRowsOrder(rows *sql.Rows) (map[string]Cart, error) {
	m := make(map[string]Cart)

	for rows.Next() {
		var cart Cart

		err := rows.Scan(&cart.ID, &cart.RestID, &cart.ClientID, &cart.State, &cart.Date)
		if err != nil {
			return m, err
		}

		strID := strconv.Itoa(cart.ID)
		m[strID] = cart
	}

	return m, nil
}

func (order *Cart) insertIntoDB() error {
	_, err := config.DB.Exec("INSERT INTO CART(rest_id, email, state_id, datetime) VALUES (?,?,?,?)", order.RestID, order.ClientID, NEW, time.Now().Format(timeLayout))
	if err != nil {
		return err
	}
	return nil
}

func (order *Cart) updateState() error {
	_, err := config.DB.Exec("UPDATE CART SET state_id = ?, date = ? WHERE order_id = ?", order.State, time.Now().Format(timeLayout), order.ID)
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
		_, err := config.DB.Exec("UPDATE CART SET state_id = ?, date = ? WHERE order_id = ?", actualState, time.Now().Format(timeLayout), cartID)
		if err != nil {
			message["state"] = ""
			return message, err
		}
	}
	message["state"] = strconv.Itoa(actualState)
	return message,nil
}