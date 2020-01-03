package model

import (
	"database/sql"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"../config"
)

// Cart schema for db
type Cart struct {
	ID       int     `json:"cart_id"`
	RestID 	 int     `json:"rest_id"`
	ItemID 	 int     `json:"item_id"`
	ClientID string  `json:"email"`
	Date     string  `json:"datetime"`
	State    int     `json:"state_id"`
	AdInfo   string  `json:"additional_info"`
	FullItem Item 	 `json:"full_item"`
}

var timeLayout string = "2006-01-02 15:04:05"

const (
	ZERO		= iota
	NEW         = iota
	ON_PROGRESS = iota
	DELIVERED   = iota
	CLOSED      = iota
)

// GetData : get all orders
func (cart *Cart) GetData() ([]*Cart, FoulError) {
	data, err := getAllOrders()
	return data, err
}

// GetDataByRestIDAndState : call getAllOrdersByRestIDAndState
func (order *Cart) GetDataByRestIDAndState(rest_id string, state string) ([]*Cart, FoulError) {

	usableID, _ := strconv.Atoi(rest_id)

	data, err := getAllOrdersByRestIDAndState(usableID, state)
	return data, err
}

// GetDataByClientAndRestID : call getAllOrdersByClientAndRestID
func (order *Cart) GetDataByClientAndRestID(email string, restID string) ([]*Cart, FoulError) {

	usableID, _ := strconv.Atoi(restID)

	data, err := getAllOrdersByClientAndRestID(email, usableID)
	return data, err
}

// UpdateStatusByOne : update Status by one
func (order *Cart) UpdateStatusByOne(orderID string) FoulError {

	usableID, _ := strconv.Atoi(orderID)

	err := incrementStateByOne(usableID)
	return err
}

// SetData : post stream into order table
func (order *Cart) SetData(stream io.Reader) FoulError {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&order)
	if err != nil {
		panic(err)
	}

	errF := order.insertIntoDB()
	if errF != nil {
		return errF
	}
	return nil
}

// UpdateState : post stream into order table
func (order *Cart) SetNewState(stream io.Reader) FoulError {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&order)
	if err != nil {
		panic(err)
	}

	errF := order.updateState()
	if errF != nil {
		return errF
	}
	return nil
}

// getAllOrders : return map with orders
func getAllOrders() ([]*Cart, FoulError) {
	var m []*Cart
	rows, err := config.DB.Query("SELECT * FROM CART")
	if err != nil {
		return m, &Foul{"DB.QUERY ERROR", err.Error(), TraceCall()}
	}
	return scanRowsOrder(rows)
}

// getAllOrdersByRestIDAndState : return map with orders
func getAllOrdersByRestIDAndState(rest_id int, state string) ([]*Cart, FoulError) {
	var m []*Cart
	rows, err := config.DB.Query("SELECT * FROM CART WHERE rest_id = ? AND state_id = ?", rest_id, state)
	if err != nil {
		return m, &Foul{"DB.QUERY ERROR", err.Error(), TraceCall()}
	}
	return scanRowsOrder(rows)
}

// getAllOrdersByClientAndRestID : return orders
func getAllOrdersByClientAndRestID(email string, restID int) ([]*Cart, FoulError) {
	var m []*Cart
	rows, err := config.DB.Query("SELECT * FROM CART WHERE rest_id = ? AND email = ? AND state_id != ?", restID, email, CLOSED)
	if err != nil {
		return m, &Foul{"DB.QUERY ERROR", err.Error(), TraceCall()}
	}
	return scanRowsOrder(rows)
}


// scanRowsOrder : scan rows
func scanRowsOrder(rows *sql.Rows) ([]*Cart, FoulError) {
	var m []*Cart

	for rows.Next() {
		var cart Cart

		err := rows.Scan(&cart.ID, &cart.RestID, &cart.ItemID, &cart.ClientID, &cart.Date, &cart.State, &cart.AdInfo)
		if err != nil {
			return m, &Foul{"DB.SCAN ERROR", err.Error(), TraceCall()}
		}

		errF := implantFullItem(&cart)
		if errF != nil {
			return m, errF
		}

		m = append(m, &cart)
	}

	return m, nil
}

// implantFullItem : find the item and add it to the cart
func implantFullItem(cart *Cart) FoulError {

	row := config.DB.QueryRow("SELECT * FROM ITEM WHERE item_id = ?", cart.ItemID)
	err := row.Scan(&cart.FullItem.ID, &cart.FullItem.RestID, &cart.FullItem.CategoryID, &cart.FullItem.Name,
		&cart.FullItem.Description, &cart.FullItem.Price, &cart.FullItem.IsEnabled)

	if err != nil {
		return &Foul{"DB.SCAN ERROR", err.Error(), TraceCall()}
	}

	return nil
}

func (order *Cart) insertIntoDB() FoulError {
	_, err := config.DB.Exec("INSERT INTO CART(rest_id, item_id, email, datetime, state_id, additional_info) VALUES (?,?,?,?,?,?)",
		order.RestID, order.ItemID, order.ClientID, time.Now().Format(timeLayout), NEW, order.AdInfo)
	if err != nil {
		return &Foul{"DB.EXEC ERROR", err.Error(), TraceCall()}
	}
	return nil
}

func (order *Cart) updateState() FoulError {
	_, err := config.DB.Exec("UPDATE CART SET state_id = ?, datetime = ? WHERE cart_id = ?", order.State, time.Now().Format(timeLayout), order.ID)
	if err != nil {
		return &Foul{"DB.EXEC ERROR", err.Error(), TraceCall()}
	}
	return nil
}

func incrementStateByOne (cartID int) FoulError {

	var actualState int

	row := config.DB.QueryRow("SELECT state_id FROM CART WHERE cart_id = ?", cartID)
	err := row.Scan(&actualState)

	if err != nil && err.Error() != "sql: no rows in result set" {
		return &Foul{"DB.EXEC ERROR", err.Error(), TraceCall()}
	}

	if actualState <= CLOSED {
		actualState += 1
		_, err := config.DB.Exec("UPDATE CART SET state_id = ?, datetime = ? WHERE cart_id = ?", actualState, time.Now().Format(timeLayout), cartID)
		if err != nil {
			return &Foul{"DB.EXEC ERROR", err.Error(), TraceCall()}
		}
	}
	return nil
}